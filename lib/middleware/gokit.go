package middleware

import (
	"context"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/httpencoder"
	"github.com/dmitrymomot/oauth2-server/lib/client"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-oauth2/oauth2/v4/errors"
)

// GokitAuthMiddleware is a middleware for gokit
func GokitAuthMiddleware(verifyFn VerifyTokenFunc) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token, ok := ctx.Value(jwt.JWTContextKey).(string)
			if token == "" || !ok {
				return nil, httpencoder.NewError(http.StatusUnauthorized, errors.ErrInvalidAccessToken, "", nil)
			}

			info, err := verifyFn(token, client.TokenTypeAccessToken)
			if err != nil || info == nil || !info.Active {
				return nil, httpencoder.NewError(http.StatusUnauthorized, errors.ErrInvalidAccessToken, "", nil)
			}

			return next(SetTokenInfoToContext(ctx, info), request)
		}
	}
}
