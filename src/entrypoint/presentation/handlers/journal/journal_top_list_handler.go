package journal

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalTopListHandler struct {
	uc *TopJournalListUseCase
}

// @inject
func NewJournalTopListHandler(uc *TopJournalListUseCase) *JournalTopListHandler {
	return &JournalTopListHandler{uc: uc}
}

// Handle
// @Tags         journal
// @Accept       json
// @Produce      json
// @Param        page         query     int false  "Page number for pagination"
// @Param        page_size    query     int false  "Number of items per page for pagination"
// @Success      200  {array} entity.JournalEntity
// @Failure      400  {object}  response.Response
// @Router       /journal/top [get]
func (this *JournalTopListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(c)

	result, err := this.uc.Execute(page, pageSize)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, result)
}
