package rate_test

import (
	"context"
	"errors"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/GenesisEducationKyiv/main-project-delveper/test/mock"
	"github.com/stretchr/testify/require"
)

func TestService_ResponseExchangeRate(t *testing.T) {
	tests := map[string]struct {
		event        event.Event
		mockProvider *mock.ExchangeRateProviderMock
		wantErr      error
	}{
		"valid event": {
			event: event.Event{
				Source:   rate.EventSource,
				Kind:     rate.EventKindRequested,
				Payload:  rate.NewCurrencyPair("USD", "EUR"),
				Response: make(chan event.Event, 1),
			},
			mockProvider: &mock.ExchangeRateProviderMock{
				GetExchangeRateFunc: func(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
					return &rate.ExchangeRate{
						Pair:  rate.NewCurrencyPair("USD", "EUR"),
						Value: 1.2,
					}, nil
				},
				StringFunc: func() string {
					return "USD/EUR"
				},
			},
			wantErr: nil,
		},

		"invalid payload type": {
			event: event.Event{
				Source:   rate.EventSource,
				Kind:     rate.EventKindRequested,
				Payload:  "USD/EUR",
				Response: make(chan event.Event, 1),
			},
			mockProvider: &mock.ExchangeRateProviderMock{
				GetExchangeRateFunc: func(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
					return nil, errors.New("invalid event: unexpected payload, expected CurrencyPairEvent: string")
				},
				StringFunc: func() string {
					return "USD/EUR"
				},
			},
			wantErr: errors.New("invalid event: unexpected payload, expected CurrencyPairEvent: string"),
		},

		"no response channel": {
			event: event.Event{
				Source:  rate.EventSource,
				Kind:    rate.EventKindRequested,
				Payload: rate.NewCurrencyPair("USD", "EUR"),
			},
			mockProvider: &mock.ExchangeRateProviderMock{
				GetExchangeRateFunc: func(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
					return &rate.ExchangeRate{
						Pair:  rate.NewCurrencyPair("USD", "EUR"),
						Value: 1.2,
					}, nil
				},
				StringFunc: func() string {
					return "USD/EUR"
				},
			},
			wantErr: errors.New("responding exchange rate event: response channel is not initialized"),
		},

		"error from GetExchangeRate": {
			event: event.Event{
				Source:   rate.EventSource,
				Kind:     rate.EventKindRequested,
				Payload:  rate.NewCurrencyPair("USD", "EUR"),
				Response: make(chan event.Event, 1),
			},
			mockProvider: &mock.ExchangeRateProviderMock{
				GetExchangeRateFunc: func(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
					return nil, errors.New("mock error")
				},
				StringFunc: func() string {
					return "USD/EUR"
				},
			},
			wantErr: errors.New("mock error"),
		},
	}

	log := logger.New(logger.WithConsoleCore(logger.LevelDebug))
	bus := event.NewBus(log)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := rate.NewService(bus, tc.mockProvider)

			err := svc.ResponseExchangeRate(context.Background(), tc.event)
			if tc.wantErr != nil {
				require.Error(t, err)
				return
			}

			require.Nil(t, err)

			select {
			case resp := <-tc.event.Response:
				r, ok := resp.Payload.(*rate.ExchangeRate)
				require.True(t, ok)
				require.Equal(t, tc.event.Payload, r.Pair)
			default:
				t.Fatal("no response event")
			}
		})
	}
}
