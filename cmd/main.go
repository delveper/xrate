package main

import (
	"context"
	"fmt"
	"github.com/delveper/gentest/sys/config"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log := logger.New()
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("Startup error", "error", err)
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	type Config struct {
		Web struct {
			Host            string        `env:"API_HOST"`
			ReadTimeout     time.Duration `env:"API_READ_TIMEOUT"`
			WriteTimeout    time.Duration `env:"API_WRITE_TIMEOUT"`
			IdleTimeout     time.Duration `env:"API_IDLE_TIMEOUT"`
			ShutdownTimeout time.Duration `env:"API_SHUTDOWN_TIMEOUT"`
		}
	}

	cfg, err := config.ParseVars[Config]()
	if err != nil {
		return fmt.Errorf("parsing config: %v", err)
	}

	/*Start*/
	log.Infow("Starting service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	mux := *new(http.Handler)

	srv := http.Server{
		Addr:         cfg.Web.Host,
		Handler:      mux,
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
