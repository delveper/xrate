/*
Package rate provides a functionality to retrieve exchange rates for digital and fiat currencies.
*/
package rate

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

var ErrInvalidCurrency = fmt.Errorf("invalid currency")

// ExchangeRateProvider is an interface for types that provide exchange rates.
//
//go:generate moq -out=../../test/mock/rate_provider.go -pkg=mock . ExchangeRateProvider
type ExchangeRateProvider interface {
	GetExchangeRate(ctx context.Context, pair CurrencyPair) (*ExchangeRate, error)
	String() string
}

type Service struct {
	bus  *event.Bus
	next *Service
	prov ExchangeRateProvider
}

// NewService constructs a new Service instance.
// Each object in the chain either handles the request or passes it to the next object in the chain.
// Services are chained in the order they are provided, with the first provider in the list being the first one called.
func NewService(bus *event.Bus, provs ...ExchangeRateProvider) *Service {
	svc := (*Service)(nil) // The last instance in the chain.

	for i := len(provs) - 1; i >= 0; i-- {
		svc = &Service{
			prov: provs[i],
			next: svc,
			bus:  bus,
		}
	}

	svc.bus.Subscribe(event.New(EventSource, EventKindRequested, nil), svc.ResponseExchangeRate)
	svc.bus.Subscribe(event.New(EventSource, EventKindFetched, nil), svc.LogExchangeRate)

	return svc
}

// GetExchangeRate attempts to get the exchange rate for a pair of currencies.
// If the Service fails to get the exchange rate, it passes the request to the next Service in the chain, if any.
func (svc *Service) GetExchangeRate(ctx context.Context, pair CurrencyPair) (xrt *ExchangeRate, err error) {
	if err := pair.Validate(); err != nil {
		return nil, err
	}

	defer func() {
		e := event.New(EventSource, EventKindFetched, ProviderResponse{Provider: svc.prov.String(), ExchangeRate: xrt})

		if err != nil {
			err = errors.Join(ErrProviderUnavailable, err)
			e = event.New(EventSource, EventKindFailed, ProviderErrorResponse{Provider: svc.prov.String(), Err: err})
		}

		if errpub := svc.bus.Publish(ctx, e); errpub != nil {
			err = errors.Join(err, errpub)
		}
	}()

	xrt, err = svc.prov.GetExchangeRate(ctx, pair)
	if err != nil && svc.next != nil {
		return svc.next.GetExchangeRate(ctx, pair)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute exchange rate providers chain: %w", err)
	}

	return xrt, nil
}
