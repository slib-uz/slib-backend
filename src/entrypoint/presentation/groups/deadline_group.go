package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/deadline"
)

type DeadlineGroup struct {
	extendDeadlineHandler *deadline.ExtendDeadlineHandler
}

// @inject
func NewDeadlineGroup(extendDeadlineHandler *deadline.ExtendDeadlineHandler) *DeadlineGroup {
	return &DeadlineGroup{extendDeadlineHandler: extendDeadlineHandler}
}

func (this *DeadlineGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/extend", this.extendDeadlineHandler.Handle)
}
