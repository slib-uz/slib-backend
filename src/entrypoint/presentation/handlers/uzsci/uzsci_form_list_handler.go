package uzsci

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/uzsciusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type UzSciFormListHandler struct {
	uc *uzsciusecases.UzSciFormListUseCase
}

// @inject
func NewUzSciFormListHandler(uc *uzsciusecases.UzSciFormListUseCase) *UzSciFormListHandler {
	return &UzSciFormListHandler{uc: uc}
}

// Handle
// @Tags         UzSci
// @Accept       json
// @Produce      json
// @Param        rating_period_id query int true "Rating period ID"
// @Success      200  {array}   entity.UzSciFormEntity
// @Router       /uzsci/forms/list [get]
func (this *UzSciFormListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	ratingPeriodID, err := strconv.ParseUint(c.QueryParam("rating_period_id"), 10, 64)
	if err != nil {
		return c.JsonResponse(400, "invalid rating_period_id")
	}

	forms, err := this.uc.Execute(uint(ratingPeriodID))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, forms)
}
