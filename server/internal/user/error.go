package user

import (
	"errors"
)

var (
	// errors that are specific to the user domain

	ErrUserNotFound = errors.New("user not found")

	ErrUserAlreadyExists = errors.New("user already exists")

	ErrUserCreateFailed = errors.New("failed to create user data")
	ErrUserReadFailed   = errors.New("failed to update read user data")
	ErrUserUpdateFailed = errors.New("failed to update user data")
	ErrUserDeleteFailed = errors.New("failed to delete user data")
)
