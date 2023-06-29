package rate

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerRate(t *testing.T) {
	tests := map[string]struct {
		getterMock Getter
		wantCode   int
		want       any
	}{
		"Valid rate": {
			getterMock: &Getter{
				GetFunc: func(context.Context) (float64, error) { return 2.5, nil },
			},
			wantCode: http.StatusOK,
			want:     Response{Rate: 2.5},
		},
		"Rate retrieval failure": {
			getterMock: &GetterMock{
				GetFunc: func(context.Context) (float64, error) { return 0.0, errors.New("unexpected error") },
			},
			wantCode: http.StatusBadRequest,
			want:     ResponseError{StatusError},
		},
	}
	log := logger.New(logger.LevelDebug)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := NewHandler(tt.getterMock, log)

			rw := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Rate)

			req, err := http.NewRequest(http.MethodGet, "/api/rate", nil)
			require.NoError(t, err)

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
