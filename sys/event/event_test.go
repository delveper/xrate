package event_test

import (
	"context"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/event"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	log := logger.New(logger.WithConsoleCore(logger.LevelDebug))
	defer log.Sync()

	bus := event.NewBus(log)

	tests := map[string]struct {
		got     event.Event
		want    event.Event
		list    func(ctx context.Context, e event.Event) error
		wantErr error
	}{
		"Simple event": {
			got: event.Event{
				Kind:    "test",
				Payload: nil,
			},
			want: event.Event{
				Kind:    "test",
				Payload: nil,
			},
			list: func(ctx context.Context, e event.Event) error {
				require.Equal(t, "test", e.Kind)
				require.Nil(t, e.Payload)
				return nil
			},
			wantErr: nil,
		},
		"Event with payload": {
			got: event.Event{
				Kind:    "test",
				Payload: "payload",
			},
			want: event.Event{
				Kind:    "test",
				Payload: "payload",
			},
			list: func(ctx context.Context, e event.Event) error {
				require.Equal(t, "test", e.Kind)
				require.Equal(t, "payload", e.Payload)
				return nil
			},
			wantErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			done := make(chan bool)
			bus.Subscribe(tc.got, func(ctx context.Context, e event.Event) error {
				defer func() { done <- true }()
				return tc.list(ctx, e)
			})

			err := bus.Publish(context.Background(), tc.got)
			require.Equal(t, tc.wantErr, err)
		})
	}
}
