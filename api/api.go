/*
Package api can be seen as the Controller layer that responsible for handling incoming HTTP requests,
applying the necessary middlewares, and delegating the requests to the appropriate handlers (Use Case Interactors).
The handlers then interact with the domain logic to process the request and generate a response.
*/
package api

import (
	"net/http"
	"os"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

// App is the main application instance.
type App struct {
	sig chan os.Signal
	log *logger.Logger
	web *web.Web
	bus *event.Bus
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

	api := App{
		sig: sig,
		log: log,
		web: web.New(sig, mws...),
		bus: event.NewBus(log),
	}

	api.Routes(
		WithRate(cfg),
		WithSubscription(cfg),
	)

	return &api
}

// Handler returns the web handler.
func (a *App) Handler() http.Handler {
	return a.web
}

// Routes applies all application routes.
func (a *App) Routes(routes ...Route) {
	for i := range routes {
		routes[i](a)
	}
}
