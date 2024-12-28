package fsp

import "errors"

var (
	// errors that are specific to the domain

	ErrTeamMemberHasRole      = errors.New("user already has the specified role")
	ErrTeamMemberHasAnAccount = errors.New("user already has an account, updated role instead")
	ErrMemorialExists         = errors.New("memorial already exists")
	ErrUserIsLastSuperAdmin   = errors.New("cannot remove the last super admin")
)
