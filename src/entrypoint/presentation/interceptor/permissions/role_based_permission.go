package permissions

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/entrypoint/presentation/app/context"
)

func RoleBasedPermission(role enum.UserRole) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			c := ctx.(*context.Context)

			if c.User == nil {
				return response.UnauthorizedError
			}

			publisherStr := c.QueryParam("publisher_id")
			if publisherStr == "" {
				if hasAdminAccess(c) {
					return next(ctx)
				}
				return response.PermissionDeniedError
			}

			publisherId, err := strconv.ParseUint(publisherStr, 10, 32)
			if err != nil {
				return response.InvalidArgument
			}

			if hasAccess(c, uint(publisherId), role) {
				return next(ctx)
			}
			return response.PermissionDeniedError
		}
	}
}

func hasAdminAccess(c *context.Context) bool {
	for _, item := range c.User.Roles {
		if item.Role == enum.RoleAdmin {
			return true
		}
	}
	return false
}

func hasAccess(c *context.Context, publisherId uint, role enum.UserRole) bool {
	for _, item := range c.User.Roles {
		if *item.PublisherID == publisherId && item.Role == role {
			return true
		}
	}
	return false
}
