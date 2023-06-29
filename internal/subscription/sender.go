package subscription

import (
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Sender is a struct that represents an email sender.
type Sender struct {
	address string
	client  *sendgrid.Client
}

// NewSender creates a new Sender instance with the provided address and API key.
func NewSender(addr, key string) *Sender {
	return &Sender{
		address: addr,
		client:  sendgrid.NewSendClient(key),
	}
}

// Send sends an email using the provided email address and rate.
func (s *Sender) Send(email Email, rate float64) error {
	subject := "Current BTC to UAH rate"

	from := mail.NewEmail("Example Use", s.address)
	to := mail.NewEmail(email.Address.Name, email.Address.String())

	textContent := "Current rate is:"
	htmlContent := fmt.Sprintf("<strong>%f</strong>", rate)

	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)

	resp, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("sending email: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		return nil

	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden:
		return fmt.Errorf("client error: %d", resp.StatusCode)

	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return fmt.Errorf("server error: %d", resp.StatusCode)

	default:
		return fmt.Errorf("unexpected: %d", resp.StatusCode)
	}
}
