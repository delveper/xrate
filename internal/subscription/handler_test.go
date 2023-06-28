package subscription

import (
	"encoding/json"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandlerSubscribe(t *testing.T) {
	cases := map[string]struct {
		email    string
		mockFunc func(m *SubscriberMock)
		wantCode int
		want     Response
	}{
		"Success": {
			email:    "email@example.com",
			mockFunc: func(m *SubscriberMock) { m.SubscribeFunc = func(Email) error { return nil } },
			wantCode: http.StatusOK,
			want:     Response{Message: StatusSubscribed},
		},
		"Email already exists": {
			email:    "exists@example.com",
			mockFunc: func(m *SubscriberMock) { m.SubscribeFunc = func(Email) error { return ErrEmailAlreadyExists } },
			wantCode: http.StatusConflict,
			want:     Response{Message: StatusError, Details: ErrEmailAlreadyExists.Error()},
		},
		// TODO(not_documented): add more cases.
	}

	log := logger.New(logger.LevelDebug)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			u := url.URL{Path: "/api/subscribe"}
			q := url.Values{"email": {tt.email}}
			u.RawQuery = q.Encode()

			req, _ := http.NewRequest(http.MethodPost, u.String(), nil)

			subMock := new(SubscriberMock)
			tt.mockFunc(subMock)

			h := NewHandler(subMock, log)

			rw := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Subscribe)

			handler.ServeHTTP(rw, req)
			require.Equal(t, tt.wantCode, rw.Code)

			gotJSON, err := io.ReadAll(rw.Body)
			require.NoError(t, err)

			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestHandlerSendEmails(t *testing.T) {
	cases := map[string]struct {
		mockFunc func(m *SubscriberMock)
		wantCode int
		want     Response
	}{
		"Successful email sending": {
			mockFunc: func(m *SubscriberMock) { m.SendEmailsFunc = func() error { return nil } },
			wantCode: http.StatusOK,
			want:     Response{Message: StatusSend},
		},
		// TODO(not_documented): add more cases.
	}

	log := logger.New(logger.LevelDebug)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/api/sendEmails", nil)

			subMock := new(SubscriberMock)
			tt.mockFunc(subMock)

			h := NewHandler(subMock, log)

			rw := httptest.NewRecorder()
			handler := http.HandlerFunc(h.SendEmails)

			handler.ServeHTTP(rw, req)
			require.Equal(t, tt.wantCode, rw.Code)

			gotJSON, err := io.ReadAll(rw.Body)
			require.NoError(t, err)

			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}
