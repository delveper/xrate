package subs

import (
	"context"
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
	t.Run("GetExchangeRate none", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		testGetAll(t, repo, os.ErrNotExist)
	})

	t.Run("Subscribe and get many", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		emails := []Subscription{
			{Subscriber{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}}, Topic{"BTC", "UAH"}},
			{Subscriber{Address: &mail.Address{Name: "Jon Doe", Address: "johndoe@example.com"}}, Topic{"BTC", "UAH"}},
			{Subscriber{Address: &mail.Address{Name: "Jane Smith", Address: "anesmith@example.com"}}, Topic{"BTC", "UAH"}},
		}

		testAdd(t, repo, nil, emails...)
		testGetAll(t, repo, nil, emails...)
	})

	t.Run("Subscribe and get many, add duplicate and get many, add duplicate", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		emails := []Subscription{
			{Subscriber{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}}, Topic{"BTC", "UAH"}},
			{Subscriber{Address: &mail.Address{Name: "Jon Doe", Address: "johndoe@example.com"}}, Topic{"BTC", "UAH"}},
			{Subscriber{Address: &mail.Address{Name: "Jane Smith", Address: "anesmith@example.com"}}, Topic{"BTC", "UAH"}},
		}

		testAdd(t, repo, nil, emails...)
		testGetAll(t, repo, nil, emails...)

		testAdd(t, repo, os.ErrExist, emails[0])
		testGetAll(t, repo, nil, emails...)

		testAdd(t, repo, os.ErrExist, emails[0])
	})

	t.Run("Subscribe one and get one, add duplicate and get one", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		subs := Subscription{Subscriber{Address: &mail.Address{Name: "Sam Johns", Address: "samjohns@example.com"}}, Topic{"BTC", "UAH"}}

		testAdd(t, repo, nil, subs)
		testGetAll(t, repo, nil, subs)

		testAdd(t, repo, os.ErrExist, subs)
		testGetAll(t, repo, nil, subs)
	})

	t.Run("Subscribe and get whole lot, get whole lot again, add duplicates randomly", func(t *testing.T) {
		repo, teardown := testSetupRepo(t)
		defer teardown()

		const wholeLot = 1_000

		subss := make([]Subscription, wholeLot)
		for i := range subss {
			subss[i] = Subscription{
				Subscriber{Address: &mail.Address{Name: fmt.Sprintf("User%d", i), Address: fmt.Sprintf("user%d@example.com", i)}},
				Topic{"BTC", "UAH"}}
		}

		testAdd(t, repo, nil, subss...)
		testGetAll(t, repo, nil, subss...)

		testGetAll(t, repo, nil, subss...)

		for i := 0; i < rand.Intn(wholeLot); i++ { //nolint:gosec
			testAdd(t, repo, os.ErrExist, subss[i])
		}
	})
}

func testSetupRepo(t *testing.T) (*Repo, func()) {
	t.Helper()
	store, teardown := filestore.TestSetup[Subscription](t)

	return NewRepo(store), teardown
}

func testAdd(t *testing.T, repo *Repo, wantErr error, want ...Subscription) {
	t.Helper()

	var errArr []error
	for _, email := range want {
		errArr = append(errArr, repo.Store(email))
	}

	require.ErrorIs(t, errors.Join(errArr...), wantErr)
}

func testGetAll(t *testing.T, repo *Repo, wantErr error, want ...Subscription) {
	t.Helper()

	got, err := repo.List(context.Background())
	require.ErrorIs(t, err, wantErr)
	require.ElementsMatch(t, want, got)
}
