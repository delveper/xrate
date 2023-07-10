package api

import (
	"net/http"
	"path"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate/curxrt"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subs"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/filestore"
)

const (
	pathRate       = "/rate"
	pathSubscribe  = "/subscribe"
	pathSendEmails = "/sendEmails"
)

// WithRate set-ups routes for handling HTTP requests related to currency exchange rates.
func WithRate(cfg ConfigAggregate) Route {
	return func(app *App) {
		grp := path.Join(cfg.Api.Path, cfg.Api.Version)
		clt := new(http.Client)
		svc := rate.NewService(event.NewBus(app.log),
			curxrt.NewProvider[curxrt.ExchangeRateHost](curxrt.Config(cfg.Rate.Provider.ExchangeRateHost), clt),
			curxrt.NewProvider[curxrt.AlphaVantage](curxrt.Config(cfg.Rate.Provider.AlphaVantage), clt),
			curxrt.NewProvider[curxrt.CoinYep](curxrt.Config(cfg.Rate.Provider.CoinYep), clt),
			curxrt.NewProvider[curxrt.CoinApi](curxrt.Config(cfg.Rate.Provider.CoinApi), clt),
			curxrt.NewProvider[curxrt.Ninjas](curxrt.Config(cfg.Rate.Provider.Ninjas), clt),
		)
		h := rate.NewHandler(svc)

		app.web.Handle(http.MethodGet, grp, pathRate, h.Rate)
	}
}

// WithSubscription set-ups routes related to subscription functionality.
func WithSubscription(cfg ConfigAggregate) Route {
	return func(app *App) {
		grp := path.Join(cfg.Api.Path, cfg.Api.Version)
		conn := filestore.New[subs.Subscriber](cfg.Subscription.Repo.Data)
		svc := subs.NewService(
			event.NewBus(app.log),
			subs.NewRepo(conn),
			subs.NewSender(cfg.Subscription.Sender.Address, cfg.Subscription.Sender.Key),
		)
		h := subs.NewHandler(svc)

		app.web.Handle(http.MethodPost, grp, pathSubscribe, h.Subscribe)
		app.web.Handle(http.MethodPost, grp, pathSendEmails, h.SendEmails)
	}
}
