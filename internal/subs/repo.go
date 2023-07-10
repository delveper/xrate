package subs

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// Storer defines the interface for storing and retrieving subscribers.
type Storer interface {
	// Store TODO: Publish context.
	Store(Subscriber) error
	// FetchAll TODO: Publish context.
	FetchAll() ([]Subscriber, error)
}

// Repo is a repository that implements the Storer interface.
type Repo struct{ Storer }

// NewRepo creates a new Repo instance with the provided Storer implementation.
func NewRepo(fileStore Storer) *Repo {
	return &Repo{Storer: fileStore}
}

// Add creates a new email subscription.
func (s *Repo) Add(ctx context.Context, email Subscriber) error {
	if err := s.Storer.Store(email); err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrEmailAlreadyExists
		}

		return fmt.Errorf("adding email subscription: %w", err)
	}

	return nil
}

// List retrieves all email subscriptions from the repository.
func (s *Repo) List(ctx context.Context) ([]Subscriber, error) {
	emails, err := s.Storer.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("getting all email subscriptions: %w", err)
	}

	return emails, nil
}
