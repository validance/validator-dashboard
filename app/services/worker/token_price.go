package worker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	database "validator-dashboard/app/db"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TokenPriceTask struct {
	db             *sqlx.DB
	newTokenPrices []database.TokenPrice
	coinGeckoIdMap map[string]string
}

type Coin struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	Slug         string  `json:"slug"`
	CurrentPrice float64 `json:"current_price"`
}

func NewTokenPriceTask(db *sqlx.DB) *TokenPriceTask {

	return &TokenPriceTask{
		db,
		nil,
		map[string]string{
			"juno": "juno-network",
		},
	}
}

func (t *TokenPriceTask) getManagedChains() []string {
	var chains []string

	queryErr := t.db.Select(
		&chains,
		`
			SELECT DISTINCT chain
			FROM delegation_history
		`,
	)

	if queryErr != nil {
		log.Err(queryErr).Msg("failed to get managed chains")
	}

	return chains
}

func (t *TokenPriceTask) createNewTokenPrices(tps []database.TokenPrice) {
	createQuery := `
		INSERT INTO token_price(chain, ticker, price)
		VALUES ($1, $2, $3)
	`

	for _, tp := range tps {
		_, err := t.db.Exec(createQuery, tp.Chain, tp.Ticker, tp.Price)
		if err != nil {
			log.Err(err).Msg("failed to create new token price data")
		}
	}
}

func (t *TokenPriceTask) getNewTokenPrices(chains []string) []database.TokenPrice {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/markets", nil)
	if err != nil {
		log.Err(err).Msg("failed to create http request")
		return nil
	}

	q := url.Values{}
	q.Add("vs_currency", "usd")

	chainIds := make([]string, len(chains))

	for i, chain := range chains {
		id := t.coinGeckoIdMap[chain]
		if id == "" {
			id = chain
		}
		chainIds[i] = id
	}
	q.Add("ids", strings.Join(chainIds, ","))

	req.Header.Set("Accepts", "application/json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("failed to receive token price from api")
		return nil
	}

	respBytes, _ := ioutil.ReadAll(resp.Body)

	coins := []Coin{}
	json.Unmarshal(respBytes, &coins)

	tokenPrices := make([]database.TokenPrice, len(coins))

	for i, coin := range coins {
		tokenPrice := database.TokenPrice{Chain: coin.Symbol, Ticker: coin.Symbol, Price: coin.CurrentPrice}
		tokenPrices[i] = tokenPrice
	}

	return tokenPrices
}

func (t *TokenPriceTask) RunTokenPriceTask() {
	chains := t.getManagedChains()

	newTokenPrices := t.getNewTokenPrices(chains)

	if newTokenPrices != nil {
		t.createNewTokenPrices(newTokenPrices)
	}
}
