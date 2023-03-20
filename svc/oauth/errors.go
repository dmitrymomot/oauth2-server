package oauth

import (
	"errors"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/httpencoder"
	oauthErrors "github.com/go-oauth2/oauth2/v4/errors"
)

// Predefined errors
var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrMethodNotAllowed   = errors.New("method_not_allowed")
	ErrInvalidAccessToken = errors.New("invalid_access_token")
	ErrUnauthorized       = errors.New("unauthorized")
)

// Error codes map
var ErrorCodes = map[error]int{
	ErrInvalidCredentials: http.StatusUnauthorized,
	ErrMethodNotAllowed:   http.StatusMethodNotAllowed,
	ErrInvalidAccessToken: http.StatusUnauthorized,
	ErrUnauthorized:       http.StatusUnauthorized,
}

// Error messages
var ErrorMessages = map[error]string{
	ErrInvalidCredentials: "Invalid credentials",
	ErrMethodNotAllowed:   "Method not allowed",
	ErrInvalidAccessToken: "Invalid access token",
	ErrUnauthorized:       "Unauthorized",
}

// NewError creates a new error
func NewError(err error) *httpencoder.ErrorResponse {
	stdErr := findError(err)
	code := findErrorCode(stdErr)
	msg := findErrMessage(stdErr)
	if stdErr == nil {
		stdErr = err
	}

	return &httpencoder.ErrorResponse{
		Code:    code,
		Error:   stdErr.Error(),
		Message: msg,
	}
}

func findError(err error) error {
	for stdErr := range ErrorCodes {
		if errors.Is(err, stdErr) {
			return stdErr
		}
	}
	for stdErr := range oauthErrors.StatusCodes {
		if errors.Is(err, stdErr) {
			return stdErr
		}
	}
	return nil
}

func findErrorCode(err error) int {
	if code, ok := ErrorCodes[err]; ok {
		return code
	}
	if code, ok := oauthErrors.StatusCodes[err]; ok {
		return code
	}
	return http.StatusInternalServerError
}

func findErrMessage(err error) string {
	if msg, ok := ErrorMessages[err]; ok {
		return msg
	}
	if msg, ok := oauthErrors.Descriptions[err]; ok {
		return msg
	}
	return http.StatusText(findErrorCode(err))
}
