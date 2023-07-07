package rate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate/adpt"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

var (
	ErrUnexpected = errors.New("unexpected error")
	ErrNotFound   = errors.New("currency not found")
)

type responseBTCExchangeRate struct {
	Rates map[string]struct {
		Value float64 `json:"value"`
	} `json:"rates"`
}

type BTCExchangeRateClient struct {
	client   adpt.HTTPClient
	endpoint string
}

func NewBTCExchangeRateClient(client adpt.HTTPClient, endpoint string) *BTCExchangeRateClient {
	return &BTCExchangeRateClient{
		client:   client,
		endpoint: endpoint,
	}
}
func (a *BTCExchangeRateClient) GetBTCExchangeRate(ctx context.Context, currency string) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("executing btc exchange request: %w", err)
	}

	defer resp.Body.Close()

	if err := web.ErrFromStatusCode(resp.StatusCode); err != nil {
		return 0, err
	}

	var data responseBTCExchangeRate
	if err := web.DecodeBody(resp.Body, &data); err != nil {
		return 0, err
	}

	rate, ok := data.Rates[strings.ToLower(currency)]
	if !ok {
		return 0, ErrNotFound
	}

	return rate.Value, nil
}
