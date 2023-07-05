/*
Package rate provides a functionality to retrieve BTC rates against Fiat currencies.
By default, the exchange rate fetched for a BTC/UAH pair.
Service could potentially be used for other pairs, including Fiat currencies in both directions.
*/
package rate

import (
	"context"
	"errors"
	"fmt"
)

const (
	CurrencyBTC = "BTC"
	CurrencyUAH = "UAH"
)

var ErrInvalidCurrency = errors.New("invalid currency pair")

// ExchangeRate represents domain exchange rate.
type ExchangeRate struct {
	Value float64
	Pair  CurrencyPair
}

func NewExchangeRate(rate float64, pair CurrencyPair) *ExchangeRate {
	return &ExchangeRate{Value: rate, Pair: pair}
}

// CurrencyPair represents a currency pair.
type CurrencyPair struct {
	Base  string
	Quote string
}

func NewCurrencyPair(base, quote string) CurrencyPair {
	return CurrencyPair{Base: base, Quote: quote}
}

func (cp CurrencyPair) String() string {
	return fmt.Sprintf("%s/%s", cp.Base, cp.Quote)
}

func (cp CurrencyPair) OK() error {
	if cp.Base != CurrencyBTC {
		return ErrInvalidCurrency
	}

	return nil
}

type BTCExchangeRateProvider interface {
	GetBTCExchangeRate(ctx context.Context, currency string) (float64, error)
}

type Service struct {
	BTCExchangeRateProvider
}

func NewService(svc BTCExchangeRateProvider) *Service {
	return &Service{BTCExchangeRateProvider: svc}
}

func (svc *Service) GetExchangeRate(ctx context.Context, pair CurrencyPair) (*ExchangeRate, error) {
	if err := pair.OK(); err != nil {
		return nil, err
	}

	val, err := svc.BTCExchangeRateProvider.GetBTCExchangeRate(ctx, pair.Quote)
	if err != nil {
		return nil, err
	}

	rate := NewExchangeRate(val, pair)

	return rate, nil
}
