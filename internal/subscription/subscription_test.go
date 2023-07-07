package subscription_test

import (
	"context"
	"errors"
	"net/mail"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription"
	"github.com/GenesisEducationKyiv/main-project-delveper/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestServiceSubscribe(t *testing.T) {
	tests := map[string]struct {
		email       subscription.Subscriber
		repo        subscription.SubscriberRepository
		rateGetter  rate.ExchangeRateService
		emailSender subscription.EmailSender
		wantErr     error
	}{
		"Successful subscription": {
			email:       subscription.Subscriber{Address: &mail.Address{Address: "test@example.com"}},
			repo:        &mock.SubscriberRepositoryMock{AddFunc: func(context.Context, subscription.Subscriber) error { return nil }},
			rateGetter:  &mock.ExchangeRateServiceMock{},
			emailSender: &mock.EmailSenderMock{},
			wantErr:     nil,
		},
		"Failed subscription due to existing email": {
			email: subscription.Subscriber{Address: &mail.Address{Address: "test@example.com"}},
			repo: &mock.SubscriberRepositoryMock{AddFunc: func(context.Context, subscription.Subscriber) error {
				return subscription.ErrEmailAlreadyExists
			}},
			rateGetter:  &mock.ExchangeRateServiceMock{},
			emailSender: &mock.EmailSenderMock{},
			wantErr:     subscription.ErrEmailAlreadyExists,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := subscription.NewService(tc.repo, tc.rateGetter, tc.emailSender)

			err := svc.Subscribe(context.Background(), tc.email)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestServiceSendEmails(t *testing.T) {
	tests := map[string]struct {
		emails      []subscription.Subscriber
		rate        float64
		repo        subscription.SubscriberRepository
		rateGetter  rate.ExchangeRateService
		emailSender subscription.EmailSender
		wantErr     error
	}{
		"Successful email sending": {
			emails: []subscription.Subscriber{
				{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}},
			},
			rate: 1.0,
			repo: &mock.SubscriberRepositoryMock{
				ListFunc: func(context.Context) ([]subscription.Subscriber, error) {
					return []subscription.Subscriber{
						{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}},
					}, nil
				},
			},
			rateGetter: &mock.ExchangeRateServiceMock{GetFunc: func(context.Context, rate.CurrencyPair) (*rate.ExchangeRate, error) {
				return &rate.ExchangeRate{Value: 1.0}, nil
			}},
			emailSender: &mock.EmailSenderMock{SendFunc: func(message subscription.Message) error { return nil }},
			wantErr:     nil,
		},
		"Failed to get rate": {
			emails: []subscription.Subscriber{
				{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}},
			},
			rate: 1.0,
			repo: &mock.SubscriberRepositoryMock{},
			rateGetter: &mock.ExchangeRateServiceMock{
				GetFunc: func(context.Context, rate.CurrencyPair) (*rate.ExchangeRate, error) {
					return nil, errors.New("getting rate: failed to get rate")
				},
			},
			emailSender: &mock.EmailSenderMock{},
			wantErr:     errors.New("getting rate: failed to get rate"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := subscription.NewService(tc.repo, tc.rateGetter, tc.emailSender)

			err := svc.SendEmails(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
