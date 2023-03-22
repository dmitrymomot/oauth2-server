package user

import (
	"context"

	"github.com/dmitrymomot/oauth2-server/internal/validator"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints collects all of the endpoints that compose a user service. It's
	// meant to be used as a helper struct, to collect all of the endpoints into a
	// single parameter.
	Endpoints struct {
		GetByID        endpoint.Endpoint
		UpdateEmail    endpoint.Endpoint
		UpdatePassword endpoint.Endpoint
		Delete         endpoint.Endpoint
	}

	UserResponse struct {
		User *User `json:"user"`
	}
)

// MakeEndpoints returns an Endpoints struct where each endpoint invokes the
// corresponding method on the provided service. Primarily useful in a server.
func MakeEndpoints(s Service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetByID:        MakeGetByIDEndpoint(s),
		UpdateEmail:    MakeUpdateEmailEndpoint(s),
		UpdatePassword: MakeUpdatePasswordEndpoint(s),
		Delete:         MakeDeleteEndpoint(s),
	}

	for _, mdw := range m {
		e.GetByID = mdw(e.GetByID)
		e.UpdateEmail = mdw(e.UpdateEmail)
		e.UpdatePassword = mdw(e.UpdatePassword)
		e.Delete = mdw(e.Delete)
	}

	return e
}

// MakeGetByIDEndpoint returns an endpoint via the passed service.
func MakeGetByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(string)
		if !ok {
			return nil, ErrInvalidRequest
		}
		u, err := s.GetByID(ctx, req)
		if err != nil {
			return nil, err
		}
		return UserResponse{User: u}, nil
	}
}

// UpdateEmailRequest is the request type for the UpdateEmail endpoint.
type UpdateEmailRequest struct {
	ID    string `json:"-"`
	Email string `json:"email" form:"email" validate:"required|email|realEmail" filter:"trim|lower|escapeJs|escapeHtml|sanitizeEmail" label:"Email address"`
}

// MakeUpdateEmailEndpoint returns an endpoint via the passed service.
func MakeUpdateEmailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateEmailRequest)
		if !ok {
			return nil, ErrInvalidRequest
		}
		u, err := s.UpdateEmail(ctx, req.ID, req.Email)
		if err != nil {
			return nil, err
		}
		return UserResponse{User: u}, nil
	}
}

// UpdatePasswordRequest is the request type for the UpdatePassword endpoint.
type UpdatePasswordRequest struct {
	ID          string `json:"-"`
	Password    string `json:"password" validate:"required" filter:"trim" label:"Current password"`
	NewPassword string `json:"new_password" validate:"required|minLen:8|maxLen:50" filter:"trim" label:"New password"`
}

// MakeUpdatePasswordEndpoint returns an endpoint via the passed service.
func MakeUpdatePasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdatePasswordRequest)
		if !ok {
			return nil, ErrInvalidRequest
		}
		if v := validator.ValidateStruct(req); len(v) > 0 {
			return nil, validator.NewValidationError(v)
		}

		if err := s.UpdatePassword(ctx, req.ID, req.Password, req.NewPassword); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeDeleteEndpoint returns an endpoint via the passed service.
func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(string)
		if !ok {
			return nil, ErrInvalidRequest
		}
		if err := s.Delete(ctx, req); err != nil {
			return nil, err
		}
		return true, nil
	}
}
