package validator

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// ErrValidation is returned when validation fails.
var ErrValidation = errors.New("validation_failed")

// ValidationError is a custom error type that contains a map of field names to
// error messages.
type ValidationError struct {
	Err    error      // The underlying error.
	Values url.Values // A map of field names to error messages.
}

// NewValidationError returns a new ValidationError.
func NewValidationError(values url.Values) *ValidationError {
	return &ValidationError{
		Err:    ErrValidation,
		Values: values,
	}
}

// Error returns the error message.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Err, e.Values)
}

// Unwrap returns the underlying error.
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// AddValue adds a new error message to the map.
func (e *ValidationError) AddValue(key, message string) *ValidationError {
	e.Values.Add(key, message)
	return e
}

// AddValues adds a new error message to the map.
func (e *ValidationError) AddValues(values url.Values) *ValidationError {
	for key, value := range values {
		e.Values[key] = value
	}
	return e
}
