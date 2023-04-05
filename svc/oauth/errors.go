package oauth

import (
	"errors"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/httpencoder"
	"github.com/dmitrymomot/oauth2-server/internal/utils"
	oauthErrors "github.com/go-oauth2/oauth2/v4/errors"
)

// Predefined errors
var (
	ErrInvalidRequest     = errors.New("invalid_request")
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrMethodNotAllowed   = errors.New("method_not_allowed")
	ErrInvalidAccessToken = errors.New("invalid_access_token")
	ErrUnauthorized       = errors.New("unauthorized")
)

// Error codes map
var ErrorCodes = map[error]int{
	ErrInvalidRequest:     http.StatusBadRequest,
	ErrInvalidCredentials: http.StatusUnauthorized,
	ErrMethodNotAllowed:   http.StatusMethodNotAllowed,
	ErrInvalidAccessToken: http.StatusUnauthorized,
	ErrUnauthorized:       http.StatusUnauthorized,

	oauthErrors.ErrInvalidRedirectURI:   http.StatusBadRequest,
	oauthErrors.ErrInvalidAuthorizeCode: http.StatusBadRequest,
	oauthErrors.ErrInvalidAccessToken:   http.StatusBadRequest,
	oauthErrors.ErrInvalidRefreshToken:  http.StatusBadRequest,
	oauthErrors.ErrExpiredAccessToken:   http.StatusBadRequest,
	oauthErrors.ErrExpiredRefreshToken:  http.StatusBadRequest,
	oauthErrors.ErrMissingCodeVerifier:  http.StatusBadRequest,
	oauthErrors.ErrMissingCodeChallenge: http.StatusBadRequest,
	oauthErrors.ErrInvalidCodeChallenge: http.StatusBadRequest,
}

// Error messages
var ErrorMessages = map[error]string{
	ErrInvalidRequest:     "Invalid request",
	ErrInvalidCredentials: "Invalid credentials",
	ErrMethodNotAllowed:   "Method not allowed",
	ErrInvalidAccessToken: "Invalid access token",
	ErrUnauthorized:       "Unauthorized",

	oauthErrors.ErrInvalidRedirectURI:   "Invalid redirect uri",
	oauthErrors.ErrInvalidAuthorizeCode: "Invalid authorize code",
	oauthErrors.ErrInvalidAccessToken:   "Invalid access token",
	oauthErrors.ErrInvalidRefreshToken:  "Invalid refresh token",
	oauthErrors.ErrExpiredAccessToken:   "Expired access token",
	oauthErrors.ErrExpiredRefreshToken:  "Expired refresh token",
	oauthErrors.ErrMissingCodeVerifier:  "Missing code verifier",
	oauthErrors.ErrMissingCodeChallenge: "Missing code challenge",
	oauthErrors.ErrInvalidCodeChallenge: "Invalid code challenge",
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
		Err:     utils.ToSnakeCase(stdErr.Error()),
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

// CodeAndMessageFrom returns http error code by error type.
// Returns (0, nil) if error is not found.
// This function can be used to get error code and message from external packages.
func CodeAndMessageFrom(err error) (int, interface{}) {
	var errCode int
	{
		if code, ok := ErrorCodes[err]; ok {
			errCode = code
		}
		if code, ok := oauthErrors.StatusCodes[err]; ok {
			errCode = code
		}
	}

	var errMsg interface{}
	{
		if msg, ok := ErrorMessages[err]; ok {
			errMsg = msg
		}
		if msg, ok := oauthErrors.Descriptions[err]; ok {
			errMsg = msg
		}
	}

	return errCode, errMsg
}
