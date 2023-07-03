package subscription_test

import (
	"context"
	"errors"
	"net/mail"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription"
	"github.com/GenesisEducationKyiv/main-project-delveper/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestServiceSubscribe(t *testing.T) {
	tests := map[string]struct {
		email       subscription.Email
		repo        subscription.EmailRepository
		rateGetter  subscription.RateGetter
		emailSender subscription.EmailSender
		wantErr     error
	}{
		"Successful subscription": {
			email:       subscription.Email{Address: &mail.Address{Address: "test@example.com"}},
			repo:        &mock.EmailRepositoryMock{AddFunc: func(subscription.Email) error { return nil }},
			rateGetter:  &mock.RateGetterMock{},
			emailSender: &mock.EmailSenderMock{},
			wantErr:     nil,
		},
		"Failed subscription due to existing email": {
			email:       subscription.Email{Address: &mail.Address{Address: "test@example.com"}},
			repo:        &mock.EmailRepositoryMock{AddFunc: func(subscription.Email) error { return subscription.ErrEmailAlreadyExists }},
			rateGetter:  &mock.RateGetterMock{},
			emailSender: &mock.EmailSenderMock{},
			wantErr:     subscription.ErrEmailAlreadyExists,
		},
		// TODO(): Failed subscription due to internal server error.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := subscription.NewService(tc.repo, tc.rateGetter, tc.emailSender)

			err := svc.Subscribe(tc.email)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestServiceSendEmails(t *testing.T) {
	tests := map[string]struct {
		emails      []subscription.Email
		rate        float64
		repo        subscription.EmailRepository
		rateGetter  subscription.RateGetter
		emailSender subscription.EmailSender
		wantErr     error
	}{
		"Successful email sending": {
			emails: []subscription.Email{
				{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}},
			},
			rate: 1.0,
			repo: &mock.EmailRepositoryMock{
				GetAllFunc: func() ([]subscription.Email, error) {
					return []subscription.Email{
						{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}},
					}, nil
				},
			},
			rateGetter:  &mock.RateGetterMock{GetFunc: func(context.Context) (float64, error) { return 1.0, nil }},
			emailSender: &mock.EmailSenderMock{SendFunc: func(subscription.Email, float64) error { return nil }},
			wantErr:     nil,
		},
		"Failed to get rate": {
			emails: []subscription.Email{
				{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}},
			},
			rate: 1.0,
			repo: &mock.EmailRepositoryMock{},
			rateGetter: &mock.RateGetterMock{
				GetFunc: func(context.Context) (float64, error) { return 0, errors.New("getting rate: failed to get rate") },
			},
			emailSender: &mock.EmailSenderMock{},
			wantErr:     errors.New("getting rate: failed to get rate"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := subscription.NewService(tc.repo, tc.rateGetter, tc.emailSender)

			err := svc.SendEmails()
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
