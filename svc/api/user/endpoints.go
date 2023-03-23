package user

import (
	"context"

	"github.com/dmitrymomot/oauth2-server/internal/validator"
	"github.com/dmitrymomot/oauth2-server/lib/middleware"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints collects all of the endpoints that compose a user service. It's
	// meant to be used as a helper struct, to collect all of the endpoints into a
	// single parameter.
	Endpoints struct {
		GetByID        endpoint.Endpoint
		GetProfile     endpoint.Endpoint
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
		GetProfile:     MakeGetProfileEndpoint(s),
		UpdateEmail:    MakeUpdateEmailEndpoint(s),
		UpdatePassword: MakeUpdatePasswordEndpoint(s),
		Delete:         MakeDeleteEndpoint(s),
	}

	for _, mdw := range m {
		e.GetByID = mdw(e.GetByID)
		e.GetProfile = mdw(e.GetProfile)
		e.UpdateEmail = mdw(e.UpdateEmail)
		e.UpdatePassword = mdw(e.UpdatePassword)
		e.Delete = mdw(e.Delete)
	}

	return e
}

// MakeGetByIDEndpoint returns an endpoint via the passed service.
func MakeGetByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tokenInfo, ok := middleware.GetTokenInfoFromContext(ctx)
		if !ok || tokenInfo == nil || tokenInfo.ClientID == "" {
			return nil, ErrForbidden
		}

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

// MakeGetProfileEndpoint returns an endpoint via the passed service.
func MakeGetProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tokenInfo, ok := middleware.GetTokenInfoFromContext(ctx)
		if !ok || tokenInfo == nil || tokenInfo.UserID == "" {
			return nil, ErrForbidden
		}

		u, err := s.GetByID(ctx, tokenInfo.UserID)
		if err != nil {
			return nil, err
		}
		return UserResponse{User: u}, nil
	}
}

// UpdateEmailRequest is the request type for the UpdateEmail endpoint.
type UpdateEmailRequest struct {
	Email string `json:"email" form:"email" validate:"required|email|realEmail" filter:"trim|lower|escapeJs|escapeHtml|sanitizeEmail" label:"Email address"`
}

// MakeUpdateEmailEndpoint returns an endpoint via the passed service.
func MakeUpdateEmailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tokenInfo, ok := middleware.GetTokenInfoFromContext(ctx)
		if !ok || tokenInfo == nil || tokenInfo.UserID == "" {
			return nil, ErrForbidden
		}

		req, ok := request.(UpdateEmailRequest)
		if !ok {
			return nil, ErrInvalidRequest
		}
		u, err := s.UpdateEmail(ctx, tokenInfo.UserID, req.Email)
		if err != nil {
			return nil, err
		}
		return UserResponse{User: u}, nil
	}
}

// UpdatePasswordRequest is the request type for the UpdatePassword endpoint.
type UpdatePasswordRequest struct {
	Password    string `json:"password" validate:"required" filter:"trim" label:"Current password"`
	NewPassword string `json:"new_password" validate:"required|minLen:8|maxLen:50" filter:"trim" label:"New password"`
}

// MakeUpdatePasswordEndpoint returns an endpoint via the passed service.
func MakeUpdatePasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tokenInfo, ok := middleware.GetTokenInfoFromContext(ctx)
		if !ok || tokenInfo == nil || tokenInfo.UserID == "" {
			return nil, ErrForbidden
		}

		req, ok := request.(UpdatePasswordRequest)
		if !ok {
			return nil, ErrInvalidRequest
		}
		if v := validator.ValidateStruct(req); len(v) > 0 {
			return nil, validator.NewValidationError(v)
		}

		if err := s.UpdatePassword(ctx, tokenInfo.UserID, req.Password, req.NewPassword); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeDeleteEndpoint returns an endpoint via the passed service.
func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tokenInfo, ok := middleware.GetTokenInfoFromContext(ctx)
		if !ok || tokenInfo == nil || tokenInfo.UserID == "" {
			return nil, ErrForbidden
		}

		if err := s.Delete(ctx, tokenInfo.UserID); err != nil {
			return nil, err
		}
		return true, nil
	}
}
