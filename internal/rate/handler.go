package rate

import (
	"encoding/json"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
)

//go:generate moq -out getter_mock_test.go . Getter
type Getter interface {
	Get() (float64, error)
}

type Handler struct {
	rate Getter
	log  *logger.Logger
}

func NewHandler(rate Getter, log *logger.Logger) *Handler {
	return &Handler{
		rate: rate,
		log:  log,
	}
}

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
