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
	"syscall"
	"time"
)

func main() {
	log := logger.New(logger.LevelDebug)
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("Startup error", "error", err)
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	type Config struct {
		Web struct {
			Host            string
			ReadTimeout     time.Duration
			WriteTimeout    time.Duration
			IdleTimeout     time.Duration
			ShutdownTimeout time.Duration
		}
		Db struct {
			DataPath string
		}
		Email struct {
			SenderAddress string
			SenderKey     string
		}
	}

	cfg, err := env.Parse[Config](".env")
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	log.Infow("Starting service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := transport.New(
		transport.Config{
			DBPath:       cfg.Db.DataPath,
			EmailAPIkey:  cfg.Email.SenderKey,
			EmailAddress: cfg.Email.SenderAddress,
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

	case sig := <-shutdown:
		log.Infow("Shutdown", "status", "shutdown started", "signal", sig)
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
