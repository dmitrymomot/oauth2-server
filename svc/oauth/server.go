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
) (*server.Server, *manage.Manager) {
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
	srv.SetTokenType("Bearer")
	srv.SetAllowGetAccessRequest(true)
	srv.SetAllowedResponseType(oauth2.Code, oauth2.Token)
	srv.SetAllowedGrantType(
		oauth2.AuthorizationCode,
		oauth2.PasswordCredentials,
		oauth2.ClientCredentials,
		oauth2.Refreshing,
		oauth2.Implicit,
	)

	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetClientAuthorizedHandler(authHandler.ClientAuthorizedHandler)
	srv.SetClientScopeHandler(authHandler.ClientScopeHandler)

	srv.SetUserAuthorizationHandler(authHandler.UserAuthorizationHandler)
	srv.SetPasswordAuthorizationHandler(authHandler.PasswordAuthorizationHandler)
	srv.SetRefreshingScopeHandler(authHandler.RefreshingScopeHandler)

	srv.SetResponseErrorHandler(authHandler.ResponseErrorHandler)
	srv.SetInternalErrorHandler(authHandler.InternalErrorHandler)

	srv.SetExtensionFieldsHandler(authHandler.ExtensionFieldsHandler)
	srv.SetAuthorizeScopeHandler(authHandler.AuthorizeScopeHandler)

	return srv, manager
}
