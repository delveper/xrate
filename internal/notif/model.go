package notif

import "fmt"

// Message represents a message to be sent by service.
type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
}

// Topic represents a value object of a topic for subscription.
type Topic struct {
	Base  string
	Quote string
}

func (t Topic) String() string {
	return fmt.Sprintf("%s/%s", t.Base, t.Quote)
}

// BaseCurrency implements CurrencyPairEvent.
func (t Topic) BaseCurrency() string {
	return t.Base
}

// QuoteCurrency implements CurrencyPairEvent.
func (t Topic) QuoteCurrency() string {
	return t.Quote
}

// ExchangeRateData represents exchange rate data for sending emails.
type ExchangeRateData struct {
	ExchangeRate float64
	Subscribers  []string
	Pair         Topic
}
