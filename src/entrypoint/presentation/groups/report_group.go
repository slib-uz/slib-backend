package groups

import (
	"github.com/labstack/echo/v4"
	report2 "slib.uz/src/entrypoint/presentation/handlers/report"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type ReportGroup struct {
	reportCreateHandler *report2.ReportCreateHandler
	reportListHandler   *report2.ReportListHandler
}

// @inject
func NewReportGroup(reportCreateHandler *report2.ReportCreateHandler, reportListHandler *report2.ReportListHandler) *ReportGroup {
	return &ReportGroup{reportCreateHandler: reportCreateHandler, reportListHandler: reportListHandler}
}

func (this ReportGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create", this.reportCreateHandler.Handle)

	adminGroup := group.Group("", permissions.AdminPermission)
	adminGroup.GET("/list", this.reportListHandler.Handle)
}
