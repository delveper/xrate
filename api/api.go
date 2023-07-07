/*
Package api provides the main App server for the application.
It includes functions for creating new App instances, handling incoming
HTTP requests, adding middlewares and applying various application routes.
*/
package api

import (
	"net/http"
	"os"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// App is the main application instance.
type App struct {
	sig chan os.Signal
	log *logger.Logger
	web *web.Web
}

// Route is a function that defines an application route.
type Route func(*App)

// New returns a new App instance with provided configuration.
func New(cfg ConfigAggregate, sig chan os.Signal, log *logger.Logger) *App {
	mws := []web.Middleware{
		web.WithLogRequest(log),
		web.WithCORS(cfg.Api.Origin),
		web.WithErrors(log),
		web.WithJSON,
		web.WithRecover(log),
	}

	web := web.New(sig, mws...)

	api := App{
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

func (a *App) Handler() http.Handler {
	return a.web
}

// Routes applies all application routes.
func (a *App) Routes(routes ...Route) {
	for i := range routes {
		routes[i](a)
	}
}
