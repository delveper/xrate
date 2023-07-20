package email

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
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

type Service struct {
	auth smtp.Auth
	tmpl *template.Template
}

func NewService(tmpl *template.Template, cfg Config) *Service {
	auth := smtp.PlainAuth("", cfg.UserName, cfg.Password, cfg.Host)

	return &Service{tmpl: tmpl, auth: auth}
}

// Send responsible for sending an email message.
func (svc *Service) Send(ctx context.Context, msg *notif.Message) error {
	var buf bytes.Buffer

	if err := svc.tmpl.ExecuteTemplate(&buf, "email", msg); err != nil {
		return fmt.Errorf("executing email template: %v", err)
	}

	return smtp.SendMail(msg.From, svc.auth, msg.From, msg.To, buf.Bytes())

}
