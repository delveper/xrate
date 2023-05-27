package subscription

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

type Sender struct {
	apiKey string
}

func NewSender(apikey string) *Sender {
	return &Sender{apiKey: apikey}
}

func (c *Sender) Send(email Email, rate float64) error {
	from := mail.NewEmail("Example Use", "rufa.matviyiv@empeek.tech")
	subject := "Current BTC to UAH rate"
	to := mail.NewEmail(email.Address.Name, email.Address.String())
	plainTextContent := "Current rate is:"
	htmlContent := fmt.Sprintf("<strong>%f</strong>", rate)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(c.apiKey)

	resp, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("sending email: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected OK, got: %d", resp.StatusCode)
	}

	return nil
}
