package email

import (
	"bytes"
	"context"
	_ "embed"
	"text/template"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/sndr"
)

type SMTPClient interface {
	Auth() error
	Connect() error
	Write([]byte) error
	Close() error
}

type Service struct {
	clt  SMTPClient
	tmpl *template.Template
}

func NewService(clt SMTPClient, tmpl *template.Template) *Service {
	return &Service{clt: clt, tmpl: tmpl}
}

// Send responsible for sending an email message.
func (svc *Service) Send(ctx context.Context, msg *sndr.Message) (err error) {
	if err := svc.clt.Connect(); err != nil {
		return err
	}

	defer func() {
		if errc := svc.clt.Close(); errc != nil {
			err = errc
		}
	}()

	if err := svc.clt.Auth(); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := svc.tmpl.ExecuteTemplate(&buf, "email", msg); err != nil {
		return err
	}

	return svc.clt.Write(buf.Bytes())
}
