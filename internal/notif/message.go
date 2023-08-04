package notif

import (
	"bytes"
	"fmt"
	"text/template"
)

// ExchangeRateContent responsible for creating exchange rate content for sending messages.
type ExchangeRateContent struct{ *template.Template }

func NewExchangeRateContent(tmpl *template.Template) *ExchangeRateContent {
	return &ExchangeRateContent{tmpl}
}

// CreateMessage creates message content using prepared templates.
func (c *ExchangeRateContent) CreateMessage(data *ExchangeRateData) (*Message, error) {
	var buf bytes.Buffer
	if err := c.ExecuteTemplate(&buf, "subject", data); err != nil {
		return nil, fmt.Errorf("executing subject template: %w", err)
	}

	subj := buf.String()
	buf.Reset()

	if err := c.ExecuteTemplate(&buf, "body", data); err != nil {
		return nil, fmt.Errorf("executing body template: %w", err)
	}

	body := buf.String()

	return NewMessage(data.Subscribers, subj, body), nil
}
