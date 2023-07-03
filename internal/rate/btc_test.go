package rate_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/stretchr/testify/require"
)

type HTTPClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (h *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	return h.DoFunc(req)
}

func TestGetBTCExchangeRate(t *testing.T) {
	tests := map[string]struct {
		mockDoFunc func(*http.Request) (*http.Response, error)
		wantRate   float64
		wantErr    error
	}{
		"Valid rate": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"rates":{"uah":{"value":2.5}}}`)),
				}, nil
			},
			wantRate: 2.5,
			wantErr:  nil,
		},
		"Rate retrieval failure": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			},
			wantRate: 0,
			wantErr:  errors.New("failed to retrieve rate"),
		},
		"Error sending request": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantRate: 0,
			wantErr:  errors.New("sending request: network error"),
		},
		"Error decoding response": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"rates":{"uah":{"value":"invalid"}}}`)),
				}, nil
			},
			wantRate: 0,
			wantErr:  errors.New("decoding response: invalid character 'i' looking for beginning of value"),
		},
		"Missing currency in response": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"rates":{}}`)),
				}, nil
			},
			wantErr: errors.New("currency not found: uah"),
		},
		"Unexpected status code": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			},
			wantRate: 0,
			wantErr:  errors.New("status code: 403"),
		},
		"Context timeout": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				time.Sleep(2 * time.Second)
				return nil, context.DeadlineExceeded
			},
			wantErr: context.DeadlineExceeded,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			client := &HTTPClientMock{DoFunc: tt.mockDoFunc}
			svc := rate.NewBTCExchangeRateClient(client, "https://api.coingecko.com/api/v3/exchange_rates")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			gotRate, err := svc.GetBTCExchangeRate(ctx, "uah")
			if tt.wantErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantRate, gotRate)
		})
	}
}
