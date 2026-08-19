package groups

import (
	"github.com/hibiken/asynqmon"
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/docs"
)

type DevelopmentGroup struct {
	asynqmonHandler *asynqmon.HTTPHandler
}

// @inject
func NewDevelopmentGroup(asynqmonHandler *asynqmon.HTTPHandler) *DevelopmentGroup {
	return &DevelopmentGroup{asynqmonHandler: asynqmonHandler}
}

func (this *DevelopmentGroup) RegisterRoutes(group *echo.Group, middlewares ...echo.MiddlewareFunc) {
	group.GET("/docs/*", docs.SwaggerWrapHandler, middlewares...)
	group.Any("/monitoring/tasks/*", echo.WrapHandler(this.asynqmonHandler), middlewares...)
}
