package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/app/api/ctrl"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/notif/email"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subs"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/env"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
)

type config struct {
	Log struct {
		Level     string   `default:"debug"`
		FilePath  string   `default:"sys.log"`
		Topic     string   `default:"logs"`
		Partition int32    `default:"0"`
		Brokers   []string `default:"localhost:9092,kafka:9092"`
		Offset    string   `default:"oldest"`
	}
	Api struct {
		Name    string `default:"rate"`
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
		Data string `default:"./data"`
	}
	Rate struct {
		Provider struct {
			ExchangeRateHost struct {
				Name     string `default:"ExchangeRateHost"`
				Endpoint string `default:"https://api.exchangerate.host/latest"`
				Header   string `default:"-"`
				Key      string `default:"-"`
			}
			Ninjas struct {
				Name     string `default:"Ninjas"`
				Endpoint string `default:"https://api.api-ninjas.com/v1/exchangerate"`
				Header   string `default:"X-Api-Key"`
				Key      string
			}
			AlphaVantage struct {
				Name     string `default:"AlphaVantage"`
				Endpoint string `default:"https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE"`
				Header   string `default:"apikey"`
				Key      string
			}
			CoinApi struct {
				Name     string `default:"CoinApi"`
				Endpoint string `default:"https://rest.coinapi.io/v1/exchangerate"`
				Header   string `default:"X-CoinAPI-Key"`
				Key      string
			}
			CoinYep struct {
				Name     string `default:"CoinYep"`
				Endpoint string `default:"https://coinyep.com/api/v1/"`
				Header   string `default:"-"`
				Key      string `default:"-"`
			}
		}
	}
	Email struct {
		Host     string `default:"smtp.ionos.com"`
		Port     string `default:"465"`
		UserName string `default:"yevhen@bilyk.dev"`
		Password string
	}
}

func main() {
	var cfg config
	if err := env.ParseTo(&cfg); err != nil {
		log.Fatalf("failed to parse env: %v", err)
	}

	log := logger.New(
		logger.WithConsoleCore(cfg.Log.Level),
		logger.WithJSONCore(cfg.Log.Level, cfg.Log.FilePath),
		logger.WithKafkaCore(cfg.Log.Level, cfg.Log.Topic, cfg.Log.Brokers...),
	)

	defer log.Sync()

	if err := run(log, &cfg); err != nil {
		log.Errorw("startup error", "error", err)
	}
}

func run(log *logger.Logger, cfg *config) error {
	log.Infow("starting service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	app, err := ctrl.New(ctrl.ConfigAggregate{
		Api: ctrl.Config(cfg.Api),
		Rate: rate.Config{
			Provider: struct{ ExchangeRateHost, Ninjas, AlphaVantage, CoinApi, CoinYep rate.ProviderConfig }{
				ExchangeRateHost: rate.ProviderConfig(cfg.Rate.Provider.ExchangeRateHost),
				Ninjas:           rate.ProviderConfig(cfg.Rate.Provider.Ninjas),
				AlphaVantage:     rate.ProviderConfig(cfg.Rate.Provider.AlphaVantage),
				CoinApi:          rate.ProviderConfig(cfg.Rate.Provider.CoinApi),
				CoinYep:          rate.ProviderConfig(cfg.Rate.Provider.CoinYep)},
		},
		Subscription: subs.Config{RepoData: cfg.Repo.Data},
		Email:        email.Config(cfg.Email),
	}, shutdown, log)

	if err != nil {
		return fmt.Errorf("configuring service: %w", err)
	}

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
