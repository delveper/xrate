package curxrt

import (
	"context"
	"net/http"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// Ninjas https://api-ninjas.com/api/exchangerate
type Ninjas struct{}

func (p Ninjas) BuildRequest(ctx context.Context, pair rate.CurrencyPair, cfg Config) (*http.Request, error) {
	return newRequest(ctx, cfg.Endpoint,
		web.WithHeader(cfg.Header, cfg.Key),
		web.WithValue("pair", pair.Base+"_"+pair.Quote),
	)
}

func (p Ninjas) ProcessResponse(resp *http.Response) (float64, error) {
	var data struct {
		CurrencyPair string  `json:"currency_pair"`
		ExchangeRate float64 `json:"exchange_rate"`
	}

	if err := web.ProcessResponse(resp, &data); err != nil {
		return 0, err
	}

	return data.ExchangeRate, nil
}
