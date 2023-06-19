package rate

import (
	"encoding/json"
	"errors"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerRate(t *testing.T) {
	cases := map[string]struct {
		mockFunc func(m *GetterMock)
		wantCode int
		want     any
	}{
		"Valid rate": {
			mockFunc: func(m *GetterMock) { m.GetFunc = func() (float64, error) { return 2.5, nil } },
			wantCode: http.StatusOK,
			want:     Response{Rate: 2.5},
		},
		"Rate retrieval failure": {
			mockFunc: func(m *GetterMock) {
				m.GetFunc = func() (float64, error) {
					return 0, errors.New("failed to retrieve rate")
				}
			},
			wantCode: http.StatusBadRequest,
			want:     ResponseError{StatusError},
		},
	}

	log := logger.New(logger.LevelDebug)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/rate", nil)

			getterMock := new(GetterMock)
			tt.mockFunc(getterMock)

			h := NewHandler(getterMock, log)

			rw := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Rate)

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
