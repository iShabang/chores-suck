package errors

// Error Represents an http service error. Provides methods for the HTTP status code and embeds the
// error interface
type Error interface {
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
