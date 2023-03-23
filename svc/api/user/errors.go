package user

import (
	"errors"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/httpencoder"
)

// Predefined errors.
var (
	ErrUserNotFound     = errors.New("user_not_found")
	ErrInvalidPassword  = errors.New("invalid_password")
	ErrInvalidRequest   = errors.New("invalid_request")
	ErrInvalidParameter = errors.New("invalid_parameter")
	ErrForbidden        = errors.New("forbidden")
)

// Error codes map
var ErrorCodes = map[error]int{
	ErrUserNotFound:     http.StatusNotFound,
	ErrInvalidPassword:  http.StatusPreconditionFailed,
	ErrInvalidRequest:   http.StatusBadRequest,
	ErrInvalidParameter: http.StatusBadRequest,
	ErrForbidden:        http.StatusForbidden,
}

// Error messages
var ErrorMessages = map[error]string{
	ErrUserNotFound:     "User not found",
	ErrInvalidPassword:  "Invalid current password",
	ErrInvalidRequest:   "Invalid request",
	ErrInvalidParameter: "Invalid parameter",
	ErrForbidden:        "Forbidden action",
}

// NewError creates a new error
func NewError(err error) *httpencoder.ErrorResponse {
	code, ok := ErrorCodes[err]
	if !ok {
		if stdErr := findError(err); stdErr != nil {
			code, ok = ErrorCodes[stdErr]
		} else {
			return nil
		}
	}

	errStr := err.Error()
	msg, ok := ErrorMessages[err]
	if !ok {
		errStr = http.StatusText(code)
		msg = err.Error()
	}

	return &httpencoder.ErrorResponse{
		Code:    code,
		Err:     errStr,
		Message: msg,
	}
}

func findError(err error) error {
	for stdErr := range ErrorCodes {
		if errors.Is(err, stdErr) {
			return stdErr
		}
	}
	return nil
}
