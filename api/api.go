/*
Package api provides the main API server for the application.
It includes functions for creating new API instances, handling incoming
HTTP requests, adding middlewares and applying various application routes.
*/
package api

import (
	"net/http"
	"os"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// API is the main application instance.
type API struct {
	sig chan os.Signal
	log *logger.Logger
	web *web.Web
}

type Route func(*API)

// New returns a new API instance with provided configuration.
func New(cfg ConfigAggregate, sig chan os.Signal, log *logger.Logger) *API {
	mws := []web.Middleware{
		web.WithLogRequest(log),
		web.WithCORS(cfg.Config.Origin),
		web.WithErrors(log),
		web.WithJSON,
		web.WithRecover(log),
	}

	web := web.New(sig, mws...)

	api := API{
		sig: sig,
		log: log,
		web: web,
	}

	api.Routes(
		WithRate(cfg),
		WithSubscription(cfg),
	)

	return &api
}

func (a *API) Handler() http.Handler {
	return a.web
}

// Routes applies all application routes.
func (a *API) Routes(routes ...Route) {
	for i := range routes {
		routes[i](a)
	}
}
