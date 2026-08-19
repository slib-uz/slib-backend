package groups

import (
	"github.com/labstack/echo/v4"
	guide2 "slib.uz/src/entrypoint/presentation/handlers/guide"
)

type GuidesGroup struct {
	guidesListHandler    *guide2.GuidesListHandler
	guideRetrieveHandler *guide2.GuideRetrieveHandler
}

// @inject
func NewGuideGroup(guideListHandler *guide2.GuidesListHandler, guideRetrieveHandler *guide2.GuideRetrieveHandler) *GuidesGroup {
	return &GuidesGroup{guidesListHandler: guideListHandler, guideRetrieveHandler: guideRetrieveHandler}
}

func (this *GuidesGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.guidesListHandler.Handle)
	group.GET("/retrieve/:guideId", this.guideRetrieveHandler.Handle)
}
