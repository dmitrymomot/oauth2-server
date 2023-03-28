package middleware

import "github.com/dmitrymomot/oauth2-server/lib/client"

// ContextKey is a key for context.
type ContextKey struct{}

// TokenInfoKey is a key for token info in context.
var TokenInfoKey = ContextKey{}

// TokenVerifier is a function interface that can be used to verify tokens.
type VerifyTokenFunc func(token string, tokenType client.TokenType) (*client.TokenInfo, error)
