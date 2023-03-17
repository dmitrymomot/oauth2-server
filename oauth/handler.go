package oauth

import (
	"context"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"golang.org/x/crypto/bcrypt"
)

type (
	Handler interface {
		AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error)
		RefreshingScopeHandler(tgr *oauth2.TokenGenerateRequest, oldScope string) (allowed bool, err error)
		PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error)
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

	handlerRepository interface {
		GetUserByEmail(ctx context.Context, email string) (repository.User, error)
		GetClientByID(ctx context.Context, id string) (repository.Client, error)
	}
)

// NewHandler creates a new oauth2 handler instance.
func NewHandler(repo handlerRepository) Handler {
	return &handler{repo: repo}
}

// AuthorizeScopeHandler check the scope of the authorization request
func (h *handler) AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	scopes := getScopeFromRequest(r)
	if scopes == "" {
		return "", nil
	}

	gt := getScopeFromRequest(r)
	if gt == "" {
		return "", errors.ErrInvalidGrant
	}

	var allowedScope string
	switch oauth2.GrantType(gt) {
	case oauth2.PasswordCredentials:
		allowedScope = h.passwordScope
	case oauth2.ClientCredentials:
		allowedScope = h.clientScope
	case oauth2.AuthorizationCode:
		allowedScope = h.codeScope
	}

	if allowedScope == "" {
		return "", nil
	}

	if !MatchScopesStrict(scopes, allowedScope) {
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
