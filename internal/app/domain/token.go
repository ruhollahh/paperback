package domain

import (
	"errors"
	"time"
)

type TokenScope string

const (
	ActivationScope TokenScope = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     TokenScope
}

func NewTokenPlaintext(tokenPlaintext string) (string, error) {
	if tokenPlaintext == "" {
		return "", errors.New("must be provided")
	}

	if len(tokenPlaintext) != 26 {
		return "", errors.New("invalid or expired token")
	}

	return tokenPlaintext, nil
}
