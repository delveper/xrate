package curxrate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

type Config struct {
	Endpoint string
	Header   string
	Key      string
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type RequestBuilder interface {
	BuildRequest(context.Context, rate.CurrencyPair) (*http.Request, error)
}

type ResponseProcessor interface {
	ProcessResponse(*http.Response) (float64, error)
}

// Provider implements rate.ExchangeRateProvider.
type Provider struct {
	// RequestBuilder and ResponseProcessor
	// will be overridden by struct that embeds Provider.
	RequestBuilder
	ResponseProcessor

	cfg Config
	clt HTTPClient
}

func NewProvider(client HTTPClient, cfg Config) Provider {
	return Provider{
		clt: client,
		cfg: cfg,
	}
}

func (p *Provider) GetExchangeRate(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
	req, err := p.BuildRequest(ctx, pair)
	if err != nil {
		return nil, err
	}

	resp, err := p.clt.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	val, err := p.ProcessResponse(resp)
	if err != nil {
		return nil, err
	}

	return rate.NewExchangeRate(val, pair), nil
}

func newRequest(ctx context.Context, endpoint string, opts ...func(r *http.Request)) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	web.ApplyRequestOptions(req, opts...)

	return req, nil
}
