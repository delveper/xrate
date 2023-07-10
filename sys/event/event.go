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
	Source   Source       `json:"source"`
	Kind     Kind         `json:"kind"`
	Data     interface{}  `json:"data"`
	Response chan<- Event `json:"-"`
}

// Source represents the source of an event.
type Source string

// Kind represents the type of event.
type Kind string

// New creates a new event with the given source, topic, and data.
func New(source Source, topic Kind, data interface{}) Event {
	return Event{
		Source:   source,
		Kind:     topic,
		Data:     data,
		Response: make(chan Event, 1),
	}
}

// Bus is an event bus that handles event subscriptions
// and dispatches events to registered lists.
type Bus struct {
	log   *logger.Logger
	mu    sync.Mutex
	lists map[Event][]Listener
}

// Listener represents a function that handles an event.
type Listener func(context.Context, Event) error

// NewBus creates a new event bus with the given logger.
func NewBus(log *logger.Logger) *Bus {
	return &Bus{
		log:   log,
		lists: make(map[Event][]Listener),
	}
}

// Publish registers a listener for the specified event.
func (b *Bus) Publish(e Event, h Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lists[e] = append(b.lists[e], h)
}

func (b *Bus) Subscribe(l Listener, events ...Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i := range events {
		b.lists[events[i]] = append(b.lists[events[i]], l)
	}
}

func (b *Bus) Dispatch(ctx context.Context, e Event) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.log.Infow("dispatch event", "status", "started", "event", e)
	defer b.log.Infow("dispatch event", "status", "completed", "event", e)

	var errArr []error

	lists := b.lists[e]

	for i := range lists {
		if err := lists[i](ctx, e); err != nil {
			b.log.Errorw("dispatch event", "err", err)
			errArr = append(errArr, err)
		}
	}

	if err := errors.Join(errArr...); err != nil {
		return fmt.Errorf("dispatch event: %w", err)
	}

	return nil
}
