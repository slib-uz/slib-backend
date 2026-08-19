package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/site"
)

type SiteGroup struct {
	statisticsHandler *site.SiteStatisticsHandler
}

// @inject
func NewSiteGroup(statisticsHandler *site.SiteStatisticsHandler) *SiteGroup {
	return &SiteGroup{statisticsHandler: statisticsHandler}
}

func (this *SiteGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/statistic", this.statisticsHandler.Handle)
}
