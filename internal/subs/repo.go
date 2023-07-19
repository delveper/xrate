package subs

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// Storer defines the interface for storing and retrieving subscribers.
type Storer interface {
	Store(Subscription) error
	FetchAll() ([]Subscription, error)
}

// Repo is a repository that implements the Storer interface.
type Repo struct{ Storer }

// NewRepo creates a new Repo instance with the provided Storer implementation.
func NewRepo(fileStore Storer) *Repo {
	return &Repo{Storer: fileStore}
}

// Add creates a new email subscription.
func (r *Repo) Add(ctx context.Context, subs Subscription) error {
	if err := r.Storer.Store(subs); err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrSubscriptionExists
		}

		return fmt.Errorf("adding subscription: %w", err)
	}

	return nil
}

// List retrieves all email subscriptions from the repository.
func (r *Repo) List(ctx context.Context) ([]Subscription, error) {
	subss, err := r.Storer.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("getting all subscriptions: %w", err)
	}

	return subss, nil
}
