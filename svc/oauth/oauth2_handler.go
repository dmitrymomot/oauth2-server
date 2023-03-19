package oauth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"golang.org/x/crypto/bcrypt"
)

type (
	Handler interface {
		ClientAuthorizedHandler(clientID string, grant oauth2.GrantType) (allowed bool, err error)
		ClientScopeHandler(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error)
		AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error)
		RefreshingScopeHandler(tgr *oauth2.TokenGenerateRequest, oldScope string) (allowed bool, err error)
		UserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error)
		PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error)
		ExtensionFieldsHandler(ti oauth2.TokenInfo) (fieldsValue map[string]interface{})
		ResponseErrorHandler(re *errors.Response)
		InternalErrorHandler(err error) (re *errors.Response)
	}

	handler struct {
		repo handlerRepository

		// default scope for auth grant types
		passwordScope string // password grant type
		clientScope   string // client_credentials grant type
		codeScope     string // authorization_code grant type
	}

	handlerOption func(h *handler)

	handlerRepository interface {
		GetUserByEmail(ctx context.Context, email string) (repository.User, error)
		GetClientByID(ctx context.Context, id string) (repository.Client, error)
		GetTokenByAccess(ctx context.Context, access string) (repository.Token, error)
	}
)

// WithPasswordScope sets the default scope for password grant type
func WithPasswordScope(scope string) handlerOption {
	return func(h *handler) {
		h.passwordScope = scope
	}
}

// WithClientScope sets the default scope for client_credentials grant type
func WithClientScope(scope string) handlerOption {
	return func(h *handler) {
		h.clientScope = scope
	}
}

// WithCodeScope sets the default scope for authorization_code grant type
func WithCodeScope(scope string) handlerOption {
	return func(h *handler) {
		h.codeScope = scope
	}
}

// NewHandler creates a new oauth2 handler instance.
func NewHandler(repo handlerRepository, opts ...handlerOption) Handler {
	h := &handler{repo: repo}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// ClientAuthorizedHandler check the client is allowed to use the grant type
func (h *handler) ClientAuthorizedHandler(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := h.repo.GetClientByID(ctx, clientID)
	if err != nil {
		return false, err
	}

	if len(client.AllowedGrants) > 0 {
		for _, g := range client.AllowedGrants {
			if g == string(grant) {
				return true, nil
			}
		}
	}

	return false, nil
}

// ClientScopeHandler check the client is allowed to use the scope
func (h *handler) ClientScopeHandler(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := h.repo.GetClientByID(ctx, tgr.ClientID)
	if err != nil {
		return false, err
	}

	if MatchScopesStrict(tgr.Scope, client.Scope) {
		return true, nil
	}

	return false, errors.ErrInvalidScope
}

// AuthorizeScopeHandler check the scope of the authorization request
func (h *handler) AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	scopes := getScopeFromRequest(r)
	if scopes == "" {
		return "", nil
	}

	if !MatchScopesStrict(scopes, h.codeScope) {
		return "", errors.ErrInvalidScope
	}

	return scopes, nil
}

// RefreshingScopeHandler check the scope of the refreshing token
func (h *handler) RefreshingScopeHandler(tgr *oauth2.TokenGenerateRequest, oldScope string) (allowed bool, err error) {
	if !MatchScopesStrict(tgr.Scope, oldScope) {
		return false, errors.ErrInvalidScope
	}

	return true, nil
}

// UserAuthorizationHandler get user id from authorization request
func (h *handler) UserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	token, err := h.ValidationBearerToken(r)
	if err != nil {
		return "", err
	}

	return token.GetUserID(), nil
}

// PasswordAuthorizationHandler get user id from username and password
func (h *handler) PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	if _, err := h.repo.GetClientByID(ctx, clientID); err != nil {
		return "", errors.ErrInvalidClient
	}

	user, err := h.repo.GetUserByEmail(ctx, username)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}

	return user.ID.String(), nil
}

// ExtensionFieldsHandler set extension fields in the token response
func (h *handler) ExtensionFieldsHandler(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {
	result := map[string]interface{}{
		"token_type": "Bearer",
		"expires_in": ti.GetAccessExpiresIn(),
	}

	if uid := ti.GetUserID(); uid != "" {
		result["user_id"] = uid
	}

	return result
}

// ResponseErrorHandler response error handing
func (h *handler) ResponseErrorHandler(re *errors.Response) {
	// do nothing
}

// InternalErrorHandler internal error handing
func (h *handler) InternalErrorHandler(err error) (re *errors.Response) {
	switch err {
	case ErrInvalidCredentials:
		return &errors.Response{
			Error:       err,
			ErrorCode:   http.StatusUnauthorized,
			Description: "Invalid credentials",
			StatusCode:  http.StatusUnauthorized,
		}
	}

	return &errors.Response{
		Error:       err,
		ErrorCode:   http.StatusInternalServerError,
		Description: "Something went wrong, please try again later or contact support",
		StatusCode:  http.StatusInternalServerError,
	}
}

// BearerAuth parse bearer token
func (h *handler) BearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access_token")
	}

	return token, token != ""
}

// ValidationBearerToken validation the bearer tokens
// https://tools.ietf.org/html/rfc6750
func (h *handler) ValidationBearerToken(r *http.Request) (oauth2.TokenInfo, error) {
	ctx := r.Context()

	accessToken, ok := h.BearerAuth(r)
	if !ok {
		return nil, ErrInvalidAccessToken
	}

	token, err := h.repo.GetTokenByAccess(ctx, accessToken)
	if err != nil {
		return nil, ErrInvalidAccessToken
	}

	return NewToken(token), nil
}
