package user

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UserGenderStatisticsHandler struct {
	uc *userusecases.UserGenderStatisticsUseCase
}

// @inject
func NewUserGenderStatisticsHandler(uc *userusecases.UserGenderStatisticsUseCase) *UserGenderStatisticsHandler {
	return &UserGenderStatisticsHandler{uc: uc}
}

// Handle godoc
// @Tags UserStats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.UserGenderStatisticsEntity
// @Router /user-stats/gender/statistic [get]
func (this *UserGenderStatisticsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	stats, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, stats)
}
