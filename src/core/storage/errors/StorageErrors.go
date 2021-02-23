package errors

import (
	"errors"
)

var (
	// ErrNotFound occurs when the storage resource requested was not found
	ErrNotFound = errors.New("Requested resource not found")

	// ErrType occurs when an unknown type is given to a repository function
	ErrType = errors.New("Unexpected type given")
)
