package adpt

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

type CoinYep = Adapter

func NewCoinYep(client HTTPClient, cfg Config) *Adapter {
	return NewAdapter(client, cfg.Endpoint)
}

type respCoinYep struct {
	BaseSymbol   string  `json:"base_symbol"`
	BaseName     string  `json:"base_name"`
	TargetSymbol string  `json:"target_symbol"`
	TargetName   string  `json:"target_name"`
	Price        string  `json:"price"`
	PriceChange  float64 `json:"price_change"`
}

func (r respCoinYep) toExchangeRate() (*rate.ExchangeRate, error) {
	val, err := strconv.ParseFloat(r.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing rate: %w", err)
	}

	xrate := rate.ExchangeRate{
		Value: val,
		Pair:  rate.NewCurrencyPair(r.BaseSymbol, r.TargetSymbol),
	}

	return &xrate, nil
}

func (a *CoinYep) GetExchangeRate(ctx context.Context, pair rate.CurrencyPair) (*rate.ExchangeRate, error) {
	resp, err := a.SendRequest(ctx,
		web.WithValue("from", pair.Base),
		web.WithValue("to", pair.Quote),
		web.WithValue("lang", "en"),
		web.WithValue("format", "json"),
	)
	if err != nil {
		return nil, err
	}

	data, err := web.ProcessResponse[respCoinYep](resp)
	if err != nil {
		return nil, err
	}

	return data.toExchangeRate()
}
