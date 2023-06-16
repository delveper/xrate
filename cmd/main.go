package main

import (
	"context"
	"fmt"
	"github.com/delveper/gentest/internal/transport"
	"github.com/delveper/gentest/sys/env"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log := logger.New(logger.LEVEL_DEBUG)
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("Startup error", "error", err)
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	type Config struct {
		Repo struct {
			Path string `env:"DB_PATH"`
		}
		Mail struct {
			APIKey  string `env:"EMAIL_KEY"`
			Address string `env:"EMAIL_ADDRESS"`
		}
		Web struct {
			Host            string        `env:"API_HOST"`
			ReadTimeout     time.Duration `env:"API_READ_TIMEOUT"`
			WriteTimeout    time.Duration `env:"API_WRITE_TIMEOUT"`
			IdleTimeout     time.Duration `env:"API_IDLE_TIMEOUT"`
			ShutdownTimeout time.Duration `env:"API_SHUTDOWN_TIMEOUT"`
		}
	}

	cfg, err := env.Parse[Config](".env")
	if err != nil {
		return fmt.Errorf("parsing config: %v", err)
	}

	log.Infow("Starting service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	api := transport.New(
		transport.Config{
			DBPath:       cfg.Repo.Path,
			EmailAPIkey:  cfg.Mail.APIKey,
			EmailAddress: cfg.Mail.Address,
		}, log)

	srv := http.Server{
		Addr:         cfg.Web.Host,
		Handler:      api.Handle(),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     log.ToStandard(),
	}

	errSrv := make(chan error, 1)
	go func() {
		log.Infow("Startup", "status", "api router started", "host", srv.Addr)
		errSrv <- srv.ListenAndServe()
	}()

	select {
	case err := <-errSrv:
		return fmt.Errorf("server error: %w", err)

	case <-shutdown:
		log.Infow("Shutdown", "status", "shutdown started")
		defer log.Infow("Shutdown", "status", "shutdown complete")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return fmt.Errorf("shuting down gracefully: %w", err)
		}
	}

	return nil
}
