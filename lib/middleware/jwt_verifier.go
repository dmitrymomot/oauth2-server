package middleware

import (
	"time"

	"github.com/dmitrymomot/oauth2-server/lib/client"
	"github.com/golang-jwt/jwt/v5"
)

// VerifyJWT verifies a token and returns the token info.
// This function is compatible with the VerifyTokenFunc interface.
func VerifyJWT(signingKey string) func(string, client.TokenType) (*client.TokenInfo, error) {
	return func(tokenString string, tokenType client.TokenType) (*client.TokenInfo, error) {
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})
		if err != nil {
			return nil, err
		}

		if !token.Valid {
			return nil, jwt.ErrTokenMalformed
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			return nil, jwt.ErrTokenRequiredClaimMissing
		}

		return castMapClaimsToTokenInfo(claims), nil
	}
}

// Cast *jwt.MapClaims to *client.TokenInfo
func castMapClaimsToTokenInfo(claims *jwt.MapClaims) *client.TokenInfo {
	result := &client.TokenInfo{
		Active: false,
	}

	if audArr, err := claims.GetAudience(); err == nil && len(audArr) > 0 {
		result.Audience = audArr[0]
		result.ClientID = audArr[0]
	}
	if sub, err := claims.GetSubject(); err == nil {
		result.Subject = sub
		result.UserID = sub
	}
	if ext, err := claims.GetExpirationTime(); err == nil && !ext.IsZero() {
		result.ExpiresAt = ext.Unix()
		if time.Now().Before(ext.Time) {
			result.Active = true
		}
	}

	return result
}
