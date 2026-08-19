package middlewares

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ClientBasicAuthMiddleware struct {
	cfg conf.ConfigAdapter
}

// @inject
func NewClientBasicAuthMiddleware(cfg conf.ConfigAdapter) *ClientBasicAuthMiddleware {
	return &ClientBasicAuthMiddleware{cfg: cfg}
}

func (this *ClientBasicAuthMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*context.Context)
		username, password, ok := ctx.Request().BasicAuth()
		if !ok {
			return ctx.JSON(response.UnauthorizedError.Status, response.UnauthorizedError)
		}

		cfgUsername, cfgSecret := this.cfg.GetClientBasicAuthCredentials()

		if username != cfgUsername || password != cfgSecret {
			return ctx.JSON(response.UnauthorizedError.Status, response.UnauthorizedError)
		}

		client := entity.NewClientBasicAuth(username, password)

		ctx.SetClient(client)

		return next(ctx)
	}
}
