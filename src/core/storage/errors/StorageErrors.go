package errors

import (
	"errors"
)

var (
	// ErrNotFound occurs when the storage resource requested was not found
	ErrNotFound = errors.New("Requested resource not found")
)
