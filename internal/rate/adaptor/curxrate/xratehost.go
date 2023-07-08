package curxrate

import (
	"context"
	"net/http"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// ExchangeRateHost https://rapidapi.com/Serply/api/exchange-rate9
type ExchangeRateHost struct{ Provider }

func NewExchangeRateHost(client HTTPClient, cfg Config) *ExchangeRateHost {
	return &ExchangeRateHost{NewProvider(client, cfg)}
}

func (p *ExchangeRateHost) BuildRequest(ctx context.Context, pair rate.CurrencyPair) (*http.Request, error) {
	return newRequest(ctx, p.cfg.Endpoint,
		web.WithValue("from", pair.Base),
		web.WithValue("to", pair.Quote),
		web.WithHeader(p.cfg.Header, p.cfg.Key),
	)
}

func (p *ExchangeRateHost) ProcessResponse(resp *http.Response) (float64, error) {
	var data struct {
		Motd struct {
			Msg string `json:"msg"`
			URL string `json:"url"`
		} `json:"motd"`
		Success    bool               `json:"success"`
		Historical bool               `json:"historical"`
		Base       string             `json:"base"`
		Date       string             `json:"date"`
		Rates      map[string]float64 `json:"rates"`
	}

	if err := web.ProcessResponse(resp, &data); err != nil {
		return 0, err
	}

	var val float64
	for _, v := range data.Rates {
		val = v
		break
	}

	return val, nil
}
