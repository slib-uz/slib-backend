package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/antiplag"
)

type AntiPlagGroup struct {
	antiplagResultsHandler *antiplag.AntiPlagResultsHandler
}

// @inject
func NewAntiPlagGroup(antiplagResultsHandler *antiplag.AntiPlagResultsHandler) *AntiPlagGroup {
	return &AntiPlagGroup{antiplagResultsHandler: antiplagResultsHandler}
}

func (this *AntiPlagGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/results", this.antiplagResultsHandler.Handle)
}
