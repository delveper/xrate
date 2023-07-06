package api

import (
	"net/http"
	"path"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/filestore"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	pathRate       = "/rate"
	pathSubscribe  = "/subscribe"
	pathSendEmails = "/sendEmails"
)

func WithRate(cfg ConfigAggregate) Route {
	return func(app *API) {
		grp := path.Join(cfg.Config.Path, cfg.Config.Version)

		client := new(http.Client)
		btcSvc := rate.NewBTCExchangeRateClient(client, cfg.Rate.Endpoint)

		svc := rate.NewService(btcSvc)
		h := rate.NewHandler(svc)

		app.web.Handle(http.MethodGet, grp, pathRate, h.Rate)
	}
}

func WithSubscription(cfg ConfigAggregate) Route {
	return func(app *API) {
		grp := path.Join(cfg.Config.Path, cfg.Config.Version)

		conn := filestore.New[subscription.Subscriber](cfg.Repo.Data)

		client := retryablehttp.NewClient()
		client.RetryMax = cfg.Rate.RetryMax
		btcSvc := rate.NewBTCExchangeRateClient(client.StandardClient(), cfg.Rate.Endpoint)

		svc := subscription.NewService(
			subscription.NewRepo(conn),
			rate.NewService(btcSvc),
			subscription.NewSender(cfg.Sender.Address, cfg.Sender.Key),
		)
		h := subscription.NewHandler(svc)

		app.web.Handle(http.MethodPost, grp, pathSendEmails, h.SendEmails)
		app.web.Handle(http.MethodPost, grp, pathSubscribe, h.Subscribe)
	}
}
