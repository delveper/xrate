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

func WithRate(cfg Config) Route {
	return func(app *API) {
		grp := path.Join(cfg.ApiConfig.Path, cfg.ApiConfig.Version)

		client := new(http.Client)
		btcSvc := rate.NewBTCExchangeRateClient(client, cfg.RateConfig.RapidApi.Endpoint)

		svc := rate.NewService(btcSvc)
		hdlr := rate.NewHandler(svc)

		app.web.Handle(http.MethodGet, grp, pathRate, hdlr.Rate)
	}
}

func WithSubscription(cfg Config) Route {
	return func(app *API) {
		grp := path.Join(cfg.ApiConfig.Path, cfg.ApiConfig.Version)

		conn := filestore.New[subscription.Subscriber](cfg.SubscriptionConfig.Data)

		client := retryablehttp.NewClient()
		client.RetryMax = cfg.RateConfig.ClientRetryMax
		btcSvc := rate.NewBTCExchangeRateClient(client.StandardClient(), cfg.RateConfig.RapidApi.Endpoint)

		svc := subscription.NewService(
			subscription.NewRepo(conn),
			rate.NewService(btcSvc),
			subscription.NewSender(cfg.EmailConfig.SenderAddress, cfg.EmailConfig.SenderKey),
		)
		hdlr := subscription.NewHandler(svc)

		app.web.Handle(http.MethodPost, grp, pathSendEmails, hdlr.SendEmails)
		app.web.Handle(http.MethodPost, grp, pathSubscribe, hdlr.Subscribe)
	}
}
