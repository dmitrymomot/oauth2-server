package middleware

import (
	"context"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/httpencoder"
	"github.com/dmitrymomot/oauth2-server/lib/client"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

// GokitAuthMiddleware is a middleware for gokit
func GokitAuthMiddleware(verifyFn VerifyTokenFunc) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token, ok := ctx.Value(jwt.JWTContextKey).(string)
			if token == "" || !ok {
				return nil, httpencoder.ErrorResponse{
					Code:      http.StatusUnauthorized,
					Err:       "unauthorized",
					Message:   "Missed or invalid access token",
					RequestID: middleware.GetReqID(ctx),
				}
			}

			info, err := verifyFn(token, client.TokenTypeAccessToken)
			if err != nil || info == nil || !info.Active {
				return nil, httpencoder.ErrorResponse{
					Code:      http.StatusUnauthorized,
					Err:       "unauthorized",
					Message:   "Access token is invalid",
					RequestID: middleware.GetReqID(ctx),
				}
			}

			return next(SetTokenInfoToContext(ctx, info), request)
		}
	}
}
