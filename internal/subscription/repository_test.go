package subscription

import (
	"net/mail"
	"os"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/filestore"
	"github.com/stretchr/testify/require"
)

func TestRepoIntegration(t *testing.T) {
	tests := map[string]struct {
		setup   func(*Repo)
		actions func(*testing.T, *Repo)
	}{
		"Add, add, get all, add, get all": {
			setup: func(repo *Repo) {
				testAdd(t, repo, nil,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
					Email{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
					Email{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}},
				)
			},
			actions: func(t *testing.T, repo *Repo) {
				t.Helper()
				testGetAll(t, repo, nil,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
					Email{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
					Email{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}},
				)
				testAdd(t, repo, os.ErrExist,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}})

				testGetAll(t, repo, nil,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
					Email{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
					Email{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}})
			},
		},
		"Add, get all, add duplicate, get all": {
			setup: func(repo *Repo) {
				testAdd(t, repo, nil,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
				)
			},
			actions: func(t *testing.T, repo *Repo) {
				t.Helper()
				testGetAll(t, repo, nil,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
				)

				testAdd(t, repo, os.ErrExist,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
				)

				testGetAll(t, repo, nil,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
				)
			},
		},
		"Get none": {
			setup: func(repo *Repo) {},
			actions: func(t *testing.T, repo *Repo) {
				t.Helper()
				var want []Email
				testGetAll(t, repo, os.ErrNotExist, want...)
			},
		},
		"Add, get all, reorder": {
			setup: func(repo *Repo) {},
			actions: func(t *testing.T, repo *Repo) {
				t.Helper()

				testAdd(t, repo, nil,
					Email{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}},
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
					Email{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
				)

				testGetAll(t, repo, nil,
					Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
					Email{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}},
					Email{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
				)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			store, teardown := filestore.TestSetup[Email](t)
			defer teardown()

			repo := NewRepo(store)
			tt.setup(repo)
			tt.actions(t, repo)
		})
	}
}
func testAdd(t *testing.T, repo *Repo, wantErr error, want ...Email) {
	t.Helper()

	var err error
	for _, email := range want {
		err = repo.Store(email)
	}

	require.Equal(t, err, wantErr)
}

func testGetAll(t *testing.T, repo *Repo, wantErr error, want ...Email) {
	t.Helper()

	got, err := repo.GetAll()
	require.ErrorIs(t, err, wantErr)

	require.ElementsMatch(t, want, got)
}
