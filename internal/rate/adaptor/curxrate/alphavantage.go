package curxrate

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// AlphaVantage https://www.alphavantage.co/documentation/#fx
type AlphaVantage struct{ Provider }

func NewAlphaVantage(client HTTPClient, cfg Config) *AlphaVantage {
	return &AlphaVantage{NewProvider(client, cfg)}
}

func (p *AlphaVantage) BuildRequest(ctx context.Context, pair rate.CurrencyPair) (*http.Request, error) {
	return newRequest(ctx, p.Provider.cfg.Endpoint,
		web.WithValue("from_currency", pair.Base),
		web.WithValue("to_currency", pair.Quote),
		web.WithValue(p.cfg.Header, p.cfg.Key),
	)
}

func (p *AlphaVantage) ProcessResponse(resp *http.Response) (float64, error) {
	var data struct {
		RealtimeCurrencyExchangeRate struct {
			FromCurrencyCode string `json:"1. From_Currency Code"`
			FromCurrencyName string `json:"2. From_Currency Name"`
			ToCurrencyCode   string `json:"3. To_Currency Code"`
			ToCurrencyName   string `json:"4. To_Currency Name"`
			ExchangeRate     string `json:"5. Exchange Rate"`
			LastRefreshed    string `json:"6. Last Refreshed"`
		} `json:"Realtime Currency Exchange Rate"`
	}

	if err := web.ProcessResponse(resp, &data); err != nil {
		return 0, err
	}

	val, err := strconv.ParseFloat(data.RealtimeCurrencyExchangeRate.ExchangeRate, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing rate: %w", err)
	}

	return val, nil
}
