package journal_rating

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalratingusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalRatingStatsHandler struct {
	uc *JournalRatingStatsUseCase
}

// @inject
func NewJournalRatingStatsHandler(uc *JournalRatingStatsUseCase) *JournalRatingStatsHandler {
	return &JournalRatingStatsHandler{uc: uc}
}

// Handle JournalRatingStatsHandler
// @Tags journal-rating
// @Accept json
// @Produce json
// @Param journalId path int true "Journal ID"
// @Success 200 {object} entity.JournalRatingStatsEntity
// @Failure 400 {object} response.Response
// @Router /journal-rating/stats/{journalId} [get]
func (this *JournalRatingStatsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	journalID, err := context2.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}

	stats, err := this.uc.Execute(uint(journalID))
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return c.JsonResponse(200, stats)
}
