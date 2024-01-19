package domain

import (
	"errors"
	"github.com/ruhollahh/paperback/pkg/validation"
	"time"
)

type User struct {
	ID             int64
	CreatedAt      time.Time
	Name           string
	Email          string
	HashedPassword []byte
	Activated      bool
	Version        int32
}

func NewName(name string) (string, error) {
	if name == "" {
		return "", errors.New("must be provided")
	}
	if len(name) > 500 {
		return "", errors.New("must not be more than 500 bytes long")
	}
	return name, nil
}

func NewEmail(email string) (string, error) {
	if email == "" {
		return "", errors.New("must be provided")
	}
	if !validation.Matches(email, validation.EmailRX) {
		return "", errors.New("must be a valid email address")
	}
	return email, nil
}

func NewPasswordPlaintext(password string) (string, error) {
	if password == "" {
		return "", errors.New("must be provided")
	}

	if len(password) < 8 {
		return "", errors.New("must be at least 8 bytes long")
	}
	if len(password) > 72 {
		return "", errors.New("must not be more than 72 bytes long")
	}
	return password, nil
}
