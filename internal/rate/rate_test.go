package rate

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceGet(t *testing.T) {
	tests := map[string]struct {
		mockDoFunc func(req *http.Request) (*http.Response, error)
		wantRate   float64
		wantErr    error
	}{
		"Valid rate": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"rates":{"UAH":{"value":2.5}}}`)),
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
			wantRate: 0.0,
			wantErr:  errors.New("failed to retrieve rate"),
		},
		"Error sending request": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantRate: 0.0,
			wantErr:  errors.New("sending request: network error"),
		},
		"Error decoding response": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"rates":{"UAH":{"value":"invalid"}}}`)),
				}, nil
			},
			wantRate: 0.0,
			wantErr:  errors.New("decoding response: invalid character 'i' looking for beginning of value"),
		},
		"Unexpected status code": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			},
			wantRate: 0.0,
			wantErr:  errors.New("status code: 403"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			svc := NewService("https://api.coingecko.com/api/v3/exchange_rates")
			svc.Client.Transport = roundTripFunc(tt.mockDoFunc)

			gotRate, err := svc.Get()
			if tt.wantErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantRate, gotRate)
		})
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
