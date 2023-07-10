package subs

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const (
	EventSource           = "subs"
	EventTopicRateRequest = "rate_request"
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
	req := event.New(EventSource, EventTopicRateRequest, toRequestEventData(topic))
	ch := make(chan event.Event)
	req.Response = ch

	if err := svc.bus.Dispatch(ctx, req); err != nil {
		return 0, err
	}

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case e := <-ch:
		resp, ok := e.Data.(ResponseEvent)
		if !ok {
			return 0, fmt.Errorf("%w: %T", ErrInvalidEvent, e)
		}

		return resp.ExchangeRate(), nil
	}
}
