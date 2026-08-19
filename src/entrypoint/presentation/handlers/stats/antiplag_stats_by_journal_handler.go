package stats

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/antiplagusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type AntiPlagStatsByJournalHandler struct {
	uc *antiplagusecases.AntiPlagStatsByJournalUseCase
}

// @inject
func NewAntiPlagStatsByJournalHandler(uc *antiplagusecases.AntiPlagStatsByJournalUseCase) *AntiPlagStatsByJournalHandler {
	return &AntiPlagStatsByJournalHandler{uc: uc}
}

// Handle godoc
// @Tags stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date in format YYYY-MM-DD"
// @Param end_date query string false "End date in format YYYY-MM-DD"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param publisherId query int false "Publisher ID"
// @Success 200 {array} entity.JournalAntiPlagStatsEntity
// @Router /stats/antiplag [get]
func (this *AntiPlagStatsByJournalHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	publisherId := context2.GetIntQueryParam(ctx, "publisherId", 0)
	startDate, endDate := context2.GetDateRangeQueryParams(ctx, "start_date", "end_date")
	page, pageSize := context2.GetPagingParams(ctx)

	stats, err := this.uc.Execute(c.User, uint(publisherId), startDate, endDate, page, pageSize)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, stats)
}
