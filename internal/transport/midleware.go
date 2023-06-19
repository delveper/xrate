package transport

import (
	"net/http"
)

// ChainMiddlewares applies a series of middleware to an HTTP handler function in a chain-like manner.
// Middlewares are applied in the reverse order that they are provided, meaning the last middleware provided
// is the first one to process the request.
func ChainMiddlewares(hdl http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		hdl = middlewares[i](hdl)
	}

	return hdl
}

// WithJSON is a middleware that sets the response content ype JSON.
func (a *API) WithJSON(hdl http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		hdl.ServeHTTP(rw, req)
	}
}

// WithCORS is a middleware that sets the Access-Control-Allow headers for CORS.
func (a *API) WithCORS(hdl http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		rw.Header().Set("Access-Control-Allow-Credentials", "false")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		hdl.ServeHTTP(rw, req)
	}
}

// WithLogRequest logs every request and sends logger instance to further handler.
func (a *API) WithLogRequest(hdl http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		a.log.Debugw("Request:",
			"method", req.Method,
			"uri", req.RequestURI,
			"user-agent", req.UserAgent(),
			"remote", req.RemoteAddr,
		)

		hdl.ServeHTTP(rw, req)
	}
}

// WithoutPanic recovers in case panic, but we won't panic.
func (a *API) WithoutPanic(hdl http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				a.log.Errorw("Recovered from panic.", "rec", rec)
				rw.WriteHeader(http.StatusInternalServerError)
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		hdl.ServeHTTP(rw, req)
	}
}

// WithMethod is a middleware that ensures that the HTTP method of the request matches the provided method.
func (a *API) WithMethod(meth string) func(http.HandlerFunc) http.HandlerFunc {
	return func(hdl http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			if req.Method != meth {
				a.log.Errorw("Method not allowed", "method", meth)
				http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

				return
			}

			hdl.ServeHTTP(rw, req)
		}
	}
}
