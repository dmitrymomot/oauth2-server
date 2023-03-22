package client

import (
	"context"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
)

// ClientCredentialsConfig represents the configuration for the client_credentials grant type
type ClientCredentialsConfig struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	Scopes       []string
}

// NewClientCredentialsClient creates a new OAuth2 client that supports the client_credentials grant type
func NewClientCredentialsClient(ctx context.Context, cfg ClientCredentialsConfig) *http.Client {
	// Set up a configuration object for the OAuth2 client
	oauth2Config := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
		Scopes:       cfg.Scopes,
	}

	return oauth2Config.Client(ctx)
}
