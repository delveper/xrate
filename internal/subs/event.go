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

type CurrencyPairEvent interface {
	BaseCurrency() string
	QuoteCurrency() string
}

// RespondSubscription handles an event send subscriptions data.
func (svc *Service) RespondSubscription(ctx context.Context, e event.Event) error {
	req, ok := e.Payload.(CurrencyPairEvent)
	if !ok {
		return fmt.Errorf("%w: unexpected payload: %T", ErrInvalidEvent, e.Payload)
	}

	subss, err := svc.List(ctx, NewTopic(req.BaseCurrency(), req.QuoteCurrency()))
	if err != nil {
		return fmt.Errorf("responding subscription event: %w", err)
	}

	if e.Response == nil {
		return fmt.Errorf("responding exchange rate event: %w", ErrInvalidEvent)
	}

	e.Response <- event.New(EventSource, EventKindResponded, Subscriptions(subss))

	return nil
}
