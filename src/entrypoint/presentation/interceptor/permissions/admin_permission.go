package permissions

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/entrypoint/presentation/app/context"
)

func AdminPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := ctx.(*context.Context)

		if c.User != nil && c.User.IsAdmin {
			return next(ctx)
		}
		return response.PermissionDeniedError
	}
}
