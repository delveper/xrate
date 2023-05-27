package subscription

import (
	"encoding/json"
	"errors"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
	"net/mail"
)

type Subscriber interface {
	Subscribe(Email) error
	SendEmails() error
}

type Handler struct {
	sub Subscriber
	log logger.Logger
}

func NewHandler(sub Subscriber, log logger.Logger) (*Handler, error) {
	if sub == nil {
		return nil, errors.New("subscriber must not be nil")
	}

	h := Handler{
		sub: sub,
		log: log,
	}

	return &h, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		switch req.URL.Path {
		case "/subscribe":
			h.Subscribe(rw, req)

		case "/sendEmails":
			h.SendEmails(rw, req)
		}
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		resp := struct{ Error string }{Error: http.StatusText(http.StatusMethodNotAllowed)}
		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			h.log.Errorw("Writing response", "error", err)
		}
	}
}

func toEmail(addr *mail.Address) Email {
	return Email{
		Address: *addr,
	}
}

// Subscribe subscribes to e-mails.
func (h *Handler) Subscribe(rw http.ResponseWriter, req *http.Request) {
	addr := req.FormValue("email")

	email, err := mail.ParseAddress(addr)
	if err != nil {
		h.log.Errorw("Invalid email", "error", err)
	}

	if err := h.sub.Subscribe(toEmail(email)); err != nil {
		h.log.Errorw("Subscription failed", "error", err)
		switch {
		case errors.Is(err, ErrEmailAlreadyExists):
			rw.WriteHeader(http.StatusConflict)
			resp := struct{ Error string }{Error: ErrEmailAlreadyExists.Error()}
			if err := json.NewEncoder(rw).Encode(resp); err != nil {
				h.log.Errorw("Writing response", "error", err)
			}
			return
		}
	}

	msg := struct{ Message string }{Message: "E-mail subscribed"}
	if err := json.NewEncoder(rw).Encode(msg); err != nil {
		h.log.Errorw("Writing response", "error", err)
	}
}

// SendEmails sends all e-mails stored in data base.
func (h *Handler) SendEmails(rw http.ResponseWriter, req *http.Request) {
	if err := h.sub.SendEmails(); err != nil {
		h.log.Errorw("Subscription failed", "error", err)
	}

	msg := struct{ Message string }{Message: "E-mail sent"}
	if err := json.NewEncoder(rw).Encode(msg); err != nil {
		h.log.Errorw("Writing response", "error", err)
	}

	h.log.Infow("E-mails sent")
}
