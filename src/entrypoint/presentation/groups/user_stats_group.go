package groups

import (
	"github.com/labstack/echo/v4"
	user2 "slib.uz/src/entrypoint/presentation/handlers/user"
)

type UserStatsGroup struct {
	genderStatisticsHandler *user2.UserGenderStatisticsHandler
	ageStatisticsHandler    *user2.UserAgeStatisticsHandler
}

// @inject
func NewUserStatsGroup(genderStatisticsHandler *user2.UserGenderStatisticsHandler, ageStatisticsHandler *user2.UserAgeStatisticsHandler) *UserStatsGroup {
	return &UserStatsGroup{genderStatisticsHandler: genderStatisticsHandler, ageStatisticsHandler: ageStatisticsHandler}
}

func (this *UserStatsGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/gender/statistic", this.genderStatisticsHandler.Handle)
	group.GET("/age/statistic", this.ageStatisticsHandler.Handle)
}
