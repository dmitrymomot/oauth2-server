package middleware

import (
	"context"

	"github.com/dmitrymomot/oauth2-server/lib/client"
	"github.com/dmitrymomot/oauth2-server/svc/oauth"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

// GokitAuthMiddleware is a middleware for gokit
func GokitAuthMiddleware(verifyFn VerifyTokenFunc) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token, ok := ctx.Value(jwt.JWTContextKey).(string)
			if token == "" || !ok {
				return nil, oauth.ErrInvalidAccessToken
			}

			info, err := verifyFn(token, client.TokenTypeAccessToken)
			if err != nil || info == nil || !info.Active {
				return nil, oauth.ErrInvalidAccessToken
			}

			return next(SetTokenInfoToContext(ctx, info), request)
		}
	}
}
