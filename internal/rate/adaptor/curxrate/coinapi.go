package curxrate

import (
	"context"
	"net/http"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// CoinApi https://docs.coinapi.io/market-data/rest-api/exchange-rates
type CoinApi struct{ Provider }

func NewCoinApi(client HTTPClient, cfg Config) *CoinApi {
	return &CoinApi{NewProvider(client, cfg)}
}

func (p *CoinApi) BuildRequest(ctx context.Context, pair rate.CurrencyPair) (*http.Request, error) {
	return newRequest(ctx, p.cfg.Endpoint,
		web.WithPath(pair.Base, pair.Quote),
		web.WithHeader(p.cfg.Header, p.cfg.Key),
	)
}

func (p *CoinApi) ProcessResponse(resp *http.Response) (float64, error) {
	var data struct {
		Time         time.Time `json:"time"`
		AssetIdBase  string    `json:"asset_id_base"`
		AssetIdQuote string    `json:"asset_id_quote"`
		Rate         float64   `json:"rate"`
	}

	if err := web.ProcessResponse(resp, &data); err != nil {
		return 0, err
	}

	return data.Rate, nil
}
