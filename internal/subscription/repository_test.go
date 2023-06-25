package subscription

import (
	"net/mail"
	"os"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/filestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoAddIntegration(t *testing.T) {
	tests := map[string]struct {
		got     []Email
		wantErr error
	}{
		"Add single email": {
			got: []Email{
				{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
			},
			wantErr: nil,
		},
		"Add multiple emails": {
			got: []Email{
				{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
				{Address: &mail.Address{Name: "Jason Johnson", Address: "jasonjohnson@example.com"}},
			},
			wantErr: nil,
		},
		"Add duplicate email": {
			got: []Email{
				{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
				{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
			},
			wantErr: ErrEmailAlreadyExists,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			store, teardown := filestore.TestSetup[Email](t)
			defer teardown()

			repo := NewRepo(store)

			var err error
			for _, item := range tt.got {
				err = repo.Add(item)
				if err != nil {
					break
				}
			}

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestRepoGetAllIntegration(t *testing.T) {
	tests := map[string]struct {
		want    []Email
		wantErr error
	}{
		"Fetch all emails": {
			want: []Email{
				{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
				{Address: &mail.Address{Name: "Jason Johnson", Address: "jasonjohnson@example.com"}},
				{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
			},
			wantErr: nil,
		},
		"Fetch with no emails": {
			want:    nil,
			wantErr: os.ErrNotExist,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			store, teardown := filestore.TestSetup[Email](t)
			defer teardown()

			repo := NewRepo(store)

			for _, item := range tt.want {
				err := repo.Add(item)
				require.NoError(t, err)
			}

			got, err := repo.GetAll()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
