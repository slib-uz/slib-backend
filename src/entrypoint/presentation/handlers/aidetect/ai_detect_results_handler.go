package aidetect

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/aidetectusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type AiDetectResultsHandler struct {
	uc *AiDetectResultsByJournalUsecase
}

// @inject
func NewAiDetectResultsHandler(uc *AiDetectResultsByJournalUsecase) *AiDetectResultsHandler {
	return &AiDetectResultsHandler{uc: uc}
}

// Handle
// @Tags         ai-detect
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        journal_id query int true "Journal ID"
// @Param        page       query int false "Page number" default(1)
// @Param        page_size  query int false "Page size" default(10)
// @Success      200 {object} response.Response
// @Router       /ai-detect/results [get]
func (this *AiDetectResultsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalID := uint(context.GetIntQueryParam(ctx, "journal_id", 0))
	page, pageSize := context.GetPagingParams(c)

	results, err := this.uc.Execute(journalID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, results)
}
