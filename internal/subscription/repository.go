package subscription

import (
	"errors"
	"io/fs"
)

type Storer interface {
	Store(name string, item Email) error
	FetchAll() ([]Email, error)
}

type Repo struct{ Storer }

func NewRepo(fileStore Storer) *Repo {
	return &Repo{fileStore}
}

func (s *Repo) Add(email Email) error {
	if err := s.Storer.Store(email.Address.String(), email); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

func (s *Repo) GetAll() ([]Email, error) {
	return s.Storer.FetchAll()
}
