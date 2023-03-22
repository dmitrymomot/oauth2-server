package middleware

import (
	"context"

	"github.com/dmitrymomot/oauth2-server/lib/client"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-oauth2/oauth2/v4/errors"
)

// GokitAuthMiddleware is a middleware for gokit
func GokitAuthMiddleware(verifier TokenVerifier) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token, ok := ctx.Value(jwt.JWTContextKey).(string)
			if token == "" || !ok {
				return nil, errors.ErrInvalidAccessToken
			}

			info, err := verifier.VerifyToken(token, client.TokenTypeAccessToken)
			if err != nil || info == nil || !info.Active {
				return nil, errors.ErrInvalidAccessToken
			}

			return next(SetTokenInfoToContext(ctx, info), request)
		}
	}
}
