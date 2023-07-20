/*
Package notif provides functionality to sending messages.
It designed for sending email notifications and open to any other use cases.
*/
package notif

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

type MessageCreator interface {
	CreateMessage(*MetaData) (*Message, error)
}

// Sender is an interface for sending messages.
type Sender interface {
	Send(context.Context, *Message) (err error)
}

type Service struct {
	bus  *event.Bus
	sndr Sender
	mc   MessageCreator
}

func NewService(bus *event.Bus, mail Sender, crt MessageCreator) *Service {
	return &Service{bus: bus, sndr: mail, mc: crt}
}

func (svc *Service) SendEmails(ctx context.Context, topic Topic) error {
	md, err := svc.RequestMetaData(ctx, topic)
	if err != nil {
		return err
	}

	msg, err := svc.mc.CreateMessage(md)
	if err != nil {
		return err
	}

	return svc.sndr.Send(ctx, msg)
}
