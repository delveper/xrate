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

func TestServiceLogExchangeRate(t *testing.T) {
	log := logger.New(logger.WithConsoleCore(logger.LevelDebug))
	bus := event.NewBus(log)
	svc := rate.NewService(bus, nil)

	tests := map[string]struct {
		event     event.Event
		wantError error
	}{
		"valid_event": {
			event: event.Event{
				Payload: rate.ProviderResponse{
					Provider:     "TestProvider",
					ExchangeRate: &rate.ExchangeRate{Pair: rate.NewCurrencyPair("USD", "EUR"), Value: 1.2},
				},
			},
			wantError: nil,
		},

		"provider_failed": {
			event: event.Event{
				Payload: rate.ProviderErrorResponse{
					Provider: "TestProvider",
					Err:      rate.ErrProviderUnavailable,
				},
			},
			wantError: nil,
		},

		"invalid_payload_type": {
			event: event.Event{
				Payload: "invalid_payload",
			},
			wantError: rate.ErrInvalidEvent,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := svc.LogExchangeRate(context.Background(), tc.event)
			require.ErrorIs(t, err, tc.wantError)
		})
	}
}

func TestServiceResponseExchangeRate(t *testing.T) {
	log := logger.New(logger.WithConsoleCore(logger.LevelDebug))
	bus := event.NewBus(log)

	tests := map[string]struct {
		event        event.Event
		mockProvider *mock.ExchangeRateProviderMock
		wantErr      error
	}{
		"valid_event": {
			event: event.New(rate.EventSource, rate.EventKindRequested, rate.NewCurrencyPair("USD", "EUR")),
			mockProvider: &mock.ExchangeRateProviderMock{
				GetExchangeRateFunc: func(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
					return &rate.ExchangeRate{Pair: rate.NewCurrencyPair("USD", "EUR"), Value: 1.2}, nil
				},
				StringFunc: func() string { return "USD/EUR" },
			},
			wantErr: nil,
		},

		"invalid_payload_type": {
			event:        event.New(rate.EventSource, rate.EventKindRequested, "invalid_payload"),
			mockProvider: nil,
			wantErr:      rate.ErrInvalidEvent,
		},

		"no_response_channel": {
			event: event.Event{
				Source:   rate.EventSource,
				Kind:     rate.EventKindFailed,
				Payload:  rate.NewCurrencyPair("USD", "EUR"),
				Response: nil,
			},
			mockProvider: &mock.ExchangeRateProviderMock{
				GetExchangeRateFunc: func(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
					return &rate.ExchangeRate{Pair: rate.NewCurrencyPair("USD", "EUR"), Value: 1.2}, nil
				},
				StringFunc: func() string { return "USD/EUR" },
			},
			wantErr: rate.ErrInvalidChannel,
		},

		"provider_failed": {
			event: event.New(rate.EventSource, rate.EventKindRequested, rate.NewCurrencyPair("USD", "EUR")),
			mockProvider: &mock.ExchangeRateProviderMock{
				GetExchangeRateFunc: func(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
					return nil, errors.New("mock error")
				},
				StringFunc: func() string { return "USD/EUR" },
			},
			wantErr: rate.ErrProviderUnavailable,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := rate.NewService(bus, tc.mockProvider)

			err := svc.ResponseExchangeRate(context.Background(), tc.event)
			require.ErrorIs(t, err, tc.wantErr)

			if tc.event.Response != nil {
				select {
				case resp := <-tc.event.Response:
					r, ok := resp.Payload.(*rate.ExchangeRate)
					require.True(t, ok)
					require.Equal(t, tc.event.Payload, r.Pair)
				default:
					require.ErrorIs(t, err, tc.wantErr)
				}
			}
		})
	}
}
