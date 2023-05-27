package rate

import (
	"encoding/json"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
)

type Getter interface {
	GetRate() (float64, error)
}

type Handler struct {
	rate Getter
	log  logger.Logger
}

func NewHandler(rate Getter, log logger.Logger) *Handler {
	return &Handler{
		rate: rate,
		log:  log,
	}
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.Rate(rw, req)

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		resp := struct{ Error string }{Error: http.StatusText(http.StatusMethodNotAllowed)}
		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			h.log.Errorw("Writing response", "error", err)
		}
	}
}

func (h *Handler) Rate(rw http.ResponseWriter, _ *http.Request) {
	rate, err := h.rate.GetRate()
	if err != nil {
		h.log.Errorw("Failed to get rate", "error", err) // 500
	}

	resp := struct{ Rate float64 }{Rate: rate}
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		h.log.Errorw("Writing response", "error", err) // 500
	}
}
