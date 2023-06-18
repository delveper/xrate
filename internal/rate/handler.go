package rate

import (
	"encoding/json"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
)

// Getter interface to get rate from external service.
//
//go:generate moq -out getter_mock_test.go . Getter
type Getter interface {
	Get() (float64, error)
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
	rate, err := h.rate.Get()
	if err != nil {
		h.log.Errorw("Failed to get rate", "error", err)
		rw.WriteHeader(http.StatusBadRequest)

		resp := struct{ Error string }{Error: "Failed to get rate"}
		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			h.log.Errorw("Writing response", "error", err)
		}

		return
	}

	resp := struct{ Rate float64 }{Rate: rate}
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		h.log.Errorw("Writing response", "error", err)
		rw.WriteHeader(http.StatusBadRequest)

		resp := struct{ Error string }{Error: "Failed to write response"}
		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			h.log.Errorw("Writing response", "error", err)
		}
	}
}
