// Package rate provides functionality to retrieve and handle exchange rates.
package rate

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

const defaultTimeout = 15 * time.Second

//go:generate moq -out=../../test/mock/getter.go -pkg=mock . ExchangeRateService

// ExchangeRateService interface to get rate from external service.
type ExchangeRateService interface {
	GetExchangeRate(ctx context.Context, currency CurrencyPair) (*ExchangeRate, error)
}

type Response struct {
	Rate float64
}

func NewResponse(rate *ExchangeRate) *Response {
	return &Response{Rate: rate.Value}
}

// Handler structure for handling rate requests.
type Handler struct {
	rate ExchangeRateService
}

// NewHandler creates a new Handler instance.
func NewHandler(rate ExchangeRateService) Handler {
	return Handler{rate: rate}
}

// Rate handles the HTTP request for the rate.
func (h *Handler) Rate(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	pair := NewCurrencyPair(
		web.FromQuery(req, "base"),
		web.FromQuery(req, "quote"),
	)

	if err := pair.OK(); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	rate, err := h.rate.GetExchangeRate(ctx, pair)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return web.NewRequestError(err, http.StatusRequestTimeout)
		}

		return err
	}

	return web.Respond(ctx, rw, NewResponse(rate), http.StatusOK)
}
