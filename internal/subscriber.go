package internal

import (
	"github.com/google/uuid"
	"net/mail"
	"time"
)

type Subscriber struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address
	SubscribedAt time.Time
}

func NewSubscriber(name string, email mail.Address) *Subscriber {
	return &Subscriber{
		ID:    uuid.New(),
		Name:  name,
		Email: email,
	}
}
