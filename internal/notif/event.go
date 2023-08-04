package notif

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const (
	EventSource        = "notif"
	EventKindRequested = "requested"
)

var ErrInvalidEvent = errors.New("invalid event")

// CurrencyPairEvent is an event for fetching a currency pair.
// In package notif it is preemptive used for showing that CurrencyPair implements it.
type CurrencyPairEvent interface {
	BaseCurrency() string
	QuoteCurrency() string
}

// ExchangeRateEvent is an event for fetching exchange rate.
type ExchangeRateEvent interface {
	ExchangeRate() float64
}

// SubscribersEvent is an event for fetching list of subscribers.
type SubscribersEvent interface {
	Subscribers() []string
}

// RequestExchangeRateData triggers fetching exchange rate data.
func (svc *Service) RequestExchangeRateData(ctx context.Context, pair Topic) (*ExchangeRateData, error) {
	e := event.New(EventSource, EventKindRequested, pair)
	if err := svc.bus.Publish(ctx, e); err != nil {
		return nil, fmt.Errorf("publishing sending email event: %w", err)
	}

	var (
		xrt   float64
		subss []string
	)

	for xrt == 0 || subss == nil {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("requesting exchange rate data: %w", ctx.Err())

		case e := <-e.Response:
			switch val := e.Payload.(type) {
			case ExchangeRateEvent:
				xrt = val.ExchangeRate()

			case SubscribersEvent:
				subss = val.Subscribers()

			default:
				return nil, fmt.Errorf("%w: unexpected payload: %T", ErrInvalidEvent, e.Payload)
			}
		}
	}

	return NewExchangeRateData(pair, xrt, subss), nil
}
