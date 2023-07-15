package subs

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const (
	EventSource         = "subs"
	EventKindSubscribed = "subscribed"
)

var ErrInvalidEvent = errors.New("invalid event")

type RequestEvent interface {
	BaseCurrency
	QuoteCurrency
}

type ResponseEvent interface {
	ExchangeRate
}

// BaseCurrency represents a base currency fetcher.
type BaseCurrency interface{ BaseCurrency() string }

// QuoteCurrency represents a quote currency fetcher.
type QuoteCurrency interface{ QuoteCurrency() string }

type ExchangeRate interface{ ExchangeRate() float64 }

// RequestEventData represents the data of a request event.
type RequestEventData struct {
	baseCurrency  string
	quoteCurrency string
}

func toRequestEventData(t Topic) *RequestEventData {
	return &RequestEventData{
		baseCurrency:  t.BaseCurrency,
		quoteCurrency: t.QuoteCurrency,
	}
}

// BaseCurrency returns the base currency of the RequestEventData.
func (rd *RequestEventData) BaseCurrency() string {
	return rd.baseCurrency
}

// QuoteCurrency returns the quote currency of the RequestEventData.
func (rd *RequestEventData) QuoteCurrency() string {
	return rd.quoteCurrency
}

func (svc *Service) RequestExchangeRate(ctx context.Context, topic Topic) (float64, error) {
	e := event.New(EventSource, "subscribed", toRequestEventData(topic))

	if err := svc.bus.Publish(ctx, e); err != nil {
		return 0, fmt.Errorf("publishing rate request event: %w", err)
	}

	select {
	case <-ctx.Done():
		return 0, ctx.Err()

	case e := <-e.Response:
		switch resp := e.Payload.(type) {
		case ResponseEvent:
			return resp.ExchangeRate(), nil
		default:
			return 0, fmt.Errorf("%w: unexpected payload: %T", ErrInvalidEvent, e.Payload)
		}
	}
}
