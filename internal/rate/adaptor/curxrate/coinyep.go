package curxrate

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// CoinYep API call was sniffed from https://coinyep.com
type CoinYep struct{ Provider }

func NewCoinYep(client HTTPClient, cfg Config) *CoinYep {
	return &CoinYep{NewProvider(client, cfg)}
}

func (p *CoinYep) BuildRequest(ctx context.Context, pair rate.CurrencyPair) (*http.Request, error) {
	return newRequest(ctx, p.Provider.cfg.Endpoint,
		web.WithValue("from", pair.Base),
		web.WithValue("to", pair.Quote),
		web.WithValue("lang", "en"),
		web.WithValue("format", "json"),
	)
}

func (p *CoinYep) ProcessResponse(resp *http.Response) (float64, error) {
	var data struct {
		BaseSymbol   string `json:"base_symbol"`
		TargetSymbol string `json:"target_symbol"`
		Price        string `json:"price"`
	}

	if err := web.ProcessResponse(resp, &data); err != nil {
		return 0, err
	}

	val, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing rate: %w", err)
	}

	return val, nil
}
