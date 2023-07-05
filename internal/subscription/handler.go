package subscription

import (
	"context"
	"errors"
	"net/http"
	"net/mail"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

const (
	StatusSend       = "emails sent"
	StatusSubscribed = "subscribed"
)

//go:generate moq -out=../../test/mock/subscriber.go -pkg=mock . SubscriptionService

// SubscriptionService is an interface for subscription service.
type SubscriptionService interface {
	Subscribe(Subscriber) error
	SendEmails() error
}

// Response is a response for subscription service.
type Response struct {
	Message string `json:"message"`
}

func NewResponse(msg string) *Response {
	return &Response{Message: msg}
}

// Handler handles subscription.
type Handler struct {
	SubscriptionService
}

// NewHandler returns a new Handler instance.
func NewHandler(ss SubscriptionService) *Handler {
	return &Handler{SubscriptionService: ss}
}

// toSubscriber converts mail.Address to Subscriber.
func toSubscriber(addr *mail.Address, topic Topic) Subscriber {
	return Subscriber{Address: addr, Topic: topic}
}

// Subscribe subscribes to e-mails.
func (h *Handler) Subscribe(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	addr := web.FromQuery(req, "email")
	if addr == "" {
		return web.NewRequestError(ErrMissingEmail, http.StatusBadRequest)
	}

	email, err := mail.ParseAddress(addr)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	sub := toSubscriber(email, "BTC/UAH")
	if err := h.SubscriptionService.Subscribe(sub); err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return web.NewRequestError(err, http.StatusConflict)
		}

		return err
	}

	return web.Respond(ctx, rw, NewResponse(StatusSubscribed), http.StatusCreated)
}

// SendEmails sends all e-mails stored in data base.
func (h *Handler) SendEmails(ctx context.Context, rw http.ResponseWriter, _ *http.Request) error {
	if err := h.SubscriptionService.SendEmails(); err != nil {
		return err
	}

	return web.Respond(ctx, rw, NewResponse(StatusSend), http.StatusOK)
}
