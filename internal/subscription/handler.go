package subscription

import (
	"encoding/json"
	"errors"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
	"net/mail"
	"strings"
)

const (
	StatusSend       = "emails sent"
	StatusSubscribed = "subscribed"
	StatusError      = "unexpected error"
)

// Subscriber is an interface for subscription service.
//
//go:generate moq -out subscriber_mock_test.go . Subscriber
type Subscriber interface {
	Subscribe(Email) error
	SendEmails() error
}

// Response is a response for subscription service.
type Response struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Handler handles subscription.
type Handler struct {
	sub Subscriber
	log *logger.Logger
}

// NewHandler returns a new Handler instance.
func NewHandler(sub Subscriber, log *logger.Logger) *Handler {
	return &Handler{
		sub: sub,
		log: log,
	}
}

// toEmail converts mail.Address to Email.
func toEmail(addr *mail.Address) Email {
	return Email{
		Address: addr,
	}
}

// Subscribe subscribes to e-mails.
func (h *Handler) Subscribe(rw http.ResponseWriter, req *http.Request) {
	req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")

	addr := req.FormValue("email")

	email, err := mail.ParseAddress(addr)
	if err != nil {
		h.log.Errorw("Invalid email", "error", err)
		// TODO(not_documented): Handle invalid email.
		return
	}

	if err := h.sub.Subscribe(toEmail(email)); err != nil {
		h.log.Errorw("Subscription failed", "error", err)

		if errors.Is(err, ErrEmailAlreadyExists) {
			rw.WriteHeader(http.StatusConflict)

			resp := Response{Message: StatusError, Details: ErrEmailAlreadyExists.Error()}
			if err := json.NewEncoder(rw).Encode(resp); err != nil {
				h.log.Errorw("Writing response", "error", err)
			}
			// TODO(not_documented): Handle other errors.
			return
		}
	}

	resp := Response{Message: StatusSubscribed}
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		h.log.Errorw("Writing response", "error", err)
	}
}

// SendEmails sends all e-mails stored in data base.
func (h *Handler) SendEmails(rw http.ResponseWriter, _ *http.Request) {
	if err := h.sub.SendEmails(); err != nil {
		h.log.Errorw("Sending failed", "error", err)
		// TODO(not_documented): Handle other errors.
		return
	}

	resp := Response{Message: StatusSend}
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		h.log.Errorw("Writing response", "error", err)
		// TODO(not_documented): Handle encoding error.
		return
	}

	h.log.Infow("E-mails sent")
}
