package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/auth"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type AuthGroup struct {
	authURLHandler          *auth.LoginURLHandler
	authLoginHandler        *auth.UserLoginHandler
	testTokenHandler        *auth.TestTokenHandler
	refreshTokenHandler     *auth.RefreshTokenHandler
	orcidAuthorizUrihandler *auth.OrcidAuthorizeURLHandler
	orcidConnectHandler     *auth.OrcidConnectHandler
	logoutHandler           *auth.LogoutHandler

	// middlewares
	developerUserBasicAuthMiddleware *middlewares.DeveloperUserBasicAuthMiddleware
}

// @inject
func NewAuthGroup(
	authURLHandler *auth.LoginURLHandler,
	authLoginHandler *auth.UserLoginHandler,
	testTokenHandler *auth.TestTokenHandler,
	refreshTokenHandler *auth.RefreshTokenHandler,
	orcidAuthorizUrihandler *auth.OrcidAuthorizeURLHandler,
	orcidConnectHandler *auth.OrcidConnectHandler,
	logoutHandler *auth.LogoutHandler,

	developerUserBasicAuthMiddleware *middlewares.DeveloperUserBasicAuthMiddleware,

) *AuthGroup {
	return &AuthGroup{
		authURLHandler:          authURLHandler,
		authLoginHandler:        authLoginHandler,
		testTokenHandler:        testTokenHandler,
		refreshTokenHandler:     refreshTokenHandler,
		orcidAuthorizUrihandler: orcidAuthorizUrihandler,
		orcidConnectHandler:     orcidConnectHandler,
		logoutHandler:           logoutHandler,

		developerUserBasicAuthMiddleware: developerUserBasicAuthMiddleware,
	}
}

func (this *AuthGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/oauth/authorize/url", this.authURLHandler.Handle)
	group.POST("/oauth/login", this.authLoginHandler.Handle)
	group.GET("/test/token/:id", this.testTokenHandler.Handle, this.developerUserBasicAuthMiddleware.Wrap)
	group.POST("/refresh-token", this.refreshTokenHandler.Handle)
	group.POST("/logout", this.logoutHandler.Handle, permissions.AuthenticatedPermission)

	group.GET("/orcid/authorize/url", this.orcidAuthorizUrihandler.Handle)
	group.POST("/orcid/connect", this.orcidConnectHandler.Handle, permissions.AuthenticatedPermission)
}
