package docs

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SwaggerWrapHandler(c echo.Context) error {
	host := c.Request().Header.Get("X-Forwarded-Host")
	if host == "" {
		host = c.Request().Host
	}
	SwaggerInfo.Host = host

	return echoSwagger.WrapHandler(c)
}
