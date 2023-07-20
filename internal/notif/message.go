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

func (c *ExchangeRateContent) CreateMessage(md *MetaData) (*Message, error) {
	tos := make([]string, len(md.subss))

	for i := range md.subss {
		tos[i] = md.subss[i].Subscriber.Address.String()
	}

	var buf bytes.Buffer
	if err := c.ExecuteTemplate(&buf, "subject", md); err != nil {
		return nil, fmt.Errorf("executing body template: %w", err)
	}

	subj := buf.String()

	buf.Reset()
	if err := c.ExecuteTemplate(&buf, "body", md); err != nil {
		return nil, fmt.Errorf("executing subject template: %w", err)
	}

	body := buf.String()

	msg := Message{
		To:      tos,
		Subject: subj,
		Body:    body,
	}

	return &msg, nil
}
