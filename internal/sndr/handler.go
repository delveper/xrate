package sndr

import (
	"context"
	"net/http"
	"time"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/web"
)

const defaultTimeout = 15 * time.Second

const StatusSend = "emails sent"

type EmailService interface {
	SendEmails(context.Context, Topic) (err error)
}

// Handler handles subscription.
type Handler struct {
	EmailService
}

type Response struct {
	Status string
}

func NewResponse(status string) *Response {
	return &Response{Status: status}
}

func NewHandler(svc EmailService) *Handler {
	return &Handler{EmailService: svc}
}

// SendEmails sends all e-mails stored in data base.
func (h *Handler) SendEmails(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var topic Topic
	if err := web.DecodeBody(req.Body, &topic); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	if err := h.EmailService.SendEmails(ctx, topic); err != nil {
		return err
	}

	return web.Respond(ctx, rw, NewResponse(StatusSend), http.StatusOK)
}
