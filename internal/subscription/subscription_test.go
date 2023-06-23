package subscription

import (
	"fmt"
	"net/mail"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=email_repository_mock_test.go -package=subscription github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription EmailRepository
//go:generate mockgen -destination=rate_getter_mock_test.go -package=subscription github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription RateGetter
//go:generate mockgen -destination=email_sender_mock_test.go -package=subscription github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription EmailSender

func TestServiceSubscribe(t *testing.T) {
	tests := map[string]struct {
		email              Email
		repoBuilder        func(*gomock.Controller) EmailRepository
		rateGetterBuilder  func(*gomock.Controller) RateGetter
		emailSenderBuilder func(*gomock.Controller) EmailSender
		wantErr            error
	}{
		"Successful subscription": {
			email: Email{Address: &mail.Address{Address: "test@example.com"}},
			repoBuilder: func(ctrl *gomock.Controller) EmailRepository {
				mock := NewMockEmailRepository(ctrl)
				mock.EXPECT().
					Add(gomock.Any()).
					Return(nil).
					Times(1)
				return mock
			},
			rateGetterBuilder:  func(ctrl *gomock.Controller) RateGetter { return NewMockRateGetter(ctrl) },
			emailSenderBuilder: func(ctrl *gomock.Controller) EmailSender { return NewMockEmailSender(ctrl) },
			wantErr:            nil,
		},
		"Failed subscription due to existing email": {
			email: Email{Address: &mail.Address{Address: "test@example.com"}},
			repoBuilder: func(ctrl *gomock.Controller) EmailRepository {
				mock := NewMockEmailRepository(ctrl)
				mock.EXPECT().
					Add(gomock.Any()).
					Return(ErrEmailAlreadyExists).
					Times(1)
				return mock
			},
			rateGetterBuilder:  func(ctrl *gomock.Controller) RateGetter { return NewMockRateGetter(ctrl) },
			emailSenderBuilder: func(ctrl *gomock.Controller) EmailSender { return NewMockEmailSender(ctrl) },
			wantErr:            fmt.Errorf("adding email subscription: %w", ErrEmailAlreadyExists),
		},
		// TODO(): Failed subscription due to internal server error.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.repoBuilder(ctrl)
			rateGetter := tc.rateGetterBuilder(ctrl)
			emailSender := tc.emailSenderBuilder(ctrl)

			svc := NewService(repo, rateGetter, emailSender)

			err := svc.Subscribe(tc.email)

			if tc.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestServiceSendEmails(t *testing.T) {
	testCases := map[string]struct {
		emails             []Email
		rate               float64
		repoBuilder        func(*gomock.Controller) EmailRepository
		rateGetterBuilder  func(*gomock.Controller) RateGetter
		emailSenderBuilder func(*gomock.Controller) EmailSender
		wantErr            error
	}{
		"Successful email sending": {
			emails: []Email{{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}}},
			rate:   1.0,
			repoBuilder: func(ctrl *gomock.Controller) EmailRepository {
				mock := NewMockEmailRepository(ctrl)
				mock.EXPECT().
					GetAll().
					Return([]Email{{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}}}, nil).
					Times(1)
				return mock
			},
			rateGetterBuilder: func(ctrl *gomock.Controller) RateGetter {
				mock := NewMockRateGetter(ctrl)
				mock.EXPECT().
					Get().
					Return(1.0, nil).
					Times(1)
				return mock
			},
			emailSenderBuilder: func(ctrl *gomock.Controller) EmailSender {
				mock := NewMockEmailSender(ctrl)
				mock.EXPECT().
					Send(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(2)
				return mock
			},
			wantErr: nil,
		},
		"Failed to get rate": {
			emails: []Email{{Address: &mail.Address{Address: "test1@example.com"}}, {Address: &mail.Address{Address: "test2@example.com"}}},
			rate:   1.0,
			repoBuilder: func(ctrl *gomock.Controller) EmailRepository {
				mock := NewMockEmailRepository(ctrl)
				mock.EXPECT().
					GetAll().
					Times(0)
				return mock
			},
			rateGetterBuilder: func(ctrl *gomock.Controller) RateGetter {
				mock := NewMockRateGetter(ctrl)
				mock.EXPECT().
					Get().
					Return(0.0, fmt.Errorf("failed to get rate")).
					Times(1)
				return mock
			},
			emailSenderBuilder: func(ctrl *gomock.Controller) EmailSender {
				mock := NewMockEmailSender(ctrl)
				mock.EXPECT().
					Send(gomock.Any(), gomock.Any()).
					Times(0)
				return mock
			},
			wantErr: fmt.Errorf("getting rate: failed to get rate"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.repoBuilder(ctrl)
			rateGetter := tc.rateGetterBuilder(ctrl)
			emailSender := tc.emailSenderBuilder(ctrl)

			svc := NewService(repo, rateGetter, emailSender)

			err := svc.SendEmails()

			if tc.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
