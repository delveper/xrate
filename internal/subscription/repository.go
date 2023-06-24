package subscription

import (
	"errors"
	"fmt"
	"io/fs"
)

// Storer defines the interface for storing and retrieving email subscriptions.
type Storer interface {
	Store(Email) error
	FetchAll() ([]Email, error)
}

// Repo is a repository that implements the Storer interface.
type Repo struct{ Storer }

// NewRepo creates a new Repo instance with the provided Storer implementation.
func NewRepo(fileStore Storer) *Repo {
	return &Repo{Storer: fileStore}
}

// Add creates a new email subscription.
func (s *Repo) Add(email Email) error {
	if err := s.Storer.Store(email); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return ErrEmailAlreadyExists
		}

		return fmt.Errorf("adding email subscription: %w", err)
	}

	return nil
}

// GetAll retrieves all email subscriptions from the repository.
func (s *Repo) GetAll() ([]Email, error) {
	emails, err := s.Storer.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("getting all email subscriptions: %w", err)
	}

	return emails, nil
}
