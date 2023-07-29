package rate

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const (
	EventSource        = "rate"
	EventKindResponded = "responded"
	EventKindRequested = "requested"
	EventKindFetched   = "fetched"
	EventKindFailed    = "failed"
)

var (
	ErrInvalidEvent   = errors.New("invalid event")
	ErrInvalidChannel = errors.New("response channel is not initialized")
)

// CurrencyPairEvent represents a currency pair fetching event.
//go:generate moq -out=../../test/mock/rate_event.go -pkg=mock . CurrencyPairEvent
type CurrencyPairEvent interface {
	BaseCurrency() string
	QuoteCurrency() string
}

// ProviderResponse represents the data of a provider response event.
type ProviderResponse struct {
	Provider     string
	ExchangeRate *ExchangeRate
}

// ProviderErrorResponse represents the data of a provider error event.
type ProviderErrorResponse struct {
	Provider string
	Err      error
}

// LogExchangeRate is an event listener designed for log the exchange rate fetching.
func (svc *Service) LogExchangeRate(ctx context.Context, e event.Event) error {
	switch e.Payload.(type) {
	case ProviderResponse, ProviderErrorResponse:
		// The Payload will be logged by event dispatcher.
	default:
		return fmt.Errorf("%w: unexpected payload: %T", ErrInvalidEvent, e.Payload)
	}

	return nil
}

// ResponseExchangeRate in an event listener that fetches the exchange rate for a requested currency pair.
func (svc *Service) ResponseExchangeRate(ctx context.Context, e event.Event) error {
	req, ok := e.Payload.(CurrencyPairEvent)
	if !ok {
		return fmt.Errorf("%w: unexpected payload, expected CurrencyPairEvent: %T", ErrInvalidEvent, e.Payload)
	}

	xrt, err := svc.GetExchangeRate(ctx, NewCurrencyPair(req.BaseCurrency(), req.QuoteCurrency()))
	if err != nil {
		return fmt.Errorf("responding exchange rate event: %w", err)
	}

	if e.Response == nil {
		return fmt.Errorf("responding exchange rate event: %w", ErrInvalidChannel)
	}

	e.Response <- event.New(EventSource, EventKindResponded, xrt)

	return nil
}
