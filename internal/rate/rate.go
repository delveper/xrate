/*
Package rate provides a functionality to retrieve exchange rates for digital and fiat currencies.
*/
package rate

import (
	"context"
	"fmt"
	"strings"
)

// ExchangeRate represents domain exchange rate.
type ExchangeRate struct {
	Value float64
	Pair  CurrencyPair
}

// NewExchangeRate creates a new ExchangeRate instance.
func NewExchangeRate(rate float64, pair CurrencyPair) *ExchangeRate {
	return &ExchangeRate{Value: rate, Pair: pair}
}

// CurrencyPair represents a currency pair.
type CurrencyPair struct {
	Base  string
	Quote string
}

// NewCurrencyPair creates a new CurrencyPair instance.
func NewCurrencyPair(base, quote string) CurrencyPair {
	return CurrencyPair{
		Base:  strings.ToUpper(base),
		Quote: strings.ToUpper(quote),
	}
}

// String converts a CurrencyPair instance to a string.
func (cp CurrencyPair) String() string {
	return fmt.Sprintf("%s/%s", cp.Base, cp.Quote)
}

// OK validates a CurrencyPair instance.
// TODO: Implement validation for all possible currency pairs.
func (cp CurrencyPair) OK() error {
	return nil
}

// ExchangeRateProvider is an interface for types that provide exchange rates.
type ExchangeRateProvider interface {
	GetExchangeRate(ctx context.Context, pair CurrencyPair) (*ExchangeRate, error)
}

type Service struct {
	prov ExchangeRateProvider
	next *Service
}

// NewService constructs a new Service instance.
// Each object in the chain either handles the request or passes it to the next object in the chain.
// Services are chained in the order they are provided, with the first provider in the list being the first one called.
func NewService(provs ...ExchangeRateProvider) *Service {
	var svc *Service

	for i := len(provs) - 1; i >= 0; i-- {
		svc = &Service{
			prov: provs[i],
			next: svc,
		}
	}

	return svc
}

// GetExchangeRate attempts to get the exchange rate for a pair of currencies.
// If the Service fails to get the exchange rate, it passes the request to the next Service in the chain, if any.
func (svc *Service) GetExchangeRate(ctx context.Context, pair CurrencyPair) (*ExchangeRate, error) {
	if err := pair.OK(); err != nil {
		return nil, err
	}

	xrt, err := svc.prov.GetExchangeRate(ctx, pair)
	if err != nil && svc.next != nil {
		return svc.next.GetExchangeRate(ctx, pair)
	}

	return xrt, nil
}
