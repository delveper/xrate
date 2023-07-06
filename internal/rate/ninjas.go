package rate

import (
	"context"
)

type NinjasAdapter = ExchangeRateProviderAdapter

func NewNinjasAdapter(client HTTPClient, cfg AdapterOption) *NinjasAdapter {
	return &NinjasAdapter{client: client, config: cfg}
}

func (a *NinjasAdapter) GetExchangeRate(ctx context.Context, pair CurrencyPair) (*ExchangeRate, error) {
	return nil, nil
}

func NewAdapterRequest(ctx context.Context, pair CurrencyPair) (*ExchangeRate, error) {
	return nil, nil
}
