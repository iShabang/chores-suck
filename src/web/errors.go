package web

import (
	"errors"
	"net/http"
)

var (
	// ErrSessionType describes a type assertion error for retrieving session values
	ErrSessionType = errors.New("session value type assertion failed")

	// ErrNotAuthorized occurs when the current session is not authorized
	ErrNotAuthorized = errors.New("authorization unsuccessfull")

	// ErrInvalidValue occurs when an invalid session value is retrieved from the current session
	ErrInvalidValue = errors.New("invalid session value")

	// ErrInvalidInput occurs when a form is submitted but the data is not valid
	ErrInvalidInput = errors.New("one or more form values was invalid")

	// ErrValueName occurs when attempting to access an invalid session value
	ErrValueName = errors.New("invalid session value name")
)

// Error Represents an http service error. Provides methods for the HTTP status code and embeds the
// error interface
type HttpError interface {
	error

	// Status should return the corresponding HTTP status code for the error
	Status() int
}

// StatusError implements the Error interface with an HTTP status code
type StatusError struct {
	Code int
	Err  error
}

// Error implements the error interface and returns the embedded error string
func (e StatusError) Error() string {
	return e.Err.Error()
}

// Status returns the HTTP status code for the error
func (e StatusError) Status() int {
	return e.Code
}

func internalError(e error) StatusError {
	return StatusError{Code: http.StatusInternalServerError, Err: e}
}

func authError(e error) StatusError {
	return StatusError{Code: http.StatusUnauthorized, Err: e}
}
