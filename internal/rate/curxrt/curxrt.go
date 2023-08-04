/*
Package curxrt provides functionality for retrieving currency exchange rates via HTTP.
*/
package curxrt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// Config defines configuration for the exchange rate provider.
type Config struct {
	Name     string
	Endpoint string
	Header   string
	Key      string
}

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// RequestBuilder is an interface for building HTTP requests for retrieving exchange rates.
type RequestBuilder interface {
	BuildRequest(context.Context, rate.CurrencyPair, Config) (*http.Request, error)
}

// ResponseProcessor is an interface for processing HTTP responses from the exchange rate provider.
type ResponseProcessor interface {
	ProcessResponse(*http.Response) (float64, error)
}

// RequestResponder designed to the build and process HTTP requests of a specific provider.
type RequestResponder interface {
	RequestBuilder
	ResponseProcessor
}

// Provider implements rate.ExchangeRateProvider interface.
// It provides functionality to get exchange rates via HTTP.
type Provider[T any] struct {
	Config
	HTTPClient
	RequestResponder
}

// NewProvider returns a new instance of specific Provider with injected dependencies implemented by RequestResponder.
func NewProvider[T RequestResponder](cfg Config, clt HTTPClient) Provider[T] {
	return Provider[T]{Config: cfg, HTTPClient: clt, RequestResponder: *new(T)}
}

// String returns the name of the provider for log.
func (p Provider[T]) String() string { return p.Name }

// GetExchangeRate retrieves the exchange rate for the specified currency pair.
func (p Provider[T]) GetExchangeRate(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
	req, err := p.BuildRequest(ctx, pair, p.Config)
	if err != nil {
		return nil, err
	}

	resp, err := p.Do(req)
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

// newRequest creates a new HTTP request with the specified context, endpoint, and request options.
func newRequest(ctx context.Context, endpoint string, opts ...func(*http.Request)) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	web.ApplyRequestOptions(req, opts...)

	return req, nil
}
