package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/uzsci"
)

type UzsciGroup struct {
	formListHandler          *uzsci.UzSciFormListHandler
	ratingPeriodListHandler  *uzsci.UzSciRatingPeriodListHandler
	createApplicationHandler *uzsci.UzSciCreateApplicationHandler
}

// @inject
func NewUzsciGroup(
	formListHandler *uzsci.UzSciFormListHandler,
	ratingPeriodListHandler *uzsci.UzSciRatingPeriodListHandler,
	createApplicationHandler *uzsci.UzSciCreateApplicationHandler,
) *UzsciGroup {
	return &UzsciGroup{
		formListHandler:          formListHandler,
		ratingPeriodListHandler:  ratingPeriodListHandler,
		createApplicationHandler: createApplicationHandler,
	}
}

func (this *UzsciGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/forms/list", this.formListHandler.Handle)
	group.GET("/rating-periods/public", this.ratingPeriodListHandler.Handle)
	group.POST("/rating-periods/:period_id/applications", this.createApplicationHandler.Handle)
}
