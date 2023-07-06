package rate

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type ExchangeRateProviderAdapter struct {
	client HTTPClient
	url    *url.URL
	header http.Header
}

type AdapterOption func(*ExchangeRateProviderAdapter) error

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func NewAdapter(client HTTPClient, opts ...AdapterOption) (*ExchangeRateProviderAdapter, error) {
	a := &ExchangeRateProviderAdapter{client: client}

	for i := range opts {
		if err := opts[i](a); err != nil {
			return nil, err
		}
	}

	return a, nil
}

func withURL(u string) AdapterOption {
	return func(a *ExchangeRateProviderAdapter) error {
		u, err := url.Parse(u)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}
		a.url = u
		return nil
	}
}

func withValue(key string, value interface{}) AdapterOption {
	return func(a *ExchangeRateProviderAdapter) error {
		values := a.url.Query()
		values.Add(key, fmt.Sprint(value))
		a.url.RawQuery = values.Encode()
		return nil
	}
}

func withPath(paths ...string) AdapterOption {
	return func(a *ExchangeRateProviderAdapter) error {
		a.url.Path = path.Join(paths...)
		return nil
	}
}

func withHeaders(pairs ...string) AdapterOption {
	return func(a *ExchangeRateProviderAdapter) error {
		if len(pairs)%2 != 0 {
			return fmt.Errorf("header pairs must contain an even number of elements")
		}
		headers := make(http.Header)
		for i := 0; i < len(pairs); i += 2 {
			key, val := pairs[i], pairs[i+1]
			headers.Add(key, val)
		}
		return nil
	}
}

func (a *ExchangeRateProviderAdapter) SendRequest(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header = a.header
	return a.client.Do(req)
}
