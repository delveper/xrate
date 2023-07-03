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

//go:generate moq -out=../../test/mocks/getter.go -pkg=mocks . ExchangeRateService

// ExchangeRateService interface to get rate from external service.
type ExchangeRateService interface {
	Get(ctx context.Context, currency CurrencyPair) (*ExchangeRate, error)
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
func NewHandler(rate ExchangeRateService) *Handler {
	return &Handler{rate: rate}
}

// Rate handles the HTTP request for the rate.
func (h *Handler) Rate(ctx context.Context, rw http.ResponseWriter, _ *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	rate, err := h.rate.Get(ctx, NewCurrencyPair(CurrencyBTC, CurrencyUAH))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return web.Respond(ctx, rw, rate, http.StatusRequestTimeout)
		}
		return err
	}

	return web.Respond(ctx, rw, NewResponse(rate), http.StatusOK)
}
