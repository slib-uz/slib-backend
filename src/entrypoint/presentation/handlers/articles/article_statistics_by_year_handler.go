package articles

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/statisticsusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleStatisticsByYearHandler struct {
	uc *statisticsusecases.ArticleStatisticsByYearUseCase
}

// @inject
func NewArticleStatisticsByYearHandler(uc *statisticsusecases.ArticleStatisticsByYearUseCase) *ArticleStatisticsByYearHandler {
	return &ArticleStatisticsByYearHandler{uc: uc}
}

// Handle godoc
// @Tags article
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2024). Defaults to current year"
// @Param journal_id query int false "Filter by journal ID"
// @Param publisher_id query int false "Filter by publisher ID"
// @Success 200 {object} entity.ArticleStatisticsByYearEntity
// @Router /articles/statistics/by-year [get]
func (this *ArticleStatisticsByYearHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	var year int
	if yearStr := ctx.QueryParam("year"); yearStr != "" {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return c.JsonResponse(400, "Invalid 'year' parameter")
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

	stats, err := this.uc.Execute(year, journalID, publisherID)
	if err != nil {
		return c.JsonResponse(400, err.Error())
	}

	return c.JsonResponse(200, stats)
}
