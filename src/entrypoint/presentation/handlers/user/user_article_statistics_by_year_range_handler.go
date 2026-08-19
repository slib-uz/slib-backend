package user

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UserArticleStatisticsByYearRangeHandler struct {
	uc *userusecases.UserArticleStatisticsByYearRangeUseCase
}

// @inject
func NewUserArticleStatisticsByYearRangeHandler(uc *userusecases.UserArticleStatisticsByYearRangeUseCase) *UserArticleStatisticsByYearRangeHandler {
	return &UserArticleStatisticsByYearRangeHandler{uc: uc}
}

// Handle
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        from query int false "From Year"
// @Param        to query int false "To Year"
// @Success      200  {object}  entity.UserArticleStatisticsByYearRangeEntity
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /user/article/statistic/by-year-range [get]
func (this *UserArticleStatisticsByYearRangeHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	fromYearStr := ctx.QueryParam("from")
	toYearStr := ctx.QueryParam("to")

	if (fromYearStr == "" && toYearStr != "") || (fromYearStr != "" && toYearStr == "") {
		return c.JsonResponse(400, "Both 'from' and 'to' parameters must be provided together, or neither")
	}

	var fromYear, toYear int
	var err error

	if fromYearStr != "" && toYearStr != "" {
		fromYear, err = strconv.Atoi(fromYearStr)
		if err != nil {
			return c.JsonResponse(400, "Invalid 'from' year parameter")
		}

		toYear, err = strconv.Atoi(toYearStr)
		if err != nil {
			return c.JsonResponse(400, "Invalid 'to' year parameter")
		}
	}

	userID := c.User.ID

	statistics, err := this.uc.Execute(userID, fromYear, toYear)
	if err != nil {
		return c.JsonResponse(400, err.Error())
	}

	return c.JsonResponse(200, statistics)
}
