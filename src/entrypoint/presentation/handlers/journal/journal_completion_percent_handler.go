package journal

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalCompletionPercentHandler struct {
	uc *JournalCompletionPercentUseCase
}

// @inject
func NewJournalCompletionPercentHandler(uc *JournalCompletionPercentUseCase) *JournalCompletionPercentHandler {
	return &JournalCompletionPercentHandler{uc: uc}
}

// Handle
// @Tags         journal
// @Accept       json
// @Produce      json
// @Param        id path uint true "Journal ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal/{id}/completion-percent [get]
func (this *JournalCompletionPercentHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	journalID, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	result, err := this.uc.Execute(uint(journalID))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
