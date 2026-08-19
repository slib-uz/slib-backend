package middlewares

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/app/context"
)

// ContextMiddleware is a middleware function to create a new context
func ContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.NewContext(c)
		return next(ctx)
	}
}
