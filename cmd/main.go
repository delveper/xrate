package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/api"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/env"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
)

func main() {
	log := logger.New(logger.LevelDebug, "./log/sys.log")
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup error", "error", err)
	}
}

func run(log *logger.Logger) error {
	var cfg struct {
		Api struct {
			Name    string `default:"gensch"`
			Path    string `default:"/api"`
			Version string `default:"v1"`
			Origin  string `default:"*"`
		}
		Web struct {
			Host            string        `default:"0.0.0.0:9999"`
			ReadTimeout     time.Duration `default:"5s"`
			WriteTimeout    time.Duration `default:"10s"`
			IdleTimeout     time.Duration `default:"60s"`
			ShutdownTimeout time.Duration `default:"15s"`
		}
		Repo struct {
			Data string `default:"/data"`
		}
		Rate struct {
			Provider struct {
				RapidApi struct { // https://rapidapi.com/Serply/api/exchange-rate9
					Endpoint string `default:"https://exchange-rate9.p.rapidapi.com/symbols"`
					Key      string
				}
				Ninjas struct { //https://api-ninjas.com/api/exchangerate
					Endpoint string `default:"https://api.api-ninjas.com/v1/exchangerate?pair"`
					Key      string
				}
				AlphaVantage struct { // https://www.alphavantage.co/documentation/#fx
					Endpoint string `default:"https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE"`
					Key      string
				}
				CoinApi struct { // https://docs.coinapi.io/market-data/rest-api/exchange-rates
					Endpoint string `default:"https://rest.coinapi.io/v1/exchangerate"`
					Key      string
				}
				CoinYep struct {
					Endpoint string `default:"https://coinyep.com/api/v1/"`
					Key      string
				}
			}
			Client struct {
				RetryMax int `default:"10"`
			}
		}
		Sender struct {
			Address string
			Key     string
		}
	}

	if err := env.ParseTo(".env", &cfg); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	log.Infow("starting service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	app := api.New(api.ConfigAggregate{
		Api: api.Config(cfg.Api),
		Rate: rate.Config{
			Provider: struct {
				RapidApi, Ninjas, AlphaVantage, CoinApi, CoinYep rate.ProviderConfig
			}{rate.ProviderConfig(cfg.Rate.Provider.RapidApi),
				rate.ProviderConfig(cfg.Rate.Provider.Ninjas),
				rate.ProviderConfig(cfg.Rate.Provider.AlphaVantage),
				rate.ProviderConfig(cfg.Rate.Provider.CoinApi),
				rate.ProviderConfig(cfg.Rate.Provider.CoinYep)},
			Client: struct{ RetryMax int }{cfg.Rate.Client.RetryMax},
		},
		Subscription: subscription.Config{
			Sender: subscription.SenderConfig(cfg.Sender),
			Repo:   subscription.RepoConfig(cfg.Repo),
		},
	},
		shutdown, log)

	srv := http.Server{
		Addr:         cfg.Web.Host,
		Handler:      app.Handler(),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     log.ToStandard(),
	}

	errSrv := make(chan error, 1)
	go func() {
		log.Infow("startup", "status", "api router started", "host", srv.Addr)
		errSrv <- srv.ListenAndServe()
	}()

	select {
	case err := <-errSrv:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("Shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return fmt.Errorf("shuting down gracefully: %w", err)
		}
	}

	return nil
}
