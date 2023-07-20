package notif

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
)

const (
	EventSource        = "notif"
	EventKindRequested = "requested"
)

var ErrInvalidEvent = errors.New("invalid event")

type MetaData struct {
	xrt   *ExchangeRate
	subss []Subscription
}

func (svc *Service) RequestMetaData(ctx context.Context, pair CurrencyPair) (*MetaData, error) {
	e := event.New(EventSource, EventKindRequested, pair)

	if err := svc.bus.Publish(ctx, e); err != nil {
		return nil, fmt.Errorf("publishing sending email event: %w", err)
	}
	var (
		xrt   *ExchangeRate
		subss []Subscription
	)

	for xrt == nil || subss == nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case e := <-e.Response:
			switch val := e.Payload.(type) {
			case *ExchangeRate:
				xrt = val
			case []Subscription:
				subss = val
			default:
				return nil, fmt.Errorf("%w: unexpected payload: %T", ErrInvalidEvent, e.Payload)
			}
		}
	}

	return &MetaData{xrt: xrt, subss: subss}, nil
}
