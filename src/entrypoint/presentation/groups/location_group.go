package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/location"
)

type LocationGroup struct {
	regionListHandler   *location.RegionListHandler
	districtListHandler *location.DistrictListHandler
}

// @inject
func NewLocationGroup(regionListHandler *location.RegionListHandler, districtListHandler *location.DistrictListHandler) *LocationGroup {
	return &LocationGroup{regionListHandler: regionListHandler, districtListHandler: districtListHandler}
}

func (this *LocationGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/regions/list", this.regionListHandler.Handle)
	group.GET("/districts/list", this.districtListHandler.Handle)
}
