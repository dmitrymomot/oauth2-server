package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/httpencoder"
	"github.com/dmitrymomot/oauth2-server/internal/kitlog"
	"github.com/go-chi/chi/v5"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

type (
	logger interface {
		Println(args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}
)

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(kitlog.NewLogger(log))),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Post("/", httptransport.NewServer(
		e.Create,
		decodeCreateRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/", httptransport.NewServer(
		e.GetByUserID,
		decodeGetByUserIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{id}", httptransport.NewServer(
		e.GetByID,
		decodeGetByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{id}", httptransport.NewServer(
		e.Delete,
		decodeDeleteRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if resp := NewError(err); resp != nil {
		return resp.Code, resp
	}

	return httpencoder.CodeAndMessageFrom(err)
}

// decodeCreateRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}

// decodeGetByUserIDRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetByUserIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

// decodeGetByIDRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrInvalidParameter
	}

	return id, nil
}

// decodeDeleteRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrInvalidParameter
	}

	return id, nil
}
