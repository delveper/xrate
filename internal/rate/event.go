package rate

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const (
	EventSource       = "rate"
	EventTypeResponse = "response"
	EventTypeFetched  = "fetched"
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
	req, ok := e.Data.(RequestEvent)
	if !ok {
		return fmt.Errorf("%w: %T", ErrInvalidEvent, e)
	}

	xrt, err := svc.GetExchangeRate(ctx, NewCurrencyPair(req.BaseCurrency(), req.QuoteCurrency()))
	if err != nil {
		return err
	}

	if e.Response == nil {
		return ErrInvalidChannel
	}

	e.Response <- event.New(EventSource, EventTypeResponse, toResponseEventData(xrt))

	return nil
}

func (svc *Service) logProviderEvent(ctx context.Context, xrt *ExchangeRate, err error) {
	data := struct {
		Provider     string
		ExchangeRate *ExchangeRate
		Error        error
	}{
		Provider:     svc.prov.String(),
		ExchangeRate: xrt,
		Error:        err,
	}

	e := event.New(EventSource, EventTypeFetched, data)

	h := func(context.Context, event.Event) error { return nil }

	svc.bus.Publish(e, h)
	svc.bus.Dispatch(ctx, e)
}
