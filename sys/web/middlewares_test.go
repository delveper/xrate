package web_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestChainMiddlewares(t *testing.T) {
	key := new(any)

	mw := func(str string) web.Middleware {
		return func(next web.Handler) web.Handler {
			return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
				val, _ := ctx.Value(key).(string)
				val += str
				ctx = context.WithValue(ctx, key, val)

				return next(ctx, rw, req)
			}
		}
	}

	var want string

	const mwNum = 10

	mws := make([]web.Middleware, mwNum)

	for i := 0; i < mwNum; i++ {
		str := strconv.Itoa(i)
		mws[i] = mw(str)
		want += str
	}

	handler := func(ctx context.Context, rw http.ResponseWriter, _ *http.Request) error {
		rw.Write([]byte(fmt.Sprint(ctx.Value(key))))
		return nil
	}

	chainedHandler := web.ChainMiddlewares(handler, mws...)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	err := chainedHandler(context.WithValue(context.Background(), key, ""), rw, req)
	require.NoError(t, err)

	require.Equal(t, want, rw.Body.String())
}

func TestMiddlewares(t *testing.T) {
	log := logger.New(logger.WithConsoleCore(logger.LevelDebug))
	defer log.Sync()

	const target = "http://example.com/foo/bar"

	t.Run("WithJSON", func(t *testing.T) {
		rw := httptest.NewRecorder()
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			rw.Write([]byte("test"))
			return nil
		}
		req := httptest.NewRequest(http.MethodGet, target, nil)

		mw := web.WithJSON(h)
		err := mw(context.Background(), rw, req)

		require.NoError(t, err)
		require.Equal(t, "application/json; charset=UTF-8", rw.Header().Get("Content-Type"))
	})

	t.Run("WithLogRequest", func(t *testing.T) {
		rw := httptest.NewRecorder()

		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			rw.Write([]byte("test"))
			return nil
		}
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.Header.Set("X-Request-ID", "")

		mw := web.WithLogRequest(log)(h)
		err := mw(context.Background(), rw, req)

		require.NoError(t, err)
		require.Equal(t, "test", rw.Body.String())

		_, err = uuid.Parse(req.Header.Get("X-Request-ID"))
		require.NoError(t, err)
	})

	t.Run("WithRecover", func(t *testing.T) {
		rw := httptest.NewRecorder()
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("test")
		}
		req := httptest.NewRequest(http.MethodGet, target, nil)

		mw := web.WithRecover(log)(h)
		err := mw(context.Background(), rw, req)

		require.NoError(t, err)
	})

	t.Run("WithCORS", func(t *testing.T) {
		rw := httptest.NewRecorder()
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			rw.Write([]byte(`{"test":"test"}`))
			return nil
		}
		req := httptest.NewRequest(http.MethodGet, target, nil)

		mw := web.WithCORS("*")(h)
		err := mw(context.Background(), rw, req)

		require.NoError(t, err)
		require.Equal(t, `{"test":"test"}`, rw.Body.String())
	})

	t.Run("WithErrors", func(t *testing.T) {
		rw := httptest.NewRecorder()
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return errors.New("test error")
		}
		req := httptest.NewRequest(http.MethodGet, target, nil)

		mw := web.WithErrors(log)(h)
		err := mw(context.Background(), rw, req)

		require.NoError(t, err)
		require.JSONEq(t, `{"error":"unknown error"}`, rw.Body.String())
	})
}
