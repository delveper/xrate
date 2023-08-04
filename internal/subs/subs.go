/*
Package subs provides functionality to manage subscriptions.
*/
package subs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const defaultTimeout = 15 * time.Second

var (
	// ErrSubscriptionExists is an error indicating that the email address already exists in the database.
	ErrSubscriptionExists = errors.New("subscription already exists")

	// ErrMissingEmail is an error indicating that the email address is missing.
	ErrMissingEmail = errors.New("missing email")

	// ErrNotFound is an error indicating that the subscription was not found in the database.
	ErrNotFound = errors.New("subscription not found")
)

// SubscriberRepository is an interface for managing email subscriptions.
//
//go:generate moq -out=../../test/mock/email_repository.go -pkg=mock . SubscriberRepository
type SubscriberRepository interface {
	Add(context.Context, Subscription) error
	List(context.Context) ([]Subscription, error)
}

// Service represents a service that manages email subscriptions and sends emails.
type Service struct {
	bus  *event.Bus
	repo SubscriberRepository
}

// NewService creates a new Service instance with the provided dependencies.
func NewService(bus *event.Bus, repo SubscriberRepository) *Service {
	svc := Service{
		bus:  bus,
		repo: repo,
	}

	svc.bus.Subscribe(event.New(EventSource, EventKindRequested, nil), svc.RespondSubscription)

	return &svc
}

// Subscribe adds a new email subscription to the repository.
func (svc *Service) Subscribe(ctx context.Context, subs Subscription) error {
	if err := svc.repo.Add(ctx, subs); err != nil {
		return fmt.Errorf("adding subscription: %w", err)
	}
	// TODO: Add event for handling new subscription.
	return nil
}

// List returns all subscriptions from the repository specified by topic.
func (svc *Service) List(ctx context.Context, topic Topic) ([]Subscription, error) {
	subscriptions, err := svc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing subscriptions: %w", err)
	}

	var n int

	for _, subs := range subscriptions {
		if subs.Topic == topic {
			subscriptions[n] = subs
			n++
		}
	}

	if n == 0 {
		return nil, fmt.Errorf("%w for topic %s", ErrNotFound, topic)
	}

	return subscriptions[:n], nil
}
