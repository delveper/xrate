package curxrate

import (
	"context"
	"net/http"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// Ninjas https://api-ninjas.com/api/exchangerate
type Ninjas struct{ Provider }

func NewNinjas(client HTTPClient, cfg Config) *Ninjas {
	return &Ninjas{NewProvider(client, cfg)}
}

func (p *Ninjas) BuildRequest(ctx context.Context, pair rate.CurrencyPair) (*http.Request, error) {
	return newRequest(ctx, p.Provider.cfg.Endpoint,
		web.WithHeader(p.cfg.Header, p.cfg.Key),
		web.WithValue("pair", pair.Base+"_"+pair.Quote),
	)
}

func (p *Ninjas) ProcessResponse(resp *http.Response) (float64, error) {
	var data struct {
		CurrencyPair string  `json:"currency_pair"`
		ExchangeRate float64 `json:"exchange_rate"`
	}

	if err := web.ProcessResponse(resp, &data); err != nil {
		return 0, err
	}

	return data.ExchangeRate, nil
}
