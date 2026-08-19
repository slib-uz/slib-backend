package antiplag

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/antiplagusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type AntiPlagResultsHandler struct {
	uc *AntiPlagResultsByJournalUsecase
}

// @inject
func NewAntiPlagResultsHandler(uc *AntiPlagResultsByJournalUsecase) *AntiPlagResultsHandler {
	return &AntiPlagResultsHandler{uc: uc}
}

// Handle
// @Tags antiplag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param journal_id query int true "Journal ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /antiplag/results [get]
func (this *AntiPlagResultsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	journalID := uint(context2.GetIntQueryParam(ctx, "journal_id", 0))
	page, pageSize := context2.GetPagingParams(c)

	results, err := this.uc.Execute(journalID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, results)
}
