package api

import (
	"net/http"
	"path"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate/curxrt"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/sndr"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/sndr/email"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/sndr/tmpl"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subs"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/filestore"
)

const (
	pathRate       = "/rate"
	pathSubscribe  = "/subscribe"
	pathSendEmails = "/sendEmails"
)

// WithRate set-ups routes for handling HTTP requests related to currency exchange rates.
func WithRate(cfg ConfigAggregate) Route {
	return func(app *App) error {
		grp := path.Join(cfg.Api.Path, cfg.Api.Version)
		clt := new(http.Client)
		provs := []rate.ExchangeRateProvider{
			curxrt.NewProvider[curxrt.ExchangeRateHost](curxrt.Config(cfg.Rate.Provider.ExchangeRateHost), clt),
			curxrt.NewProvider[curxrt.AlphaVantage](curxrt.Config(cfg.Rate.Provider.AlphaVantage), clt),
			curxrt.NewProvider[curxrt.CoinYep](curxrt.Config(cfg.Rate.Provider.CoinYep), clt),
			curxrt.NewProvider[curxrt.CoinApi](curxrt.Config(cfg.Rate.Provider.CoinApi), clt),
			curxrt.NewProvider[curxrt.Ninjas](curxrt.Config(cfg.Rate.Provider.Ninjas), clt),
		}
		svc := rate.NewService(app.bus, provs...)
		h := rate.NewHandler(svc)

		app.web.Handle(http.MethodGet, grp, pathRate, h.Rate)

		return nil
	}
}

// WithSubscription set-ups routes related to subscription functionality.
func WithSubscription(cfg ConfigAggregate) Route {
	return func(app *App) error {
		grp := path.Join(cfg.Api.Path, cfg.Api.Version)
		conn := filestore.New[subs.Subscription](cfg.Subscription.Repo.Data)
		repo := subs.NewRepo(conn)
		svc := subs.NewService(app.bus, repo)
		h := subs.NewHandler(svc)

		app.web.Handle(http.MethodPost, grp, pathSubscribe, h.Subscribe)

		return nil
	}
}

// WithNotification set-ups routes related to notification functionality.
func WithNotification(cfg ConfigAggregate) Route {
	return func(app *App) error {
		grp := path.Join(cfg.Api.Path, cfg.Api.Version)
		clt := email.NewSMTPClient(cfg.Notification)
		t, err := tmpl.Load()
		if err != nil {
			return err
		}

		mail := email.NewService(clt, t)
		cont := sndr.NewExchangeRateContent(t)

		svc := sndr.NewService(app.bus, mail, cont)
		h := sndr.NewHandler(svc)

		app.web.Handle(http.MethodPost, grp, pathSendEmails, h.SendEmails)

		return nil
	}
}
