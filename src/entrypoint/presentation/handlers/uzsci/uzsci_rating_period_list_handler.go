package uzsci

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/uzsciusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type UzSciRatingPeriodListHandler struct {
	uc *uzsciusecases.UzSciRatingPeriodListUseCase
}

// @inject
func NewUzSciRatingPeriodListHandler(uc *uzsciusecases.UzSciRatingPeriodListUseCase) *UzSciRatingPeriodListHandler {
	return &UzSciRatingPeriodListHandler{uc: uc}
}

// Handle
// @Tags         UzSci
// @Accept       json
// @Produce      json
// @Param        is_active  query  bool  false  "Filter by active status"
// @Success      200  {array}   entity.UzSciRatingPeriodEntity
// @Router       /uzsci/rating-periods/public [get]
func (this *UzSciRatingPeriodListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	var isActive *bool
	if raw := c.QueryParam("is_active"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return response.NewFailResponse(400, "invalid is_active query param")
		}
		isActive = &parsed
	}

	periods, err := this.uc.Execute(isActive)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, periods)
}
