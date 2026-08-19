package stats

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/statisticsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ReviewStageStatisticsHandler struct {
	uc *statisticsusecases.ReviewStageStatisticsUseCase
}

// @inject
func NewReviewStageStatisticsHandler(uc *statisticsusecases.ReviewStageStatisticsUseCase) *ReviewStageStatisticsHandler {
	return &ReviewStageStatisticsHandler{uc: uc}
}

// Handle godoc
// @Tags stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param journal_id query int false "Filter by journal ID"
// @Success 200 {object} response.Response
// @Router /stats/review-stage [get]
func (this *ReviewStageStatisticsHandler) Handle(ctx echo.Context) error {
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
