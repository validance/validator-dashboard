package worker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
	"validator-dashboard/app/config"
	database "validator-dashboard/app/db"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TokenPriceTask struct {
	db             *sqlx.DB
	newTokenPrices []database.TokenPrice
}

type CurrentPriceResponse struct {
	Usd float64 `json:"usd"`
}

type MarketDataResponse struct {
	CurrentPrice CurrentPriceResponse `json:"current_price"`
}

type CoinHistoryResponse struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Symbol     string             `json:"symbol"`
	Slug       string             `json:"slug"`
	MarketData MarketDataResponse `json:"market_data"`
}

func NewTokenPriceTask(db *sqlx.DB) *TokenPriceTask {

	return &TokenPriceTask{
		db,
		nil,
	}
}

func (t *TokenPriceTask) createNewTokenPrices(tps []*database.TokenPrice) {
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

func (t *TokenPriceTask) getNewTokenPrice(chainId string) (*database.TokenPrice, error) {
	var tokenPrice *database.TokenPrice

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/"+chainId+"/history", nil)
	if err != nil {
		log.Err(err).Msg("failed to create http request")
		return nil, err
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("02-01-2006")

	q := url.Values{}
	q.Add("date", yesterday)
	q.Add("localization", "false")

	req.Header.Set("Accepts", "application/json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("failed to receive token price from api")
		return nil, err
	}

	respBytes, _ := ioutil.ReadAll(resp.Body)

	respBody := CoinHistoryResponse{}
	json.Unmarshal(respBytes, &respBody)

	tokenPrice = &database.TokenPrice{
		Chain:  respBody.Symbol,
		Ticker: respBody.Symbol,
		Price:  respBody.MarketData.CurrentPrice.Usd,
	}

	return tokenPrice, nil
}

func (t *TokenPriceTask) RunTokenPriceTask() {
	var newTokenPrices []*database.TokenPrice

	chainIds := config.GetConfig().CoingeckoIds

	tasksNum := len(chainIds)
	var wg sync.WaitGroup

	wg.Add(tasksNum)

	for _, chainId := range chainIds {
		chainId := chainId
		go func() {
			defer wg.Done()

			newTokenPrice, _ := t.getNewTokenPrice(chainId)
			if newTokenPrice != nil {
				newTokenPrices = append(newTokenPrices, newTokenPrice)
			}

		}()
	}

	wg.Wait()

	t.createNewTokenPrices(newTokenPrices)
}
