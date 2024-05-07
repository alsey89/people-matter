package company

import (
	"errors"
)

var (
	ErrUserExists  = errors.New("user already exists")
	ErrEmptyTable  = errors.New("no entries in the table")
	ErrNoRowsFound = errors.New("no rows found")
)
