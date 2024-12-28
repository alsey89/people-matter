package identity

import (
	"errors"
)

var (
	// errors that are specific to the auth domain
	ErrUserNotFound             = errors.New("user not found")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrEmailAlreadyInUse        = errors.New("email already in use")
	ErrEmailNotConfirmed        = errors.New("user not confirmed")
	ErrEmailAlreadyConfirmed    = errors.New("user already confirmed")
	ErrNewPasswordIsOldPassword = errors.New("new password is the same as the old password")

	ErrRoleNotFound     = errors.New("no roles found")
	ErrInvalidRoleLevel = errors.New("invalid role level")

	ErrUserHasFSPRole     = errors.New("user already has Tenant role")
	ErrUserHasApplication = errors.New("user already has application")

	ErrMemorialNotFound    = errors.New("memorial not found")
	ErrUserHasMemorialRole = errors.New("user already has a memorial role in this memorial")
)

var (
	ErrCodeUserNotFound             = "ERR_USER_NOT_FOUND"
	ErrCodeInvalidCredentials       = "ERR_INVALID_CREDENTIALS"
	ErrCodeEmailAlreadyInUse        = "ERR_EMAIL_ALREADY_IN_USE"
	ErrCodeEmailNotConfirmed        = "ERR_EMAIL_NOT_CONFIRMED"
	ErrCodeEmailAlreadyConfirmed    = "ERR_EMAIL_ALREADY_CONFIRMED"
	ErrCodeNewPasswordIsOldPassword = "ERR_NEW_PASSWORD_IS_OLD_PASSWORD"

	ErrCodeRoleNotFound     = "ERR_ROLE_NOT_FOUND"
	ErrCodeInvalidRoleLevel = "ERR_INVALID_ROLE_LEVEL"

	ErrCodeUserHasFSPRole     = "ERR_USER_HAS_FSP_ROLE"
	ErrCodeUserHasApplication = "ERR_USER_HAS_APPLICATION"

	ErrCodeMemorialNotFound    = "ERR_MEMORIAL_NOT_FOUND"
	ErrCodeUserHasMemorialRole = "ERR_USER_HAS_MEMORIAL_ROLE"
)
