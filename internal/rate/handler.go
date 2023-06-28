// Package rate provides functionality to retrieve and handle exchange rates.
package rate

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
)

// Getter interface to get rate from external service.
//
//go:generate moq -out getter_mock_test.go . Getter
type Getter interface {
	Get(context.Context) (float64, error)
}

// StatusError is used to communicate the error to the client.
const StatusError = "unexpected error"

const defaultTimeout = 5 * time.Second

type Response struct {
	Rate float64
}

type ResponseError struct {
	Error string
}

// Handler structure for handling rate requests.
type Handler struct {
	rate Getter
	log  *logger.Logger
}

// NewHandler creates a new Handler instance.
func NewHandler(rate Getter, log *logger.Logger) *Handler {
	return &Handler{
		rate: rate,
		log:  log,
	}
}

// Rate handles the HTTP request for the rate.
func (h *Handler) Rate(rw http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	rate, err := h.rate.Get(ctx)
	if err != nil {
		h.log.Errorw("Failed to get rate", "error", err)
		rw.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(rw).Encode(ResponseError{StatusError}); err != nil {
			h.log.Errorw("Writing response", "error", err)
		}

		return
	}

	if err := json.NewEncoder(rw).Encode(Response{rate}); err != nil {
		h.log.Errorw("Writing response", "error", err)
		rw.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(rw).Encode(ResponseError{StatusError}); err != nil {
			h.log.Errorw("Writing response", "error", err)
		}
	}
}
