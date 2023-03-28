package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// Introspect returns the token introspection response
func Introspect(endpoint string) func(string, TokenType) (*TokenInfo, error) {
	return func(token string, tokenType TokenType) (*TokenInfo, error) {
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
