package middleware

import (
	"time"

	"github.com/dmitrymomot/oauth2-server/lib/client"
	"github.com/golang-jwt/jwt/v5"
)

// JwtVerifier is a token verifier that verifies JWT tokens.
type JwtVerifier struct {
	signingKey string
}

// NewJwtVerifier creates a new JwtVerifier.
func NewJwtVerifier(signingKey string) *JwtVerifier {
	return &JwtVerifier{
		signingKey: signingKey,
	}
}

// VerifyToken verifies a token.
func (j *JwtVerifier) VerifyToken(tokenString string, tokenType client.TokenType) (*client.TokenInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return castMapClaimsToTokenInfo(claims), nil
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
