package oauth

import (
	"context"
	"net/http"

	"github.com/go-oauth2/oauth2/v4"
)

type (
	// handlerLogger is a decorator for oauth.handler.
	handlerLogger struct {
		Handler
		log oauthLogger
	}

	oauthLogger interface {
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}
)

// NewHandlerLogger returns a new handlerLogger.
func NewHandlerLogger(h Handler, log oauthLogger) Handler {
	return &handlerLogger{
		Handler: h,
		log:     log,
	}
}

// ClientAuthorizedHandler check the client is authorized or not
// and logs the request and response.
func (h *handlerLogger) ClientAuthorizedHandler(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
	h.log.Debugf("ClientAuthorizedHandler: clientID=%s, grant=%s", clientID, grant)

	allowed, err = h.Handler.ClientAuthorizedHandler(clientID, grant)
	if err != nil {
		h.log.Errorf("ClientAuthorizedHandler: %v", err)
		return false, err
	}

	h.log.Debugf("ClientAuthorizedHandler: allowed=%t", allowed)
	return allowed, nil
}

// ClientScopeHandler check the scope of the client
// and logs the request and response.
func (h *handlerLogger) ClientScopeHandler(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
	h.log.Debugf("ClientScopeHandler: request=%+v", tgr)

	allowed, err = h.Handler.ClientScopeHandler(tgr)
	if err != nil {
		h.log.Errorf("ClientScopeHandler: %v", err)
		return false, err
	}

	h.log.Debugf("ClientScopeHandler: allowed=%t", allowed)
	return allowed, nil
}

// AuthorizeScopeHandler logs the request and response.
func (h *handlerLogger) AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	h.log.Debugf("AuthorizeScopeHandler: request=%+v", r)

	scope, err = h.Handler.AuthorizeScopeHandler(w, r)
	if err != nil {
		h.log.Errorf("AuthorizeScopeHandler: %v", err)
		return "", err
	}

	h.log.Debugf("AuthorizeScopeHandler: scope=%s", scope)
	return scope, nil
}

// RefreshingScopeHandler check the scope of the refreshing token
// and logs the request and response.
func (h *handlerLogger) RefreshingScopeHandler(tgr *oauth2.TokenGenerateRequest, oldScope string) (allowed bool, err error) {
	h.log.Debugf("RefreshingScopeHandler: request=%+v", tgr)

	allowed, err = h.Handler.RefreshingScopeHandler(tgr, oldScope)
	if err != nil {
		h.log.Errorf("RefreshingScopeHandler: %v", err)
		return false, err
	}

	h.log.Debugf("RefreshingScopeHandler: allowed=%t", allowed)
	return allowed, nil
}

// UserAuthorizationHandler logs the request and response.
func (h *handlerLogger) UserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	h.log.Debugf("UserAuthorizationHandler: request=%+v", r)

	userID, err = h.Handler.UserAuthorizationHandler(w, r)
	if err != nil {
		h.log.Errorf("UserAuthorizationHandler: %v", err)
		return "", err
	}

	h.log.Debugf("UserAuthorizationHandler: user_id=%s", userID)
	return userID, nil
}

// PasswordAuthorizationHandler get user id from username and password
// and logs the request and response.
func (h *handlerLogger) PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	h.log.Debugf("PasswordAuthorizationHandler: clientID=%s, username=%s, password=%s", clientID, username, password)

	userID, err = h.Handler.PasswordAuthorizationHandler(ctx, clientID, username, password)
	if err != nil {
		h.log.Errorf("PasswordAuthorizationHandler: %v", err)
		return "", err
	}

	h.log.Debugf("PasswordAuthorizationHandler: user_id=%s", userID)
	return userID, nil
}

// ExtensionFieldsHandler logs the request and response.
func (h *handlerLogger) ExtensionFieldsHandler(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {
	h.log.Debugf("ExtensionFieldsHandler: token_info=%+v", ti)
	fieldsValue = h.Handler.ExtensionFieldsHandler(ti)
	h.log.Debugf("ExtensionFieldsHandler: fields_value=%+v", fieldsValue)
	return fieldsValue
}
