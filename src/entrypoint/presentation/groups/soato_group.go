package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/soato"
)

type SoatoGroup struct {
	listHandler *soato.SoatoListHandler
}

// @inject
func NewSoatoGroup(listHandler *soato.SoatoListHandler) *SoatoGroup {
	return &SoatoGroup{listHandler: listHandler}
}

func (this *SoatoGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.listHandler.Handle)
}
