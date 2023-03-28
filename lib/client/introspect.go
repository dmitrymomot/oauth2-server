package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// Verifier is a function interface for verifying tokens.
type Verifier func(ctx context.Context, token string, tokenType TokenType) (*TokenInfo, error)

// Introspect returns the token introspection response
func Introspect(endpoint string) Verifier {
	return func(ctx context.Context, token string, tokenType TokenType) (*TokenInfo, error) {
		// use the default http client, because token introspection does not require
		// any authentication headers.
		resp, err := http.Post(
			endpoint,
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
}
