package user

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UserAgeStatisticsHandler struct {
	uc *userusecases.UserAgeStatisticsUseCase
}

// @inject
func NewUserAgeStatisticsHandler(uc *userusecases.UserAgeStatisticsUseCase) *UserAgeStatisticsHandler {
	return &UserAgeStatisticsHandler{uc: uc}
}

// Handle godoc
// @Tags UserStats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.UserAgeStatisticsEntity
// @Router /user-stats/age/statistic [get]
func (this *UserAgeStatisticsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	stats, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, stats)
}
