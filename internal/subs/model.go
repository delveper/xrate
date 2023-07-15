package subs

import (
	"net/mail"
)

// Subscription represents aggregate subscription.
type Subscription struct {
	Subscriber Subscriber
	Topic      Topic
}

// Subscriber represents an entity that subscribes to emails.
type Subscriber struct {
	// TODO: extend entity.
	Address *mail.Address
}

// Topic represents a value object of a topic for subscription.
type Topic struct {
	BaseCurrency  string
	QuoteCurrency string
}

// Message represents an email message.
type Message struct {
	From    *mail.Address
	To      *mail.Address
	Subject string
	Body    string
}

func NewSubscriber(address *mail.Address) Subscriber {
	return Subscriber{Address: address}
}

func NewTopic(base, quote string) Topic {
	return Topic{
		BaseCurrency:  base,
		QuoteCurrency: quote,
	}
}

func NewMessage(subject, body string, to *mail.Address) Message {
	return Message{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}
