package groups

import "github.com/labstack/echo/v4"

type Group interface {
	RegisterRoutes(group *echo.Group)
}
