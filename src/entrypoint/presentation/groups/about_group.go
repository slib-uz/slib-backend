package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/about"
)

type AboutGroup struct {
	aboutListHandler     *about.AboutListHandler
	aboutRetrieveHandler *about.AboutRetrieveHandler
}

// @inject
func NewAboutGroup(aboutListHandler *about.AboutListHandler, aboutRetrieveHandler *about.AboutRetrieveHandler) *AboutGroup {
	return &AboutGroup{aboutListHandler: aboutListHandler, aboutRetrieveHandler: aboutRetrieveHandler}
}

func (this *AboutGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.aboutListHandler.Handle)
	group.GET("/detail/:id", this.aboutRetrieveHandler.Handle)
}
