package subscription

import (
	"encoding/json"
	"errors"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
	"net/mail"
	"strings"
)

//go:generate moq -out subscriber_mock_test.go . Subscriber
type Subscriber interface {
	Subscribe(Email) error
	SendEmails() error
}

type Handler struct {
	sub Subscriber
	log *logger.Logger
}

func NewHandler(sub Subscriber, log *logger.Logger) *Handler {
	return &Handler{
		sub: sub,
		log: log,
	}
}

func toEmail(addr *mail.Address) Email {
	return Email{
		Address: addr,
	}
}

// Subscribe subscribes to e-mails.
func (h *Handler) Subscribe(rw http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, "/") {
		req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
	}

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
func (h *Handler) SendEmails(rw http.ResponseWriter, _ *http.Request) {
	if err := h.sub.SendEmails(); err != nil {
		h.log.Errorw("Subscription failed", "error", err)
	}

	msg := struct{ Message string }{Message: "E-mail sent"}
	if err := json.NewEncoder(rw).Encode(msg); err != nil {
		h.log.Errorw("Writing response", "error", err)
	}

	h.log.Infow("E-mails sent")
}
