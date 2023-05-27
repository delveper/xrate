package subscription

import (
	"errors"
	"io/fs"
)

type Storer interface {
	Store(name string, item Email) error
	FetchAll() ([]Email, error)
}

type Store struct{ Storer }

func NewStore(fileStore Storer) *Store {
	return &Store{fileStore}
}

func (s *Store) Add(email Email) error {
	if err := s.Storer.Store(email.Address.String(), email); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

func (s *Store) GetAll() ([]Email, error) {
	return s.Storer.FetchAll()
}
