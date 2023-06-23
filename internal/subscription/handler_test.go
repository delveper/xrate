package subscription

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=subscriber_mock_test.go -package=subscription github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription Subscriber

func TestHandlerSubscribe(t *testing.T) {
	tests := map[string]struct {
		email             string
		subscriberBuilder func(controller *gomock.Controller) Subscriber
		wantCode          int
		want              Response
	}{
		"Success": {
			email: "email@example.com",
			subscriberBuilder: func(ctrl *gomock.Controller) Subscriber {
				mock := NewMockSubscriber(ctrl)
				mock.EXPECT().
					Subscribe(gomock.Any()).
					Return(nil).
					Times(1)
				return mock
			},
			wantCode: http.StatusOK,
			want:     Response{Message: StatusSubscribed},
		},
		"Email already exists": {
			email: "exists@example.com",
			subscriberBuilder: func(ctrl *gomock.Controller) Subscriber {
				mock := NewMockSubscriber(ctrl)
				mock.EXPECT().
					Subscribe(gomock.Any()).
					Return(ErrEmailAlreadyExists).
					Times(1)
				return mock
			},
			wantCode: http.StatusConflict,
			want:     Response{Message: StatusError, Details: ErrEmailAlreadyExists.Error()},
		},
		// TODO(not_documented): add case for internal server error.
	}

	log := logger.New(logger.LevelDebug)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest(http.MethodPost, "/api/subscribe", nil)
			require.NoError(t, err)

			q := req.URL.Query()
			q.Add("email", tt.email)
			req.URL.RawQuery = q.Encode()

			sub := tt.subscriberBuilder(ctrl)

			h := NewHandler(sub, log)

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
		subscriberBuilder func(controller *gomock.Controller) Subscriber
		wantCode          int
		want              Response
	}{
		"Successful email sending": {
			subscriberBuilder: func(ctrl *gomock.Controller) Subscriber {
				mock := NewMockSubscriber(ctrl)
				mock.EXPECT().
					SendEmails().
					Return(nil).
					Times(1)
				return mock
			},
			wantCode: http.StatusOK,
			want:     Response{Message: StatusSend},
		},
		// TODO(not_documented): Add case for internal server error.
	}

	log := logger.New(logger.LevelDebug)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest(http.MethodPost, "/api/sendEmails", nil)
			require.NoError(t, err)

			sub := tt.subscriberBuilder(ctrl)
			h := NewHandler(sub, log)

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
