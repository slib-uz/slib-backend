package journal_editorial

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journaleditorialusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalEditorialListHandler struct {
	uc *journaleditorialusecases.JournalEditorialListUseCase
}

// @inject
func NewJournalEditorialListHandler(uc *journaleditorialusecases.JournalEditorialListUseCase) *JournalEditorialListHandler {
	return &JournalEditorialListHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-editorial
// @Accept       json
// @Produce      json
// @Param        journalId  path   int  true   "Journal ID"
// @Param        page       query  int  false  "Page number"
// @Param        page_size  query  int  false  "Items per page"
// @Success      200  {array}  entity.JournalEditorialEntity
// @Failure      400  {object}  response.Response
// @Router       /journal-editorial/{journalId}/list [get]
func (this *JournalEditorialListHandler) Handle(ctx echo.Context) error {
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
