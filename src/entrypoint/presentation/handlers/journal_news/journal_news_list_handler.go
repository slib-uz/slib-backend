package journal_news

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journalnewsusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalNewsListHandler struct {
	uc *journalnewsusecases.JournalNewsListUseCase
}

// @inject
func NewJournalNewsListHandler(uc *journalnewsusecases.JournalNewsListUseCase) *JournalNewsListHandler {
	return &JournalNewsListHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-news
// @Accept       json
// @Produce      json
// @Param        journalId  path   int  true   "Journal ID"
// @Param        page       query  int  false  "Page number"
// @Param        page_size  query  int  false  "Items per page"
// @Success      200  {array}  entity.JournalNewsEntity
// @Failure      400  {object}  response.Response
// @Router       /journal-news/{journalId}/list [get]
func (this *JournalNewsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalIdValue, err := context.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}
	journalID := uint(journalIdValue)

	page, pageSize := context.GetPagingParams(c)

	result, err := this.uc.Execute(journalID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
