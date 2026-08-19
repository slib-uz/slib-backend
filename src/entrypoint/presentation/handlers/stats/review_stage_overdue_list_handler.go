package stats

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/statisticsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ReviewStageOverdueListHandler struct {
	uc *statisticsusecases.ReviewStageOverdueListUseCase
}

// @inject
func NewReviewStageOverdueListHandler(uc *statisticsusecases.ReviewStageOverdueListUseCase) *ReviewStageOverdueListHandler {
	return &ReviewStageOverdueListHandler{uc: uc}
}

// Handle godoc
// @Tags stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param journal_id query int false "Filter by journal ID"
// @Success 200 {object} response.Response
// @Router /stats/review-stage/overdue/list [get]
func (this *ReviewStageOverdueListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)

	var journalID *uint
	if val := ctx.QueryParam("journal_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		rid := uint(id)
		journalID = &rid
	}

	items, err := this.uc.Execute(page, pageSize, journalID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, items)
}
