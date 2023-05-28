package subscription

import (
	"errors"
	"net/mail"
)

var ErrEmailAlreadyExists = errors.New("email address is already in the database")

type Email struct {
	Address *mail.Address
}

type EmailRepository interface {
	Add(Email) error
	GetAll() ([]Email, error)
}

type RateGetter interface {
	Get() (float64, error)
}

type EmailSender interface {
	Send(Email, float64) error
}

type Service struct {
	repo EmailRepository
	rate RateGetter
	mail EmailSender
}

func NewService(repo EmailRepository, rate RateGetter, mail EmailSender) *Service {
	return &Service{
		repo: repo,
		rate: rate,
		mail: mail,
	}
}

func (svc *Service) Subscribe(email Email) error {
	if err := svc.repo.Add(email); err != nil {
		return err
	}

	return nil
}

func (svc *Service) SendEmails() error {
	rate, err := svc.rate.Get()
	if err != nil {
		return err
	}

	emails, err := svc.repo.GetAll()
	if err != nil {
		return err
	}

	var errArr []error
	for _, email := range emails {
		err := svc.mail.Send(email, rate)
		if err != nil {
			errArr = append(errArr, err)
		}
	}

	if errArr != nil {
		return errors.Join(errArr...)
	}

	return nil
}
