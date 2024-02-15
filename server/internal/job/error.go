package job

import (
	"errors"
)

var (
	ErrEmptyTable  = errors.New("no entries in the table")
	ErrNoRowsFound = errors.New("no rows found")
)
