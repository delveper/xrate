package sndr

import (
	"bytes"
	"fmt"
	"text/template"
)

type ExchangeRateContent struct{ *template.Template }

func NewExchangeRateContent(tmpl *template.Template) *ExchangeRateContent {
	return &ExchangeRateContent{tmpl}
}

func (c *ExchangeRateContent) Create(md *MetaData) (*Message, error) {
	tos := make([]string, len(md.subss))

	for i := range md.subss {
		tos[i] = md.subss[i].Subscriber.Address.String()
	}

	var buf bytes.Buffer
	if err := c.ExecuteTemplate(&buf, "subject", md); err != nil {
		return nil, fmt.Errorf("executing subjecttemplate: %w", err)
	}

	subj := buf.String()

	buf.Reset()
	if err := c.ExecuteTemplate(&buf, "body", md); err != nil {
		return nil, fmt.Errorf("executing subjecttemplate: %w", err)
	}

	body := buf.String()

	msg := Message{
		To:      tos,
		Subject: subj,
		Body:    body,
	}

	return &msg, nil
}
