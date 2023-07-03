package api

import (
	"net/http"
	"path"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

const (
	pathRate       = "/rate"
	pathSubscribe  = "/subscribe"
	pathSendEmails = "/sendEmails"
)

func WithRate(cfg Config) Route {
	return func(app *API) {
		grp := path.Join(cfg.ApiConfig.Path, cfg.ApiConfig.Version)
		// TODO: Bind all dependencies.
		dummy := *new(web.Handler)

		app.web.Handle(http.MethodGet, grp, pathRate, dummy)
	}
}

func WithSubscription(cfg Config) Route {
	return func(app *API) {
		grp := path.Join(cfg.ApiConfig.Path, cfg.ApiConfig.Version)
		// TODO: Bind all dependencies.
		dummy := *new(web.Handler)

		app.web.Handle(http.MethodPost, grp, pathSendEmails, dummy)
		app.web.Handle(http.MethodPost, grp, pathSubscribe, dummy)
	}
}
