package common

import (
	"errors"
)

var (
	// errors that are common to all services
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidInput     = errors.New("invalid input")
	ErrResourceNotFound = errors.New("resource not found")
)
