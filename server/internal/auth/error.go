package auth

import (
	"errors"
)

var (
	// errors that are specific to the auth domain

	ErrEmailNotAvailable = errors.New("email not available")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrIncorrectPassword  = errors.New("incorrect password")

	ErrTokenGenerationFailed  = errors.New("failed to generate token")
	ErrTokenExpired           = errors.New("token expired")
	ErrInvalidToken           = errors.New("invalid token")
	ErrAuthenticationRequired = errors.New("authentication required")

	ErrAccessDenied = errors.New("access denied")
)
