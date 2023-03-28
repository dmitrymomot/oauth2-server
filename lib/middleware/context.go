package middleware

import (
	"context"

	"github.com/dmitrymomot/oauth2-server/lib/client"
)

// SetTokenInfoToContext sets token info to context.
func SetTokenInfoToContext(ctx context.Context, info *client.TokenInfo) context.Context {
	return context.WithValue(ctx, TokenInfoKey, info)
}

// GetTokenInfoFromContext gets token info from context.
func GetTokenInfoFromContext(ctx context.Context) (*client.TokenInfo, bool) {
	info, ok := ctx.Value(TokenInfoKey).(*client.TokenInfo)
	return info, ok
}

// GetClientIDFromContext gets client id from context.
func GetClientIDFromContext(ctx context.Context) (string, bool) {
	info, ok := GetTokenInfoFromContext(ctx)
	if !ok {
		return "", false
	}
	return info.ClientID, info.ClientID != ""
}

// GetUserIDFromContext gets user id from context.
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	info, ok := GetTokenInfoFromContext(ctx)
	if !ok {
		return "", false
	}
	return info.UserID, info.UserID != ""
}
