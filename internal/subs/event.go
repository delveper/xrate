package subs

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const (
	EventSource        = "subs"
	EventKindRequested = "requested"
	EventKindResponded = "responded"
)

var ErrInvalidEvent = errors.New("invalid event")

// RespondSubscription handles an event send subscriptions data.
func (svc *Service) RespondSubscription(ctx context.Context, e event.Event) error {
	subscriptions, err := svc.repo.List(ctx)
	if err != nil {
		return fmt.Errorf("responding subscription event: %w", err)
	}

	if e.Response == nil {
		return fmt.Errorf("responding exchange rate event: %w", ErrInvalidEvent)
	}

	e.Response <- event.New(EventSource, EventKindResponded, subscriptions)

	return nil
}
