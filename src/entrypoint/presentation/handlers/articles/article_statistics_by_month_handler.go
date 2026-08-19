package articles

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/statisticsusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleStatisticsByMonthHandler struct {
	uc *statisticsusecases.ArticleStatisticsByMonthUseCase
}

// @inject
func NewArticleStatisticsByMonthHandler(uc *statisticsusecases.ArticleStatisticsByMonthUseCase) *ArticleStatisticsByMonthHandler {
	return &ArticleStatisticsByMonthHandler{uc: uc}
}

// Handle godoc
// @Tags article
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2024)"
// @Param month query int false "Month number (1-12)"
// @Param journal_id query int false "Filter by journal ID"
// @Param publisher_id query int false "Filter by publisher ID"
// @Success 200 {object} entity.ArticleStatisticsByMonthEntity
// @Router /articles/statistics/by-month [get]
func (this *ArticleStatisticsByMonthHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	yearStr := ctx.QueryParam("year")
	monthStr := ctx.QueryParam("month")

	if (yearStr == "" && monthStr != "") || (yearStr != "" && monthStr == "") {
		return c.JsonResponse(400, "Both 'year' and 'month' parameters must be provided together, or neither")
	}

	var year, month int
	var err error

	if yearStr != "" && monthStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return c.JsonResponse(400, "Invalid 'year' parameter")
		}

		month, err = strconv.Atoi(monthStr)
		if err != nil {
			return c.JsonResponse(400, "Invalid 'month' parameter")
		}
	}

	var journalID *uint
	if val := ctx.QueryParam("journal_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		rid := uint(id)
		journalID = &rid
	}

	var publisherID *uint
	if val := ctx.QueryParam("publisher_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		rid := uint(id)
		publisherID = &rid
	}

	stats, err := this.uc.Execute(year, month, journalID, publisherID)
	if err != nil {
		return c.JsonResponse(400, err.Error())
	}

	return c.JsonResponse(200, stats)
}
