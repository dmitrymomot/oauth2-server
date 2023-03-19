package httpencoder

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/validator"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
)

type (
	logger interface {
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}

	// Error represents an error response
	ErrorResponse struct {
		Code      int         `json:"code"`
		Error     string      `json:"error"`
		Message   string      `json:"message,omitempty"`
		Details   interface{} `json:"details,omitempty"`
		RequestID string      `json:"request_id,omitempty"`
	}
)

// EncodeError ...
func EncodeError(l logger, codeAndMessageFrom func(err error) (int, interface{})) httptransport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			l.Warnf("encode nil error: %#v", err)
			return
		}

		code, msg := codeAndMessageFrom(err)
		if code >= http.StatusInternalServerError {
			// Log only unexpected errors
			l.Errorf("http transport error: %v", err)
		}

		var resp ErrorResponse
		switch msg.(type) {
		case ErrorResponse:
			resp = msg.(ErrorResponse)
		case *ErrorResponse:
			resp = *msg.(*ErrorResponse)
		case *validator.ValidationError:
			resp = ErrorResponse{
				Code:    http.StatusPreconditionFailed,
				Error:   msg.(*validator.ValidationError).Err.Error(),
				Message: "Validation error. See the details.",
				Details: msg.(*validator.ValidationError).Values,
			}
		case validator.ValidationError:
			resp = ErrorResponse{
				Code:    http.StatusPreconditionFailed,
				Error:   msg.(validator.ValidationError).Err.Error(),
				Message: "Validation error. See the details.",
				Details: msg.(validator.ValidationError).Values,
			}
		default:
			resp = ErrorResponse{
				Code:    code,
				Error:   http.StatusText(code),
				Message: fmt.Sprintf("%v", msg),
				Details: nil,
			}
		}
		resp.RequestID = middleware.GetReqID(ctx)

		w.Header().Set(ContentTypeHeader, ContentType)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// CodeAndMessageFrom helper
func CodeAndMessageFrom(err error) (int, interface{}) {
	if err == nil {
		return http.StatusOK, nil
	}

	if errors.Is(err, validator.ErrValidation) {
		return http.StatusPreconditionFailed, err
	}

	if errors.Is(err, jwt.ErrTokenContextMissing) {
		return http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)
	}

	if errors.Is(err, jwt.ErrTokenExpired) ||
		errors.Is(err, jwt.ErrTokenInvalid) ||
		errors.Is(err, jwt.ErrTokenMalformed) ||
		errors.Is(err, jwt.ErrTokenNotActive) ||
		errors.Is(err, jwt.ErrUnexpectedSigningMethod) {
		return http.StatusUnauthorized, err.Error()
	}

	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, err.Error()
	}

	switch err {
	case jwt.ErrTokenContextMissing:
		return http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)
	case jwt.ErrTokenExpired,
		jwt.ErrTokenInvalid,
		jwt.ErrTokenMalformed,
		jwt.ErrTokenNotActive,
		jwt.ErrUnexpectedSigningMethod:
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
