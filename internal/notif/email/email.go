package email

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/smtp"
	"text/template"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/notif"
)

type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
}

// Service represents an email service.
type Service struct {
	cfg  Config
	auth smtp.Auth
	tmpl *template.Template
}

func NewService(tmpl *template.Template, cfg Config) *Service {
	auth := smtp.PlainAuth("", cfg.UserName, cfg.Password, cfg.Host)
	return &Service{cfg: cfg, tmpl: tmpl, auth: auth}
}

// Send responsible for sending an email message.
func (svc *Service) Send(ctx context.Context, msg *notif.Message) error {
	var buf bytes.Buffer

	msg.From = svc.cfg.UserName

	if err := svc.tmpl.ExecuteTemplate(&buf, "email", msg); err != nil {
		return fmt.Errorf("executing email template: %v", err)
	}

	addr := net.JoinHostPort(svc.cfg.Host, svc.cfg.Port)

	if err := smtp.SendMail(addr, svc.auth, msg.From, msg.To, buf.Bytes()); err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
