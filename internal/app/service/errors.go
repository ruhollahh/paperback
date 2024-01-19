package service

import "errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrBadRequest     = errors.New("bad request")
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)
