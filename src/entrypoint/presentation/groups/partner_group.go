package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/partner"
)

type PartnerGroup struct {
	partnersListHandler *partner.PartnersListHandler
}

// @inject
func NewPartnerGroup(partnerListHandler *partner.PartnersListHandler) *PartnerGroup {
	return &PartnerGroup{partnersListHandler: partnerListHandler}
}

func (this *PartnerGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.partnersListHandler.Handle)
}
