package rate

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=getter_mock_test.go -package=rate github.com/GenesisEducationKyiv/main-project-delveper/internal/rate Getter

func TestHandlerRate(t *testing.T) {
	tests := map[string]struct {
		getterBuilder func(*gomock.Controller) Getter
		wantCode      int
		want          any
	}{
		"Valid rate": {
			getterBuilder: func(ctrl *gomock.Controller) Getter {
				mock := NewMockGetter(ctrl)
				mock.EXPECT().
					Get().
					Return(2.5, nil).
					Times(1)
				return mock
			},
			wantCode: http.StatusOK,
			want:     Response{Rate: 2.5},
		},
		"Rate retrieval failure": {
			getterBuilder: func(ctrl *gomock.Controller) Getter {
				mock := NewMockGetter(ctrl)
				mock.EXPECT().
					Get().
					Return(0.0, errors.New("failed to retrieve rate")).
					Times(1)
				return mock
			},
			wantCode: http.StatusBadRequest,
			want:     ResponseError{StatusError},
		},
	}

	log := logger.New(logger.LevelDebug)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest(http.MethodGet, "/api/rate", nil)
			require.NoError(t, err)

			getter := tt.getterBuilder(ctrl)
			h := NewHandler(getter, log)

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
