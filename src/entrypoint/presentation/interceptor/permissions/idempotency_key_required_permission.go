package permissions

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/headers"
)

func IdempotencyKeyRequiredPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		key := ctx.Request().Header.Get(headers.HeaderIdempotencyKey)

		if key == "" {
			return ctx.JSON(400, map[string]any{"error": "Idempotency key is required"})
		}

		return next(ctx)
	}
}
