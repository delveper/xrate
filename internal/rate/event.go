package rate

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const (
	EventSource         = "rate"
	EventKindResponded  = "responded"
	EventKindSubscribed = "subscribed"
	EventKindFetched    = "fetched"
	EventKindFailed     = "failed"
)

var (
	ErrInvalidEvent   = errors.New("invalid event")
	ErrInvalidChannel = errors.New("response channel is not initialized")
)

// RequestEvent represents a request event to fetch the exchange rate.
type RequestEvent interface {
	BaseCurrency
	QuoteCurrency
}

// BaseCurrency represents a base currency fetcher.
type BaseCurrency interface{ BaseCurrency() string }

// QuoteCurrency represents a quote currency fetcher.
type QuoteCurrency interface{ QuoteCurrency() string }

// ResponseEventData represents the data of a response event.
type ResponseEventData struct {
	baseCurrency  string
	quoteCurrency string
	exchangeRate  float64
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

func toResponseEventData(xrt *ExchangeRate) *ResponseEventData {
	return &ResponseEventData{
		baseCurrency:  xrt.Pair.Base,
		quoteCurrency: xrt.Pair.Quote,
		exchangeRate:  xrt.Value,
	}
}

// BaseCurrency returns the base currency of the ResponseEventData.
func (rd *ResponseEventData) BaseCurrency() string {
	return rd.baseCurrency
}

// QuoteCurrency returns the quote currency of the ResponseEventData.
func (rd *ResponseEventData) QuoteCurrency() string {
	return rd.quoteCurrency
}

// ExchangeRate returns the exchange rate of the ResponseEventData.
func (rd *ResponseEventData) ExchangeRate() float64 {
	return rd.exchangeRate
}

// RespondExchangeRate handles an event and fetches the exchange rate for a requested currency pair.
func (svc *Service) RespondExchangeRate(ctx context.Context, e event.Event) error {
	req, ok := e.Payload.(RequestEvent)
	if !ok {
		return fmt.Errorf("%w: %T", ErrInvalidEvent, e)
	}

	xrt, err := svc.GetExchangeRate(ctx, NewCurrencyPair(req.BaseCurrency(), req.QuoteCurrency()))
	if err != nil {
		return fmt.Errorf("responding exchange rate event: %w", err)
	}

	if e.Response == nil {
		return fmt.Errorf("responding exchange rate event: %w", ErrInvalidChannel)
	}

	e.Response <- event.New(EventSource, EventKindResponded, toResponseEventData(xrt))

	return nil
}

// LogExchangeRate handles an event and logs the exchange rate for a requested currency pair.
func (svc *Service) LogExchangeRate(ctx context.Context, e event.Event) error {
	switch e.Payload.(type) {
	case ProviderResponse, ProviderErrorResponse:
		// The Payload will be logged by event dispatcher.
	default:
		return fmt.Errorf("logging provider response: %w", ErrInvalidEvent)
	}

	return nil
}
