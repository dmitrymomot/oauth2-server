package client

import (
	"context"

	"github.com/dmitrymomot/oauth2-server/lib/middleware"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints collects all of the endpoints that compose a client service. It's
	// meant to be used as a helper struct, to collect all of the endpoints into a
	// single parameter.
	Endpoints struct {
		Create      endpoint.Endpoint
		GetByID     endpoint.Endpoint
		GetByUserID endpoint.Endpoint
		Delete      endpoint.Endpoint
	}

	ClientResponse struct {
		Client  *Client   `json:"client,omitempty"`
		Clients []*Client `json:"clients,omitempty"`
	}
)

// MakeEndpoints returns an Endpoints struct where each endpoint invokes the
// corresponding method on the provided service. Primarily useful in a server.
func MakeEndpoints(s Service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		Create:      MakeCreateEndpoint(s),
		GetByID:     MakeGetByIDEndpoint(s),
		Delete:      MakeDeleteEndpoint(s),
		GetByUserID: MakeGetByUserIDEndpoint(s),
	}

	for _, mdw := range m {
		e.Create = mdw(e.Create)
		e.GetByID = mdw(e.GetByID)
		e.Delete = mdw(e.Delete)
		e.GetByUserID = mdw(e.GetByUserID)
	}

	return e
}

// CreateRequest is a request for the Create method.
type CreateRequest struct {
	Domain string `json:"domain" validate:"required|fullUrl" filter:"trim|lower|escapeJs|escapeHtml" label:"Domain"`
	Public bool   `json:"is_public" validate:"bool" label:"Is Public"`
}

// MakeCreateEndpoint returns an endpoint via the passed service.
func MakeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tokenInfo, ok := middleware.GetTokenInfoFromContext(ctx)
		if !ok || tokenInfo == nil || tokenInfo.UserID == "" {
			return nil, ErrForbidden
		}

		req, ok := request.(CreateRequest)
		if !ok {
			return nil, ErrInvalidRequest
		}

		client, err := s.Create(ctx, tokenInfo.UserID, req.Domain, req.Public)
		if err != nil {
			return nil, err
		}

		return ClientResponse{Client: client}, nil
	}
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
		client, err := s.GetByID(ctx, req)
		if err != nil {
			return nil, err
		}

		if tokenInfo.ClientID != client.ID && tokenInfo.UserID != client.UserID {
			return nil, ErrForbidden
		}

		return ClientResponse{Client: client}, nil
	}
}

// MakeGetByUserIDEndpoint returns an endpoint via the passed service.
func MakeGetByUserIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tokenInfo, ok := middleware.GetTokenInfoFromContext(ctx)
		if !ok || tokenInfo == nil || tokenInfo.UserID == "" {
			return nil, ErrForbidden
		}

		clients, err := s.GetByUserID(ctx, tokenInfo.UserID)
		if err != nil {
			return nil, err
		}

		return ClientResponse{Clients: clients}, nil
	}
}

// MakeDeleteEndpoint returns an endpoint via the passed service.
func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tokenInfo, ok := middleware.GetTokenInfoFromContext(ctx)
		if !ok || tokenInfo == nil || tokenInfo.UserID == "" {
			return nil, ErrForbidden
		}

		req, ok := request.(string)
		if !ok {
			return nil, ErrInvalidRequest
		}

		client, err := s.GetByID(ctx, req)
		if err != nil {
			return nil, err
		}

		if tokenInfo.ClientID == client.ID && tokenInfo.UserID != client.UserID {
			return nil, ErrForbidden
		}

		if err := s.Delete(ctx, client.ID); err != nil {
			return nil, err
		}

		return true, nil
	}
}
