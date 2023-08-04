package subs

import (
	"net/mail"
)

type Subscriptions []Subscription

// Subscription represents aggregate subscription.
type Subscription struct {
	Subscriber Subscriber
	Topic      Topic
}

// Subscriber represents an entity that subscribes to emails.
type Subscriber struct {
	//TODO: extend entity.
	Address *mail.Address
}

// Topic represents a value object of a topic for subscription.
type Topic = CurrencyPair

// CurrencyPair represents a value object of a currency pair for subscription.
type CurrencyPair struct {
	Base  string
	Quote string
}

// Subscribers implements SubscribersEvent.
func (subss Subscriptions) Subscribers() []string {
	list := make([]string, len(subss))
	for i := range subss {
		list[i] = subss[i].Subscriber.Address.String()
	}

	return list
}

func NewSubscriber(address *mail.Address) Subscriber {
	return Subscriber{Address: address}
}

func NewTopic(base, quote string) Topic {
	return Topic{
		Base:  base,
		Quote: quote,
	}
}
