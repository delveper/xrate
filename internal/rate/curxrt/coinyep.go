package curxrt

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// CoinYep API call was sniffed from https://coinyep.com
type CoinYep struct{}

func (p CoinYep) BuildRequest(ctx context.Context, pair rate.CurrencyPair, cfg Config) (*http.Request, error) {
	return newRequest(ctx, cfg.Endpoint,
		web.WithValue("from", pair.Base),
		web.WithValue("to", pair.Quote),
		web.WithValue("lang", "en"),
		web.WithValue("format", "json"),
	)
}

func (p CoinYep) ProcessResponse(resp *http.Response) (float64, error) {
	var data struct {
		BaseSymbol   string  `json:"base_symbol"`
		BaseName     string  `json:"base_name"`
		TargetSymbol string  `json:"target_symbol"`
		TargetName   string  `json:"target_name"`
		Price        string  `json:"price"`
		PriceChange  float64 `json:"price_change"`
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
