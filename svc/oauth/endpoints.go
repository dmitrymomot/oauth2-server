package oauth

import (
	"context"

	"github.com/dmitrymomot/oauth2-server/internal/validator"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints collects all of the endpoints that compose a auth service. It's
	// meant to be used as a helper struct, to collect all of the endpoints into a
	// single parameter.
	Endpoints struct {
		Authorize   endpoint.Endpoint
		Token       endpoint.Endpoint
		RevokeToken endpoint.Endpoint
	}
)

// Init endpoints for auth service
func InitEndpoints(s oauth2Server) Endpoints {
	return Endpoints{
		Authorize: MakeAuthorizeEndpoint(s),
		// Token:       MakeTokenEndpoint(s),
		// RevokeToken: MakeRevokeTokenEndpoint(s),
	}
}

type (
	AuthorizeRequest struct{}
)

// MakeAuthorizeEndpoint returns an endpoint via the passed service.
func MakeAuthorizeEndpoint(s oauth2Server) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(AuthorizeRequest)
		if !ok {
			return nil, ErrInvalidRequest
		}
		if v := validator.ValidateStruct(req); len(v) > 0 {
			return nil, validator.NewValidationError(v)
		}

		return nil, nil
	}
}
