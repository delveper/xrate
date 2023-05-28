package subscription

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/delveper/gentest/sys/logger"
)

func TestHandler_Subscribe(t *testing.T) {
	tests := map[string]struct {
		email    string
		mockSub  func(m *SubscriberMock)
		wantCode int
	}{
		"valid email": {
			email: "email@example.com",
			mockSub: func(m *SubscriberMock) {
				m.SubscribeFunc = func(email Email) error {
					return nil
				}
			},
			wantCode: http.StatusOK,
		},
		"invalid email": {
			email: "invalidemail",
			mockSub: func(m *SubscriberMock) {
				m.SubscribeFunc = func(email Email) error {
					return nil
				}
			},
			wantCode: http.StatusBadRequest,
		},
		"email already exists": {
			email: "exists@example.com",
			mockSub: func(m *SubscriberMock) {
				m.SubscribeFunc = func(email Email) error {
					return ErrEmailAlreadyExists
				}
			},
			wantCode: http.StatusConflict,
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/subscribe", strings.NewReader("email="+tt.email))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			mockSub := new(SubscriberMock)
			tt.mockSub(mockSub)

			log := logger.New("debug")
			h := NewHandler(mockSub, log)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Subscribe)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}
		})
	}
}

func TestHandler_SendEmails(t *testing.T) {
	tests := map[string]struct {
		mockSub  func(m *SubscriberMock)
		wantCode int
	}{
		"successful email sending": {
			mockSub: func(m *SubscriberMock) {
				m.SendEmailsFunc = func() error {
					return nil
				}
			},
			wantCode: http.StatusOK,
		},
		"failed email sending": {
			mockSub: func(m *SubscriberMock) {
				m.SendEmailsFunc = func() error {
					return errors.New("failed to send emails")
				}
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/sendEmails", nil)

			mockSub := new(SubscriberMock)
			tt.mockSub(mockSub)

			log := logger.New("debug")
			h := NewHandler(mockSub, log)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.SendEmails)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}
		})
	}
}
