package adpt

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

var (
	ErrUnexpected = errors.New("unexpected error")
	ErrNotFound   = errors.New("currency not found")
)

type Adapter struct {
	client   HTTPClient
	endpoint string
}

type Config struct {
	Endpoint string
	Key      string
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func NewAdapter(client HTTPClient, endpoint string) *Adapter {
	return &Adapter{client: client, endpoint: endpoint}
}

func (a *Adapter) SendRequest(ctx context.Context, opts ...web.RequestOption) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	web.ApplyRequestOptions(req, opts...)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}
