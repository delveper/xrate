package web

import (
	"context"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/google/uuid"
)

// Middleware is a middleware that implements a series of middleware to an HTTP h function in a chain-like manner.
type Middleware = func(Handler) Handler

// ChainMiddlewares chains a series of Middlewares applied from top to bottom order for better readability.
// For example, if you provide slice [MiddlewareA, MiddlewareB, MiddlewareC], the actual execution order
// would be MiddlewareA -> MiddlewareB -> MiddlewareC.
func ChainMiddlewares(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

// WithJSON is a middleware that sets the response content ype JSON.
func WithJSON(h Handler) Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

		return h(ctx, rw, req)
	}
}

// WithLogRequest logs every request.
func WithLogRequest(log *logger.Logger) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			start := time.Now()

			defer func() {
				var id string
				if id = FromHeader(req, "X-Request-ID", ""); id == "" {
					id = uuid.New().String()
				}

				req.Header.Set("X-Request-ID", id)

				log.Debugw("request completed",
					"id", id,
					"uri", req.RequestURI,
					"method", req.Method,
					"duration", time.Since(start),
				)
			}()

			return h(ctx, rw, req)
		}
	}
}

// WithRecover recovers application from panic with logging stack trace.
func WithRecover(log *logger.Logger) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			defer func() {
				if rec := recover(); rec != nil {
					log.Errorw("recovered from panic",
						"rec", rec,
						"trace", string(debug.Stack()),
					)
				}
			}()

			return h(ctx, rw, req)
		}
	}
}

// WithCORS is a middleware that ensures that the HTTP
// method of the request matches the provided method.
func WithCORS(origins ...string) Middleware {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
	headers := []string{"Accept", "Authorization", "Accept-Encoding", "Content-Type", "Content-Length", "X-CSRF-Token", "X-Request-ID"}

	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			rw.Header().Set("Access-Control-Allow-Origin", strings.Join(origins, ", "))
			rw.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
			rw.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))

			return h(ctx, rw, req)
		}
	}
}

// WithErrors is a middleware that wraps an HTTP h to provide centralized error handling.
func WithErrors(log *logger.Logger) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := h(ctx, rw, req); err != nil {
				log.Errorw("error", "message", err)

				if _, ok := IsError[*ShutdownError](err); ok {
					return err
				}

				resp, code := defineErrorResponse(err)
				if err := Respond(ctx, rw, resp, code); err != nil {
					return err
				}
			}

			return nil
		}
	}
}

// defineErrorResponse determines the HTTP response message and status code based on the provided error.
// If the error is of an unknown type, it returns a generic internal server error message and status code 500.
func defineErrorResponse(err error) (resp ErrorResponse, code int) {
	resp.Error = ErrUnknownError.Error()
	code = http.StatusInternalServerError

	if err, ok := IsError[*RequestError](err); ok {
		resp.Error = err.Error()
		code = err.StatusCode
	}

	return resp, code
}
