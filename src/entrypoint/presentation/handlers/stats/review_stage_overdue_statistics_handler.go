package stats

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/statisticsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ReviewStageOverdueStatisticsHandler struct {
	uc *statisticsusecases.ReviewStageOverdueStatisticsUseCase
}

// @inject
func NewReviewStageOverdueStatisticsHandler(uc *statisticsusecases.ReviewStageOverdueStatisticsUseCase) *ReviewStageOverdueStatisticsHandler {
	return &ReviewStageOverdueStatisticsHandler{uc: uc}
}

// Handle godoc
// @Tags stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param journal_id query int false "Filter by journal ID"
// @Success 200 {object} response.Response
// @Router /stats/review-stage/overdue [get]
func (this *ReviewStageOverdueStatisticsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	var journalID *uint
	if val := ctx.QueryParam("journal_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		rid := uint(id)
		journalID = &rid
	}

	stats, err := this.uc.Execute(journalID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, stats)
}
