package main

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/notif"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/notif/email"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/notif/tmpl"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
)

func main() {
	log := logger.New(logger.LevelDebug, "sys.log")

	cfg := email.Config{
		Host:     "smtp.ionos.com",
		Port:     "465",
		UserName: "yevhen@bilyk.dev",
		Password: "jncjcbrjptgbc.y",
	}

	/*
		addr := net.JoinHostPort(cfg.Host, cfg.Port)
		tlsCfg := tls.Config{InsecureSkipVerify: false, ServerName: cfg.Host}
		conn, err := tls.Dial("tcp", addr, &tlsCfg)
		if err != nil {
			log.Errorf("connecting SMTP server: %v", err)
		}

		clt, err := smtp.NewClient(conn, cfg.Host)

		auth := smtp.PlainAuth("", cfg.UserName, cfg.Password, cfg.Host)
		if err := clt.Auth(auth); err != nil {
			log.Fatalf("starting %s authentication: %v", auth, err)
		}

		if err := clt.Mail("yevhen@bilyk.dev"); err != nil {
			log.Fatal("mail: %v", err)
		}

		if err := clt.Rcpt("yevhen@bilyk.dev"); err != nil {
			log.Fatal("rcpt: %v", err)
		}

		if _, err = clt.Data(); err != nil {
			log.Fatal("data: %v", err)
		}

		if err := clt.Quit(); err != nil {
			log.Fatal("quit: %v", err)
		}
	*/
	clt := email.NewSMTPClient(cfg)
	t, err := tmpl.Load()
	if err != nil {
		log.Fatalf("failed to load template: %v", err)
	}
	msg := &notif.Message{
		From:    "yevhen@bilyk.dev",
		To:      []string{"yevhen@bilyk.dev"},
		Subject: "Hello",
		Body:    "Hello",
	}

	mail := email.NewService(clt, t)

	if err := mail.Send(context.Background(), msg); err != nil {
		log.Fatalf("failed to send email: %v", err)
	}
}
