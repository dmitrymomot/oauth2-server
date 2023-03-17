package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
)

// NewOauth2Server initializes the OAuth2 server.
func NewOauth2Server(
	jwtGen oauth2.AccessGenerate,
	codeGen oauth2.AuthorizeGenerate,
	tokenStorage oauth2.TokenStore,
	clientStorage oauth2.ClientStore,
	authHandler Handler,
) *server.Server {
	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetClientTokenCfg(manage.DefaultClientTokenCfg)
	manager.SetPasswordTokenCfg(manage.DefaultPasswordTokenCfg)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	manager.MapTokenStorage(tokenStorage)
	manager.MapClientStorage(clientStorage)
	manager.MapAccessGenerate(jwtGen)
	manager.MapAuthorizeGenerate(codeGen)

	// Create OAuth2 server
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetAllowedGrantType(
		oauth2.AuthorizationCode,
		oauth2.PasswordCredentials, // TODO: disable in production
		oauth2.ClientCredentials,
		oauth2.Refreshing,
		oauth2.Implicit,
	)
	srv.SetAllowedResponseType(oauth2.Code, oauth2.Token)
	srv.SetInternalErrorHandler(authHandler.InternalErrorHandler)
	srv.SetResponseErrorHandler(authHandler.ResponseErrorHandler)
	srv.SetAuthorizeScopeHandler(authHandler.AuthorizeScopeHandler)
	srv.SetRefreshingScopeHandler(authHandler.RefreshingScopeHandler)
	srv.SetPasswordAuthorizationHandler(authHandler.PasswordAuthorizationHandler)

	return srv
}
