// Package subscription provides functionality to manage subscriptions.
package subs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const defaultTimeout = 15 * time.Second

const (
	currencyBTC = "BTC"
	currencyUAH = "UAH"
)

var (
	// ErrSubscriptionExists is an error indicating that the email address already exists in the database.
	ErrSubscriptionExists = errors.New("subscription already exists")

	// ErrMissingEmail is an error indicating that the email address is missing.
	ErrMissingEmail = errors.New("missing email")
)

//go:generate moq -out=../../test/mock/email_repository.go -pkg=mock . SubscriberRepository

// SubscriberRepository is an interface for managing email subscriptions.
type SubscriberRepository interface {
	Add(context.Context, Subscription) error
	List(context.Context) ([]Subscription, error)
}

//go:generate moq -out=../../test/mock/email_sender.go -pkg=mock . EmailSender

// EmailSender is an interface for sending emails.
type EmailSender interface {
	Send(Message) error
}

// Service represents a service that manages email subscriptions and sends emails.
type Service struct {
	bus  *event.Bus
	repo SubscriberRepository
	mail EmailSender
}

// NewService creates a new Service instance with the provided dependencies.
func NewService(bus *event.Bus, repo SubscriberRepository, mail EmailSender) *Service {
	return &Service{
		bus:  bus,
		repo: repo,
		mail: mail,
	}
}

// Subscribe adds a new email subscription to the repository.
func (svc *Service) Subscribe(ctx context.Context, subs Subscription) error {
	if err := svc.repo.Add(ctx, subs); err != nil {
		return fmt.Errorf("adding subscription: %w", err)
	}
	// TODO: Add event for handling new subscription.
	return nil
}

// SendEmails sends emails to all subscribers using the current rate.
func (svc *Service) SendEmails(ctx context.Context, topic Topic) error {
	val, err := svc.RequestExchangeRate(ctx, topic)
	if err != nil {
		return fmt.Errorf("requesting exchange rate: %w", err)
	}

	subscriptions, err := svc.repo.List(ctx)
	if err != nil {
		return fmt.Errorf("listing subscriptions: %w", err)
	}

	var n int
	for _, subs := range subscriptions {
		if subs.Topic == topic {
			subscriptions[n] = subs
			n++
		}
	}

	if n == 0 {
		return fmt.Errorf("no subscriptions for topic %s found", topic)
	}

	subject := fmt.Sprintf("%s exchange rate at %s", topic, time.Now().Format(time.Stamp))
	body := fmt.Sprintf("Current exhange rate: %f", val)

	var errs []error

	for _, subs := range subscriptions[:n] {
		msg := NewMessage(
			subject,
			body,
			subs.Subscriber.Address,
		)

		if err := svc.mail.Send(msg); err != nil {
			errs = append(errs, err)
		}
	}

	if errs != nil {
		return fmt.Errorf("sending emails: %w", errors.Join(errs...))
	}

	return nil
}
