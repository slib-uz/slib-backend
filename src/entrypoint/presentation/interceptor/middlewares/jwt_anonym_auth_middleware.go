package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/service"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/headers"
)

type JwAnonymAuthMiddleware struct {
	anonymousService *service.AnonymousUserTokenService
}

// @inject
func NewJwAnonymAuthMiddleware(anonymousService *service.AnonymousUserTokenService) *JwAnonymAuthMiddleware {
	return &JwAnonymAuthMiddleware{anonymousService: anonymousService}
}

func (this *JwAnonymAuthMiddleware) Call(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := ctx.(*context.Context)

		if c.User != nil {
			return next(c)
		}
		token := this.checkToken(c)

		c.Response().Before(func() {
			c.Response().Header().Set(headers.HeaderAnonymToken, token)
		})

		return next(c)
	}
}

func (this *JwAnonymAuthMiddleware) checkToken(c *context.Context) string {
	token := c.Request().Header.Get(headers.HeaderAnonymToken)

	if token == "" {
		return this.generateNewAnonymToken(c)
	}

	subject, err := this.anonymousService.VerifyToken(token)
	if err == nil {
		c.AnonymousID = subject
		return token
	}

	return this.generateNewAnonymToken(c)
}

func (this *JwAnonymAuthMiddleware) generateNewAnonymToken(c *context.Context) string {
	anonymID := uuid.NewString()
	c.AnonymousID = anonymID
	return this.anonymousService.GenerateToken(anonymID)
}
