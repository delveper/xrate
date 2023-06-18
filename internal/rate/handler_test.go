package rate

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/delveper/gentest/sys/logger"
)

func TestHandlerRate(t *testing.T) {
	cases := map[string]struct {
		mockGetter func(m *GetterMock)
		wantCode   int
		wantRate   float64
		wantError  string
	}{
		"valid rate": {
			mockGetter: func(m *GetterMock) {
				m.GetFunc = func() (float64, error) {
					return 2.5, nil
				}
			},
			wantCode:  http.StatusOK,
			wantRate:  2.5,
			wantError: "",
		},
		"rate retrieval failure": {
			mockGetter: func(m *GetterMock) {
				m.GetFunc = func() (float64, error) {
					return 0, errors.New("failed to retrieve rate")
				}
			},
			wantCode:  http.StatusBadRequest,
			wantRate:  0,
			wantError: "Failed to get rate",
		},
	}

	for key, tt := range cases {
		t.Run(key, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/rate", nil)

			mockGetter := new(GetterMock)
			tt.mockGetter(mockGetter)

			log := logger.New("debug")
			h := NewHandler(mockGetter, log)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Rate)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}

			if tt.wantCode == http.StatusOK {
				var resp struct{ Rate float64 }
				err := json.NewDecoder(rr.Body).Decode(&resp)
				if err != nil {
					t.Fatal("failed to parse response")
				}

				if resp.Rate != tt.wantRate {
					t.Errorf("handler returned wrong rate: got %v want %v",
						resp.Rate, tt.wantRate)
				}
			} else {
				var resp struct{ Error string }
				err := json.NewDecoder(rr.Body).Decode(&resp)
				if err != nil {
					t.Fatal("failed to parse error response")
				}

				if resp.Error != tt.wantError {
					t.Errorf("handler returned wrong error: got %v want %v",
						resp.Error, tt.wantError)
				}
			}
		})
	}
}
