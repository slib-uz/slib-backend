package middlewares

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ClientBasicAuthDBMiddleware struct {
	clientAuthService *service.ClientAuthService
}

// @inject
func NewClientBasicAuthDBMiddleware(clientAuthService *service.ClientAuthService) *ClientBasicAuthDBMiddleware {
	return &ClientBasicAuthDBMiddleware{
		clientAuthService: clientAuthService,
	}
}

func (m *ClientBasicAuthDBMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*context.Context)

		clientID, clientSecret, ok := ctx.Request().BasicAuth()
		if !ok {
			return ctx.JSON(response.UnauthorizedError.Status, response.UnauthorizedError)
		}

		client, err := m.clientAuthService.Authenticate(clientID, clientSecret)
		if err != nil {
			return ctx.JSON(response.UnauthorizedError.Status, response.UnauthorizedError)
		}

		ctx.SetClient(client)

		return next(ctx)
	}
}
