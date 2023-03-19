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
		log logger
	}

	logger interface {
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}
)

// NewHandlerLogger returns a new handlerLogger.
func NewHandlerLogger(h Handler, log logger) Handler {
	return &handlerLogger{
		Handler: h,
		log:     log,
	}
}

// AuthorizeScopeHandler logs the request and response.
func (h *handlerLogger) AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	h.log.Debugf("AuthorizeScopeHandler: request=%+v", r)
	scope, err = h.Handler.AuthorizeScopeHandler(w, r)
	h.log.Debugf("AuthorizeScopeHandler: response=%s, err=%v", scope, err)
	return scope, err
}

// RefreshingScopeHandler check the scope of the refreshing token
// and logs the request and response.
func (h *handlerLogger) RefreshingScopeHandler(tgr *oauth2.TokenGenerateRequest, oldScope string) (allowed bool, err error) {
	h.log.Debugf("RefreshingScopeHandler: request=%+v", tgr)
	allowed, err = h.Handler.RefreshingScopeHandler(tgr, oldScope)
	if err != nil {
		h.log.Errorf("RefreshingScopeHandler: %v", err)
	}
	h.log.Debugf("RefreshingScopeHandler: response=%t, err=%v", allowed, err)
	return allowed, err
}

// PasswordAuthorizationHandler get user id from username and password
// and logs the request and response.
func (h *handlerLogger) PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	h.log.Debugf("PasswordAuthorizationHandler: clientID=%s, username=%s, password=%s", clientID, username, password)
	userID, err = h.Handler.PasswordAuthorizationHandler(ctx, clientID, username, password)
	if err != nil {
		h.log.Errorf("PasswordAuthorizationHandler: %v", err)
	}
	h.log.Debugf("PasswordAuthorizationHandler: response=%s, err=%v", userID, err)
	return userID, err
}
