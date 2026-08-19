package user

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UserArticleStatisticsByYearHandler struct {
	uc *userusecases.UserArticleStatisticsByYearUseCase
}

// @inject
func NewUserArticleStatisticsByYearHandler(uc *userusecases.UserArticleStatisticsByYearUseCase) *UserArticleStatisticsByYearHandler {
	return &UserArticleStatisticsByYearHandler{uc: uc}
}

// Handle
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        year path int true "Year"
// @Success      200  {object}  entity.UserArticleStatisticsByYearEntity
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /user/article/statistic/by-year/{year} [get]
func (this *UserArticleStatisticsByYearHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	// Get year from path parameter
	yearStr := ctx.Param("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.JsonResponse(400, "Invalid year parameter")
	}

	// Get user ID from context (authenticated user)
	userID := c.User.ID

	statistics, err := this.uc.Execute(userID, year)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, statistics)
}
