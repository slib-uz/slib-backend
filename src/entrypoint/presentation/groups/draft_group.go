package groups

import (
	"github.com/labstack/echo/v4"
	draft2 "slib.uz/src/entrypoint/presentation/handlers/draft"
)

type DraftGroup struct {
	getHandler  *draft2.GetDraftHandler
	saveHandler *draft2.DraftSaveHandler
}

// @inject
func NewDraftGroup(getHandler *draft2.GetDraftHandler, saveHandler *draft2.DraftSaveHandler) *DraftGroup {
	return &DraftGroup{getHandler: getHandler, saveHandler: saveHandler}
}

func (this *DraftGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/get", this.getHandler.Handle)
	group.POST("/save", this.saveHandler.Handle)
}
