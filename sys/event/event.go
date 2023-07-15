package event

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
)

// ErrInvalidEvent is an error indicating that the event is invalid.
var ErrInvalidEvent = errors.New("invalid event")

// Event represents an event in the system.
type Event struct {
	Source   Source      `json:"source"`
	Kind     Kind        `json:"kind"`
	Data     interface{} `json:"data"`
	Response chan Event  `json:"-"`
}

// Source represents the source of an event.
type Source = string

// Kind represents the type of event.
type Kind = string

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

// Register registers a listener for the specified source and kind.
func (b *Bus) Register(k Kind, l Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.log.Infow("register listener", "kind", k)

	b.lists[k] = append(b.lists[k], l)
}

// Dispatch dispatches an event to all listeners registered for its source and kind.
func (b *Bus) Dispatch(ctx context.Context, e Event) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.log.Infow("dispatch event", "status", "started", "event", e)
	defer b.log.Infow("dispatch event", "status", "completed", "event", e)

	var errs []error

	if lists, ok := b.lists[e.Kind]; ok {
		for _, l := range lists {
			if err := l(ctx, e); err != nil {
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
