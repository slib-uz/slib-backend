package permissions

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/entrypoint/presentation/app/context"
)

func RolePermission(next echo.HandlerFunc, allowedRoles ...enum.UserRole) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := ctx.(*context.Context)

		if c.User != nil {
			userRoles := c.User.Roles

			for _, allowedRole := range allowedRoles {
				for _, userRole := range userRoles {
					if userRole.Role == allowedRole {
						return next(ctx)
					}
				}
			}

		}

		return response.PermissionDeniedError
	}
}
