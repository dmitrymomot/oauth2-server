package oauth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dmitrymomot/oauth2-server/internal/httpencoder"
	"github.com/dmitrymomot/oauth2-server/internal/session"
	"github.com/go-chi/chi/v5"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-oauth2/oauth2/v4"
)

type (
	oauth2Server interface {
		HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error
		HandleTokenRequest(w http.ResponseWriter, r *http.Request) error
	}

	logger interface {
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}

	tokenStoreManager interface {
		RemoveAccessToken(ctx context.Context, access string) error
		RemoveRefreshToken(ctx context.Context, refresh string) error
		LoadAccessToken(ctx context.Context, access string) (oauth2.TokenInfo, error)
		LoadRefreshToken(ctx context.Context, refresh string) (oauth2.TokenInfo, error)
	}
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func MakeHTTPHandler(srv oauth2Server, ts tokenStoreManager, log logger, loginURI string) http.Handler {
	r := chi.NewRouter()
	errEncoder := httpencoder.EncodeError(log, codeAndMessageFrom)

	r.Post("/token", httpTokenHandler(srv, errEncoder))
	r.HandleFunc("/authorize", httpAuthorizeHandler(srv, errEncoder, loginURI))
	r.Post("/revoke", httpRevokeTokenHandler(ts, errEncoder))
	r.Post("/introspect", httpIntrospectTokenHandler(ts, errEncoder))

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if resp := NewError(err); resp != nil {
		return resp.Code, resp
	}
	return httpencoder.CodeAndMessageFrom(err)
}

// httpTokenHandler returns an http.HandlerFunc that makes a set of endpoints
// available on predefined paths.
func httpTokenHandler(s oauth2Server, errEncoder httptransport.ErrorEncoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.HandleTokenRequest(w, r); err != nil {
			errEncoder(r.Context(), err, w)
			return
		}
	}
}

// httpAuthorizeHandler returns an http.HandlerFunc that makes a set of endpoints
// available on predefined paths.
func httpAuthorizeHandler(s oauth2Server, errEncoder httptransport.ErrorEncoder, loginURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			errEncoder(r.Context(), ErrMethodNotAllowed, w)
			return
		}

		if !session.IsLoggedIn(r, w) {
			session.StoreRedirectData(r, w)
			session.StoreReturnURI(r, w, r.URL.String())

			http.Redirect(w, r, loginURI, http.StatusFound)
			return
		} else {
			r.Form = session.GetRedirectData(r, w)
		}

		if err := s.HandleAuthorizeRequest(w, r); err != nil {
			errEncoder(r.Context(), err, w)
			return
		}

		if err := session.Logout(r, w); err != nil {
			log.Printf("failed to logout: %v", err)
		}

		log.Println("user redirected to", r.FormValue("redirect_uri"))
	}
}

// httpRevokeTokenHandler returns an http.HandlerFunc that makes a set of endpoints
// available on predefined paths.
func httpRevokeTokenHandler(ts tokenStoreManager, errEncoder httptransport.ErrorEncoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			errEncoder(r.Context(), err, w)
			return
		}

		token := r.PostForm.Get("token")
		tokenType := r.PostForm.Get("token_type_hint")

		switch tokenType {
		case "access_token":
			if err := ts.RemoveAccessToken(r.Context(), token); err != nil {
				errEncoder(r.Context(), err, w)
				return
			}
		case "refresh_token":
			if err := ts.RemoveRefreshToken(r.Context(), token); err != nil {
				errEncoder(r.Context(), err, w)
				return
			}
		default:
			if err := ts.RemoveAccessToken(r.Context(), token); err != nil {
				if err := ts.RemoveRefreshToken(r.Context(), token); err != nil {
					errEncoder(r.Context(), err, w)
					return
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

type (
	IntrospectResponse struct {
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
)

// httpIntrospectTokenHandler returns an http.HandlerFunc that makes a set of endpoints
// available on predefined paths.
func httpIntrospectTokenHandler(ts tokenStoreManager, errEncoder httptransport.ErrorEncoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			errEncoder(r.Context(), err, w)
			return
		}

		token := r.PostForm.Get("token")
		tokenType := r.PostForm.Get("token_type_hint")

		var (
			ti     oauth2.TokenInfo
			err    error
			active bool
			expAt  int64
			iat    int64
		)

		switch tokenType {
		case "access_token":
			ti, err = ts.LoadAccessToken(r.Context(), token)
			if err == nil && ti != nil {
				active = ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).After(time.Now())
				expAt = ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Unix()
				iat = ti.GetAccessCreateAt().Unix()
			}
		case "refresh_token":
			ti, err = ts.LoadRefreshToken(r.Context(), token)
			if err == nil && ti != nil {
				active = ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).After(time.Now())
				expAt = ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Unix()
				iat = ti.GetRefreshCreateAt().Unix()
			}
		default:
			ti, err = ts.LoadAccessToken(r.Context(), token)
			if err == nil && ti != nil {
				active = ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).After(time.Now())
				expAt = ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Unix()
				iat = ti.GetAccessCreateAt().Unix()
				tokenType = "access_token"
			} else {
				ti, err = ts.LoadRefreshToken(r.Context(), token)
				if err == nil && ti != nil {
					active = ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).After(time.Now())
					expAt = ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Unix()
					iat = ti.GetRefreshCreateAt().Unix()
					tokenType = "refresh_token"
				}
			}
		}

		if err != nil {
			errEncoder(r.Context(), err, w)
			return
		}

		if err := json.NewEncoder(w).Encode(IntrospectResponse{
			Active:    active,
			Scope:     ti.GetScope(),
			ClientID:  ti.GetClientID(),
			UserID:    ti.GetUserID(),
			TokenType: tokenType,
			ExpiresAt: expAt,
			IssuedAt:  iat,
			NotBefore: iat,
			Subject:   ti.GetUserID(),
			Audience:  ti.GetClientID(),
		}); err != nil {
			errEncoder(r.Context(), err, w)
			return
		}
	}
}
