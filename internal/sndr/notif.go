/*
Package sndr provides functionality to sending messages.
It designed for sending email notifications and open to any other use cases.
*/
package sndr

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

type Creator interface {
	Create(*MetaData) (*Message, error)
}

// Sender is an interface for sending messages.
type Sender interface {
	Send(context.Context, *Message) (err error)
}

type Service struct {
	bus  *event.Bus
	mail Sender
	crt  Creator
}

func NewService(bus *event.Bus, mail Sender, crt Creator) *Service {
	return &Service{bus: bus, mail: mail, crt: crt}
}

func (svc *Service) SendEmails(ctx context.Context, topic Topic) error {
	md, err := svc.RequestMetaData(ctx, topic)
	if err != nil {
		return err
	}

	msg, err := svc.crt.Create(md)
	if err != nil {
		return err
	}

	return svc.mail.Send(ctx, msg)
}
