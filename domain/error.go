package domain

import "errors"

var (
	ErrBadRequest          = errors.New("bad request")
	ErrForbidenAccess      = errors.New("forbidden access")
	ErrNotFound            = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")

	ErrInvalidCredential = errors.New("invalid credential")
	ErrInvalidRole       = errors.New("invalid role")
	ErrInvalidId         = errors.New("invalid id")
	ErrDuplicateUser     = errors.New("duplicate entry")
)
