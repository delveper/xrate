package subs_test

import (
	"context"
	"net/mail"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subs"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/GenesisEducationKyiv/main-project-delveper/test/mock"
	"github.com/stretchr/testify/require"
)

func TestServiceRespondSubscription(t *testing.T) {
	tests := map[string]struct {
		event    event.Event
		repo     subs.SubscriberRepository
		wantSubs []subs.Subscription
		wantErr  error
	}{
		"valid_event_with_subscriptions": {
			event: event.New(subs.EventSource, subs.EventKindRequested, &mock.CurrencyPairEventMock{
				BaseCurrencyFunc:  func() string { return "USD" },
				QuoteCurrencyFunc: func() string { return "EUR" },
			}),
			repo: &mock.SubscriberRepositoryMock{
				ListFunc: func(ctx context.Context) ([]subs.Subscription, error) {
					return []subs.Subscription{
						{Subscriber: subs.Subscriber{Address: &mail.Address{Name: "user_0", Address: "user_0@example.com"}}, Topic: subs.NewTopic("USD", "EUR")},
						{Subscriber: subs.Subscriber{Address: &mail.Address{Name: "user_1", Address: "user_1@example.com"}}, Topic: subs.NewTopic("USD", "EUR")},
					}, nil
				},
			},
			wantSubs: []subs.Subscription{
				{Subscriber: subs.Subscriber{Address: &mail.Address{Name: "user_0", Address: "user_0@example.com"}}, Topic: subs.NewTopic("USD", "EUR")},
				{Subscriber: subs.Subscriber{Address: &mail.Address{Name: "user_1", Address: "user_1@example.com"}}, Topic: subs.NewTopic("USD", "EUR")},
			},
			wantErr: nil,
		},

		"invalid_payload_type": {
			event: event.New(subs.EventSource, subs.EventKindRequested, "invalid_payload"),
			repo: &mock.SubscriberRepositoryMock{
				ListFunc: func(ctx context.Context) ([]subs.Subscription, error) {
					return []subs.Subscription{
						{Subscriber: subs.Subscriber{Address: &mail.Address{Name: "user_0", Address: "user_0@example.com"}}, Topic: subs.NewTopic("USD", "EUR")},
					}, nil
				},
			},
			wantSubs: nil,
			wantErr:  subs.ErrInvalidEvent,
		},

		"error_from_list": {
			event: event.New(subs.EventSource, subs.EventKindRequested, &mock.CurrencyPairEventMock{
				BaseCurrencyFunc:  func() string { return "USD" },
				QuoteCurrencyFunc: func() string { return "EUR" },
			}),
			repo: &mock.SubscriberRepositoryMock{
				ListFunc: func(ctx context.Context) ([]subs.Subscription, error) {
					return nil, subs.ErrNotFound
				},
			},
			wantSubs: nil,
			wantErr:  subs.ErrNotFound,
		},

		"no_response_channel": {
			event: event.New(subs.EventSource, subs.EventKindRequested, &mock.CurrencyPairEventMock{
				BaseCurrencyFunc:  func() string { return "USD" },
				QuoteCurrencyFunc: func() string { return "EUR" },
			}),
			repo: &mock.SubscriberRepositoryMock{
				ListFunc: func(ctx context.Context) ([]subs.Subscription, error) {
					return []subs.Subscription{
						{Subscriber: subs.Subscriber{Address: &mail.Address{Name: "user_0", Address: "user_0@example.com"}}, Topic: subs.NewTopic("USD", "EUR")},
					}, nil
				},
			},
			wantSubs: []subs.Subscription{
				{Subscriber: subs.Subscriber{Address: &mail.Address{Name: "user_0", Address: "user_0@example.com"}}, Topic: subs.NewTopic("USD", "EUR")},
			},
			wantErr: nil,
		},
	}
	log := logger.New(logger.WithConsoleCore(logger.LevelDebug))
	bus := event.NewBus(log)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := subs.NewService(bus, tc.repo)

			err := svc.RespondSubscription(context.Background(), tc.event)
			require.ErrorIs(t, err, tc.wantErr)

			if tc.event.Response != nil {
				select {
				case resp := <-tc.event.Response:
					subscriptions, ok := resp.Payload.(subs.Subscriptions)
					require.True(t, ok)
					require.EqualValues(t, tc.wantSubs, []subs.Subscription(subscriptions))
				default:
					require.ErrorIs(t, err, tc.wantErr)
				}
			}
		})
	}
}
