package sndr

import (
	"net/mail"
)

// Message represents a message to be sent by service.
type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
}

// ExchangeRate represents exchange rate.
type ExchangeRate struct {
	Value float64
	Pair  CurrencyPair
}

// CurrencyPair represents a currency pair.
type CurrencyPair struct {
	Base  string
	Quote string
}

// Subscription represents aggregate subscription.
type Subscription struct {
	Subscriber Subscriber
	Topic      Topic
}

// Topic represents a value object of a topic for subscription.
type Topic = CurrencyPair

// Subscriber represents an entity that subscribes to emails.
type Subscriber struct {
	Address *mail.Address
}
