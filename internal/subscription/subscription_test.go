package subscription

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServiceSubscribe(t *testing.T) {
	tests := map[string]struct {
		email           Email
		repoMock        EmailRepository
		rateGetterMock  RateGetter
		emailSenderMock EmailSender
		wantErr         error
	}{
		"Successful subscription": {
			email:           Email{Address: &mail.Address{Address: "test@example.com"}},
			repoMock:        &mocks.EmailRepositoryMock{AddFunc: func(Email) error { return nil }},
			rateGetterMock:  &mocks.RateGetterMock{},
			emailSenderMock: &mocks.EmailSenderMock{},
			wantErr:         nil,
		},
		"Failed subscription due to existing email": {
			email:           Email{Address: &mail.Address{Address: "test@example.com"}},
			repoMock:        &mocks.EmailRepositoryMock{AddFunc: func(Email) error { return ErrEmailAlreadyExists }},
			rateGetterMock:  &mocks.RateGetterMock{},
			emailSenderMock: &mocks.EmailSenderMock{},
			wantErr:         fmt.Errorf("adding email subscription: %w", ErrEmailAlreadyExists),
		},
		// TODO(): Failed subscription due to internal server error.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := NewService(tc.repoMock, tc.rateGetterMock, tc.emailSenderMock)

			err := svc.Subscribe(tc.email)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestServiceSendEmails(t *testing.T) {
	testCases := map[string]struct {
		emails      []Email
		rate        float64
		repo        EmailRepository
		rateGetter  RateGetter
		emailSender EmailSender
		wantErr     error
	}{
		"Successful email sending": {
			emails: []Email{{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}}},
			rate:   1.0,
			repo: &mocks.EmailRepositoryMock{
				GetAllFunc: func() ([]Email, error) {
					return []Email{{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}}}, nil
				},
			},
			rateGetter:  &mocks.RateGetterMock{GetFunc: func(context.Context) (float64, error) { return 1.0, nil }},
			emailSender: &mocks.EmailSenderMock{SendFunc: func(Email, float64) error { return nil }},
			wantErr:     nil,
		},
		"Failed to get rate": {
			emails: []Email{{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}}},
			rate:   1.0,
			repo:   &mocks.EmailRepositoryMock{},
			rateGetter: &mocks.RateGetterMock{
				GetFunc: func(context.Context) (float64, error) { return 0, errors.New("getting rate: failed to get rate") },
			},
			emailSender: &mocks.EmailSenderMock{},
			wantErr:     errors.New("getting rate: failed to get rate"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			svc := NewService(tc.repo, tc.rateGetter, tc.emailSender)

			err := svc.SendEmails()
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
