package filestore

import (
	"os"
	"testing"

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
			wantErr: os.ErrExist,
		},
		"List valid items": {
			items:   []item{{Name: "item1", Value: 1}, {Name: "item2", Value: 2}, {Name: "item3", Value: 3}},
			wantErr: nil,
		},
		"Invalid item": {
			items:   []item{{}},
			wantErr: os.ErrInvalid,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store, teardown := TestSetup[item](t)
			defer teardown()

			var err error
			for _, item := range tc.items {
				err = store.Store(item)
				if err != nil {
					break
				}
			}

			require.ErrorIs(t, err, tc.wantErr)
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
			wantErr: os.ErrNotExist,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store, teardown := TestSetup[item](t)
			defer teardown()

			for _, item := range tc.got {
				err := store.Store(item)
				require.NoError(t, err)
			}

			fetchedItems, err := store.FetchAll()
			require.ErrorIs(t, err, tc.wantErr)

			require.ElementsMatch(t, tc.want, fetchedItems)
		})
	}
}
