package journal_rating

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalratingusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalRatingListHandler struct {
	uc *JournalRatingListUseCase
}

// @inject
func NewJournalRatingListHandler(uc *JournalRatingListUseCase) *JournalRatingListHandler {
	return &JournalRatingListHandler{uc: uc}
}

// Handle JournalRatingListHandler
// @Tags journal-rating
// @Accept json
// @Produce json
// @Param journalId path int true "Journal ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param ordering query string false "Ordering" Enums(created_at,-created_at)
// @Success 200 {array} entity.JournalRatingEntity
// @Failure 400 {object} response.Response
// @Router /journal-rating/list/{journalId} [get]
func (this *JournalRatingListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	ordering := ctx.QueryParam("ordering")

	journalIdValue, err := context2.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}
	journalID := uint(journalIdValue)

	paging, err := this.uc.Execute(journalID, page, pageSize, ordering)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, paging)
}
