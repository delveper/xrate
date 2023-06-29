package subscription

import (
	"errors"
	"fmt"
	"math/rand"
	"net/mail"
	"os"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/filestore"
	"github.com/stretchr/testify/require"
)

func TestRepoIntegration(t *testing.T) {
	t.Run("Get none", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		testGetAll(t, repo, os.ErrNotExist)
	})

	t.Run("Add and get many", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		emails := []Email{
			{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}},
			{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
			{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
		}

		testAdd(t, repo, nil, emails...)
		testGetAll(t, repo, nil, emails...)
	})

	t.Run("Add and get many, add duplicate and get many, add duplicate", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		emails := []Email{
			{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}},
			{Address: &mail.Address{Name: "Jane Smith", Address: "janesmith@example.com"}},
			{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}},
		}

		testAdd(t, repo, nil, emails...)
		testGetAll(t, repo, nil, emails...)

		testAdd(t, repo, os.ErrExist, emails[0])
		testGetAll(t, repo, nil, emails...)

		testAdd(t, repo, os.ErrExist, emails[0])
	})

	t.Run("Add one and get one, add duplicate and get one", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		email := Email{Address: &mail.Address{Name: "John Doe", Address: "johndoe@example.com"}}

		testAdd(t, repo, nil, email)
		testGetAll(t, repo, nil, email)

		testAdd(t, repo, os.ErrExist, email)
		testGetAll(t, repo, nil, email)
	})

	t.Run("Add and get whole lot, get whole lot again, add duplicates randomly", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		const wholeLot = 1_000

		emails := make([]Email, wholeLot)
		for i := range emails {
			emails[i] = Email{Address: &mail.Address{
				Name:    fmt.Sprintf("User%d", i),
				Address: fmt.Sprintf("user%d@example.com", i)},
			}
		}

		testAdd(t, repo, nil, emails...)
		testGetAll(t, repo, nil, emails...)

		testGetAll(t, repo, nil, emails...)

		for i := 0; i < rand.Intn(wholeLot); i++ { //nolint:gosec
			testAdd(t, repo, os.ErrExist, emails[i])
		}
	})
}

func testSetupRepo(t *testing.T) (*Repo, func()) {
	t.Helper()
	store, teardown := filestore.TestSetup[Email](t)

	return NewRepo(store), teardown
}

func testAdd(t *testing.T, repo *Repo, wantErr error, want ...Email) {
	t.Helper()

	var errArr []error
	for _, email := range want {
		errArr = append(errArr, repo.Store(email))
	}

	require.ErrorIs(t, errors.Join(errArr...), wantErr)
}

func testGetAll(t *testing.T, repo *Repo, wantErr error, want ...Email) {
	t.Helper()

	got, err := repo.GetAll()
	require.ErrorIs(t, err, wantErr)
	require.ElementsMatch(t, want, got)
}
