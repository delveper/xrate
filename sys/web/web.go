/*
Package web contains implementation of a web framework.

The Web struct is the main type of the package and consists of:

	mux: The httprouter router that is used to route HTTP requests to handlers.
	mws: A slice of middlewares that are applied to all HTTP requests before they are handled by the handler function.
	sig: A channel that is used to receive shutdown signals.
*/
package web

import (
	"context"
	"net/http"
	"os"
	"path"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

// Handler is responsible for handling HTTP requests.
type Handler = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error

// Web is a web framework.
type Web struct {
	mux *httprouter.Router
	mws []Middleware
	sig chan os.Signal
}

// New creates a new Web struct.
func New(sig chan os.Signal, mds ...Middleware) *Web {
	return &Web{
		mux: httprouter.New(),
		mws: mds,
		sig: sig,
	}
}

// ServeHTTP Serves HTTP requests.
func (w *Web) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w.mux.ServeHTTP(rw, req)
}

// Handle registers a handler function for a specific HTTP method and path.
func (w *Web) Handle(meth string, grp string, pth string, hdlr Handler, mws ...Middleware) {
	hdlr = ChainMiddlewares(hdlr, append(w.mws, mws...)...)

	fn := func(rw http.ResponseWriter, req *http.Request) {
		if err := hdlr(req.Context(), rw, req); err != nil {
			if _, ok := IsError[*ShutdownError](err); ok {
				w.Shutdown()
				return
			}
		}
	}

	pth = path.Clean(path.Join("/", grp, pth))

	w.mux.HandlerFunc(meth, pth, fn)
}

// Shutdown shutdowns the web application.
func (w *Web) Shutdown() {
	w.sig <- syscall.SIGTERM
}
