package articles

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/statisticsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleStatisticsByStudyFieldHandler struct {
	uc *statisticsusecases.ArticleStatisticsByStudyFieldUseCase
}

// @inject
func NewArticleStatisticsByStudyFieldHandler(uc *statisticsusecases.ArticleStatisticsByStudyFieldUseCase) *ArticleStatisticsByStudyFieldHandler {
	return &ArticleStatisticsByStudyFieldHandler{uc: uc}
}

// Handle godoc
// @Tags article
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param journal_id query int false "Filter by journal ID"
// @Success 200 {array} entity.ArticleStatisticsByStudyFieldEntity
// @Router /articles/statistics/by-study-field [get]
func (this *ArticleStatisticsByStudyFieldHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)

	var journalID *uint
	if val := ctx.QueryParam("journal_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		rid := uint(id)
		journalID = &rid
	}

	stats, err := this.uc.Execute(page, pageSize, journalID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, stats)
}
