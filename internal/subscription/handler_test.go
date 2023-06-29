package subscription_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription/mocks"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/stretchr/testify/require"
)

func TestHandlerSubscribe(t *testing.T) {
	testCases := map[string]struct {
		email      string
		subscriber subscription.Subscriber
		wantCode   int
		want       subscription.Response
	}{
		"Success": {
			email:      "email@example.com",
			subscriber: &mocks.SubscriberMock{SubscribeFunc: func(subscription.Email) error { return nil }},
			wantCode:   http.StatusOK,
			want:       subscription.Response{Message: subscription.StatusSubscribed},
		},
		"Email already exists": {
			email: "exists@example.com",
			subscriber: &mocks.SubscriberMock{SubscribeFunc: func(subscription.Email) error {
				return subscription.ErrEmailAlreadyExists
			}},
			wantCode: http.StatusConflict,
			want:     subscription.Response{Message: subscription.StatusError, Details: subscription.ErrEmailAlreadyExists.Error()},
		},
	}

	log := logger.New(logger.LevelDebug)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rw := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/api/subscribe", nil)
			require.NoError(t, err)

			q := req.URL.Query()
			q.Add("email", tc.email)
			req.URL.RawQuery = q.Encode()

			h := http.HandlerFunc(subscription.NewHandler(tc.subscriber, log).Subscribe)
			h.ServeHTTP(rw, req)
			require.Equal(t, tc.wantCode, rw.Code)

			gotJSON, err := io.ReadAll(rw.Body)
			require.NoError(t, err)

			wantJSON, err := json.Marshal(tc.want)
			require.NoError(t, err)

			log.Infof("%s", gotJSON)
			require.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestHandlerSendEmails(t *testing.T) {
	tests := map[string]struct {
		subscriber subscription.Subscriber
		wantCode   int
		want       subscription.Response
	}{
		"Successful email sending": {
			subscriber: &mocks.SubscriberMock{SendEmailsFunc: func() error { return nil }},
			wantCode:   http.StatusOK,
			want:       subscription.Response{Message: subscription.StatusSend},
		},
		// TODO(not_documented): Add case for internal server error.
	}

	log := logger.New(logger.LevelDebug)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/sendEmails", nil)
			require.NoError(t, err)

			h := subscription.NewHandler(tc.subscriber, log)

			rw := httptest.NewRecorder()
			handler := http.HandlerFunc(h.SendEmails)

			handler.ServeHTTP(rw, req)
			require.Equal(t, tc.wantCode, rw.Code)

			gotJSON, err := io.ReadAll(rw.Body)
			require.NoError(t, err)

			wantJSON, err := json.Marshal(tc.want)
			require.NoError(t, err)

			require.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}
