package client

import "fmt"

// TokenType is a type of token.
type TokenType string

// Token types.
const (
	TokenTypeAccessToken  TokenType = "access_token"
	TokenTypeRefreshToken TokenType = "refresh_token"
)

// TokenInfo is a struct that contains information about a token.
type TokenInfo struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
	Audience  string `json:"aud,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	TokenID   string `json:"jti,omitempty"`
}

// ErrorResponse is a struct that contains an error message.
type ErrorResponse struct {
	Code      int         `json:"code"`
	Err       string      `json:"error"`
	Message   string      `json:"message,omitempty"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// Error implements the error interface.
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Code, e.Err, e.Message)
}
