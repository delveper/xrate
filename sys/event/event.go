package event

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
)

// Event represents an event in the system.
type Event struct {
	Source   Source      `json:"source"`
	Kind     Kind        `json:"kind"`
	Payload  interface{} `json:"payload,omitempty"`
	Response chan Event  `json:"-"`
}

// Source represents the source of an event.
type Source = string

// Kind represents the type of event.
type Kind = string

// New creates a new event with the given source, topic, and data.
func New(source Source, kind Kind, payload interface{}) Event {
	return Event{
		Source:   source,
		Kind:     kind,
		Payload:  payload,
		Response: make(chan Event, 2),
	}
}

// Bus is an event bus that handles event subscriptions
// and dispatches events to registered lists.
type Bus struct {
	log   *logger.Logger
	mu    sync.Mutex
	lists map[Kind][]Listener
}

// Listener represents a function that handles an event.
type Listener func(context.Context, Event) error

// NewBus creates a new event bus with the given logger.
func NewBus(log *logger.Logger) *Bus {
	return &Bus{
		log:   log,
		lists: make(map[Kind][]Listener),
	}
}

// Subscribe registers a listener for the specified source and kind.
func (b *Bus) Subscribe(e Event, l Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.log.Infow("subscribe listener", "event", e)
	b.lists[e.Kind] = append(b.lists[e.Kind], l)
}

// Publish dispatches an event to all listeners registered for its source and kind.
func (b *Bus) Publish(ctx context.Context, e Event) error {
	b.log.Infow("dispatch event", "status", "started", "event", e)
	defer b.log.Infow("dispatch event", "status", "completed", "event", e)

	var errs []error

	if lists, ok := b.lists[e.Kind]; ok {
		for i := range lists {
			if err := lists[i](ctx, e); err != nil {
				b.log.Errorw("dispatch event", "err", err)
				errs = append(errs, err)
			}
		}
	}

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("dispatch event: %w", err)
	}

	return nil
}
