package rate

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// Clienter represents the behavior of http.Client needed by Service.
type Clienter interface {
	Do(req *http.Request) (*http.Response, error)
}

// MockClient is a mock implementation of Clienter.
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestService_Get(t *testing.T) {
	tests := map[string]struct {
		mockDoFunc func(req *http.Request) (*http.Response, error)
		wantRate   float64
		wantErr    bool
	}{
		"valid rate": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"rates":{"UAH":{"value":2.5}}}`)),
				}, nil
			},
			wantRate: 2.5,
			wantErr:  false,
		},
		"rate retrieval failure": {
			mockDoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			},
			wantRate: 0,
			wantErr:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := &Service{Client: &http.Client{}}
			s.Client.Transport = roundTripFunc(tt.mockDoFunc)

			gotRate, err := s.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRate != tt.wantRate {
				t.Errorf("Service.Get() = %v, want %v", gotRate, tt.wantRate)
			}
		})
	}
}
