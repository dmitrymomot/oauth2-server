package middleware

import (
	"net/http"
	"strings"

	"github.com/dmitrymomot/oauth2-server/internal/httpencoder"
	"github.com/dmitrymomot/oauth2-server/lib/client"
	"github.com/go-chi/chi/v5/middleware"
)

// AuthMiddleware is a middleware that checks if the request is authorized.
func AuthMiddleware(verifier TokenVerifier) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := getBearerToken(r)
			if token == "" {
				httpencoder.EncodeResponse(r.Context(), w, httpencoder.ErrorResponse{
					Code:      http.StatusUnauthorized,
					Err:       "unauthorized",
					Message:   "Missed or invalid access token",
					RequestID: middleware.GetReqID(r.Context()),
				})
				return
			}

			info, err := verifier.VerifyToken(token, client.TokenTypeAccessToken)
			if err != nil || info == nil || !info.Active {
				httpencoder.EncodeResponse(r.Context(), w, httpencoder.ErrorResponse{
					Code:      http.StatusUnauthorized,
					Err:       "unauthorized",
					Message:   "Access token is invalid",
					RequestID: middleware.GetReqID(r.Context()),
				})
				return
			}

			ctx := SetTokenInfoToContext(r.Context(), info)
			r = r.WithContext(ctx)

			// Call the next middleware/handler in chain
			next.ServeHTTP(w, r)
		})
	}
}

// get bearer token from request
func getBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return ""
	}

	return authHeaderParts[1]
}
