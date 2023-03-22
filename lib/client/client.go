package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	// use the default http client, because token introspection does not require
	// any authentication headers.
	resp, err := http.Post(
		c.introspectEndpoint,
		"application/x-www-form-urlencoded",
		strings.NewReader(url.Values{
			"token":           {token},
			"token_type_hint": {string(tokenType)},
		}.Encode()),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var ti TokenInfo
		if err := json.NewDecoder(resp.Body).Decode(&ti); err != nil {
			return nil, err
		}
		return &ti, nil
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return nil, err
	}
	return nil, errResp
}

// VerifyTokenViaIntrospect verifies the token via introspection endpoint.
func (c *client) VerifyTokenViaIntrospect(token string, tokenType TokenType) (*TokenInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.Introspect(ctx, token, tokenType)
}
