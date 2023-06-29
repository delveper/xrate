package rate_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate/mocks"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerRate(t *testing.T) {
	tests := map[string]struct {
		getterMock rate.Getter
		wantCode   int
		want       any
	}{
		"Valid rate": {
			getterMock: &mocks.GetterMock{
				GetFunc: func(context.Context) (float64, error) { return 2.5, nil },
			},
			wantCode: http.StatusOK,
			want:     rate.Response{Rate: 2.5},
		},
		"Rate retrieval failure": {
			getterMock: &mocks.GetterMock{
				GetFunc: func(context.Context) (float64, error) { return 0.0, errors.New("unexpected error") },
			},
			wantCode: http.StatusBadRequest,
			want:     rate.ResponseError{rate.StatusError},
		},
	}
	log := logger.New(logger.LevelDebug)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			h := rate.NewHandler(tc.getterMock, log)

			rw := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Rate)

			req, err := http.NewRequest(http.MethodGet, "/api/rate", nil)
			require.NoError(t, err)

			handler.ServeHTTP(rw, req)
			require.Equal(t, tc.wantCode, rw.Code)

			gotJSON, err := io.ReadAll(rw.Body)
			require.NoError(t, err)

			wantJSON, err := json.Marshal(tc.want)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}
