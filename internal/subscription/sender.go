package subscription

import (
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailClient is a struct that represents an email sender.
type EmailClient struct {
	address string
	client  *sendgrid.Client
}

// NewSender creates a new EmailClient instance with the provided address and API key.
func NewSender(addr, key string) *EmailClient {
	return &EmailClient{
		address: addr,
		client:  sendgrid.NewSendClient(key),
	}
}

// Send sends an email using the provided email address and rate.
func (s *EmailClient) Send(msg Message) error {
	from := mail.NewEmail("Victoria Ray", s.address)
	to := mail.NewEmail(msg.To.Name, msg.To.String())

	email := mail.NewSingleEmailPlainText(from, msg.Subject, to, msg.Body)

	resp, err := s.client.Send(email)
	if err != nil {
		return fmt.Errorf("sending email: %v", err)
	}

	return web.ErrFromStatusCode(resp.StatusCode)
}
