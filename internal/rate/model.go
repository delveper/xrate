package rate

import (
	"fmt"
	"strings"
)

// ExchangeRate represents exchange rate.
type ExchangeRate struct {
	Value float64
	Pair  CurrencyPair
}

// CurrencyPair represents a currency pair.
type CurrencyPair struct {
	Base  string
	Quote string
}

// NewExchangeRate creates a new ExchangeRate instance.
func NewExchangeRate(rate float64, pair CurrencyPair) *ExchangeRate {
	return &ExchangeRate{Value: rate, Pair: pair}
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

// Validate validates a CurrencyPair instance.
// TODO: Improve validation for all possible currency pairs.
func (cp CurrencyPair) Validate() error {
	if cp.Base == "" || cp.Quote == "" {
		return fmt.Errorf("%w: %+v", ErrInvalidCurrency, cp)
	}

	return nil
}
