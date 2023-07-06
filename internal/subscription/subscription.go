// Package subscription provides functionality to manage subscriptions.
package subscription

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
)

const defaultTimeout = 5 * time.Second

var (
	// ErrEmailAlreadyExists is an error indicating that the email address already exists in the database.
	ErrEmailAlreadyExists = errors.New("email address exists")

	// ErrMissingEmail is an error indicating that the email address is missing.
	ErrMissingEmail = errors.New("missing email")
)

// Subscriber represents an entity that subscribes to emails.
type Subscriber struct {
	Address *mail.Address
	Topic   Topic
}

func NewSubscriber(address *mail.Address, topic Topic) *Subscriber {
	return &Subscriber{Address: address, Topic: topic}
}

// Topic represents a topic that subscribes to emails.
type Topic = string

// Message represents an email message.
type Message struct {
	From    *mail.Address
	To      *mail.Address
	Subject string
	Body    string
}

func NewMessage(subject, body string, to *mail.Address) Message {
	return Message{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

//go:generate moq -out=../../test/mock/email_repository.go -pkg=mock . SubscriberRepository

// SubscriberRepository is an interface for managing email subscriptions.
type SubscriberRepository interface {
	Add(Subscriber) error
	List() ([]Subscriber, error)
}

//go:generate moq -out=../../test/mock/email_sender.go -pkg=mock . EmailSender

// EmailSender is an interface for sending emails.
type EmailSender interface {
	Send(Message) error
}

// Service represents a service that manages email subscriptions and sends emails.
type Service struct {
	rate rate.ExchangeRateService
	repo SubscriberRepository
	mail EmailSender
}

// NewService creates a new Service instance with the provided dependencies.
func NewService(repo SubscriberRepository, rate rate.ExchangeRateService, mail EmailSender) *Service {
	return &Service{
		repo: repo,
		rate: rate,
		mail: mail,
	}
}

// Subscribe adds a new email subscription to the repository.
func (svc *Service) Subscribe(sub Subscriber) error {
	if sub.Topic == "" {
		sub.Topic = rate.NewCurrencyPair(rate.CurrencyBTC, rate.CurrencyUAH).String()
	}

	if err := svc.repo.Add(sub); err != nil {
		return fmt.Errorf("adding subscription: %w", err)
	}

	return nil
}

// SendEmails sends emails to all subscribers using the current rate.
func (svc *Service) SendEmails() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	pair := rate.NewCurrencyPair(rate.CurrencyBTC, rate.CurrencyUAH)

	rate, err := svc.rate.Get(ctx, pair)
	if err != nil {
		return err
	}

	subscribers, err := svc.repo.List()
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s exchange rate at %s", pair, time.Now().Format(time.Stamp))
	body := fmt.Sprintf("Current exhange rate: %f", rate.Value)

	var errArr []error

	for _, sub := range subscribers {
		msg := NewMessage(
			subject,
			body,
			sub.Address,
		)

		if err := svc.mail.Send(msg); err != nil {
			errArr = append(errArr, err)
		}
	}

	if errArr != nil {
		return errors.Join(errArr...)
	}

	return nil
}
