package client

import (
	"net/http"
	"strings"
)

// SetHTTPClient sets the HTTP client to use for requests.
func SetHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// SetUserAPIEndpoint sets the endpoint for the User API.
func SetUserAPIEndpoint(endpoint string) ClientOption {
	return func(c *client) {
		if endpoint != "" && strings.HasPrefix(endpoint, "http") {
			c.userAPIEndpoint = endpoint
		}
	}
}

// SetClientAPIEndpoint sets the endpoint for the Client API.
func SetClientAPIEndpoint(endpoint string) ClientOption {
	return func(c *client) {
		if endpoint != "" && strings.HasPrefix(endpoint, "http") {
			c.clientAPIEndpoint = endpoint
		}
	}
}

// SetIntrospectEndpoint sets the endpoint for the Introspect API.
func SetIntrospectEndpoint(endpoint string) ClientOption {
	return func(c *client) {
		if endpoint != "" && strings.HasPrefix(endpoint, "http") {
			c.introspectEndpoint = endpoint
		}
	}
}
