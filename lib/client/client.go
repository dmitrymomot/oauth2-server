package client

import (
	"context"
	"net/http"
)

type (
	Client interface {
		Introspect(ctx context.Context, token string, tokenType TokenType) (*TokenInfo, error)
	}

	// ClientOption is a function that configures a Client.
	ClientOption func(*client)

	client struct {
		httpClient *http.Client

		introspectEndpoint string
		userAPIEndpoint    string
		clientAPIEndpoint  string
	}
)

// NewClient returns a new Client.
func NewClient(opts ...ClientOption) Client {
	c := &client{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Introspect returns the token introspection response
func (c *client) Introspect(ctx context.Context, token string, tokenType TokenType) (*TokenInfo, error) {
	return Introspect(c.introspectEndpoint)(ctx, token, tokenType)
}
