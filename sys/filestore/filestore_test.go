package filestore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	type item struct {
		Name  string
		Value int
	}

	tests := map[string]struct {
		name    string
		items   []item
		dir     string
		wantErr error
	}{
		"Valid item": {
			items:   []item{{Name: "item1", Value: 1}},
			wantErr: nil,
		},
		"Duplicate item": {
			items:   []item{{Name: "item1", Value: 1}, {Name: "item1", Value: 1}},
			wantErr: ErrItemExists,
		},
		"List valid items": {
			items:   []item{{Name: "item1", Value: 1}, {Name: "item2", Value: 2}, {Name: "item3", Value: 3}},
			wantErr: nil,
		},
		"Invalid item": {
			items:   []item{{}},
			wantErr: ErrInvalidItem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, teardown := TestSetup[item](t)
			defer teardown()

			var err error
			for _, item := range tt.items {
				err = store.Store(item)
				if err != nil {
					break
				}
			}

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestFetchAll(t *testing.T) {
	type item struct {
		Name  string
		Value int
	}

	tests := map[string]struct {
		name    string
		got     []item
		want    []item
		wantErr error
	}{
		"Fetch single item": {
			got:     []item{{Name: "item1", Value: 1}},
			want:    []item{{Name: "item1", Value: 1}},
			wantErr: nil,
		},
		"Fetch multiple items": {
			got:     []item{{Name: "item1", Value: 1}, {Name: "item2", Value: 2}, {Name: "item3", Value: 3}},
			want:    []item{{Name: "item1", Value: 1}, {Name: "item2", Value: 2}, {Name: "item3", Value: 3}},
			wantErr: nil,
		},
		"Error fetching items": {
			wantErr: ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, teardown := TestSetup[item](t)
			defer teardown()

			for _, item := range tt.got {
				err := store.Store(item)
				require.NoError(t, err)
			}

			fetchedItems, err := store.FetchAll()
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, fetchedItems)
		})
	}
}
