package subscription

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/stretchr/testify/require"
)

func TestHandlerSubscribe(t *testing.T) {
	testCases := map[string]struct {
		email      string
		subscriber Subscriber
		wantCode   int
		want       Response
	}{
		"Success": {
			email:      "email@example.com",
			subscriber: &SubscriberMock{SubscribeFunc: func(Email) error { return nil }},
			wantCode:   http.StatusOK,
			want:       Response{Message: StatusSubscribed},
		},
		"Email already exists": {
			email: "exists@example.com",
			subscriber: &SubscriberMock{SubscribeFunc: func(Email) error {
				return ErrEmailAlreadyExists
			}},
			wantCode: http.StatusConflict,
			want:     Response{Message: StatusError, Details: ErrEmailAlreadyExists.Error()},
		},
	}

	log := logger.New(logger.LevelDebug)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/subscribe", nil)
			require.NoError(t, err)

			q := req.URL.Query()
			q.Add("email", tc.email)
			req.URL.RawQuery = q.Encode()

			h := NewHandler(tc.subscriber, log)

			rw := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Subscribe)

			handler.ServeHTTP(rw, req)
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
	testCases := map[string]struct {
		subscriber Subscriber
		wantCode   int
		want       Response
	}{
		"Successful email sending": {
			subscriber: &SubscriberMock{SendEmailsFunc: func() error { return nil }},
			wantCode:   http.StatusOK,
			want:       Response{Message: StatusSend},
		},
		// TODO(not_documented): Add case for internal server error.
	}

	log := logger.New(logger.LevelDebug)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/sendEmails", nil)
			require.NoError(t, err)

			h := NewHandler(tc.subscriber, log)

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
