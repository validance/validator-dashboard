package worker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	database "validator-dashboard/app/db"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TokenPriceTask struct {
	db             *sqlx.DB
	newTokenPrices []database.TokenPrice
}

type Quote struct {
	Price float64 `json:"price"`
}

type Currency struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Symbol string           `json:"symbol"`
	Slug   string           `json:"slug"`
	Quote  map[string]Quote `json:"quote"`
}

type QuotesApiResponseBody struct {
	Currencies map[string]Currency `json:"data"`
}

func NewTokenPriceTask(db *sqlx.DB) *TokenPriceTask {

	return &TokenPriceTask{
		db,
		nil,
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

func (t *TokenPriceTask) getNewTokenPrices(slugs []string) []database.TokenPrice {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sandbox-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Err(err).Msg("failed to create http request")
		return nil
	}

	q := url.Values{}
	for _, slug := range slugs {
		q.Add("slug", slug)
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("failed to receive token price from api")
		return nil
	}

	respBytes, _ := ioutil.ReadAll(resp.Body)

	respBody := QuotesApiResponseBody{}
	json.Unmarshal(respBytes, &respBody)

	tokenPrices := make([]database.TokenPrice, len(respBody.Currencies))

	for slug, currency := range respBody.Currencies {
		tokenPrice := database.TokenPrice{Chain: slug, Ticker: currency.Symbol, Price: currency.Quote["USD"].Price}
		tokenPrices = append(tokenPrices, tokenPrice)
	}

	return tokenPrices
}

func (t *TokenPriceTask) RunTokenPriceTask() {
	tasksNum := 1
	var wg sync.WaitGroup

	wg.Add(tasksNum)

	go func() {
		defer wg.Done()
		chains := t.getManagedChains()
		newTokenPrices := t.getNewTokenPrices(chains)
		if newTokenPrices != nil {
			t.createNewTokenPrices(newTokenPrices)
		}
	}()

	wg.Wait()
}
