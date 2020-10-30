package goodreads

import (
	"errors"
	"fmt"
)

// ErrUnexpectedResponse is returned when an API call received a response with a status code it did not expect.
type ErrUnexpectedResponse struct {
	Code int
}

func (err ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("unexpected status code: %d", err.Code)
}

var _ error = ErrUnexpectedResponse{}

// ErrAPIKeyNotSet is returned when an API call was attempted without an API key set when required.
type ErrAPIKeyNotSet struct{}

func (err ErrAPIKeyNotSet) Error() string {
	return "goodreads API key not set"
}

var _ error = ErrAPIKeyNotSet{}

// ErrNotFound is returned when an API call could not find a requested resource.
type ErrNotFound struct{}

func (err ErrNotFound) Error() string {
	return "not found"
}

var _ error = ErrNotFound{}

// IsNotFound returns true if err is an ErrNotFound.
func IsNotFound(err error) bool {
	var e ErrNotFound

	return errors.As(err, &e)
}
