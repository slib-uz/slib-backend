package journal_editorial

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journaleditorialusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalEditorialListByRoleHandler struct {
	uc *journaleditorialusecases.JournalEditorialListByRoleUseCase
}

// @inject
func NewJournalEditorialListByRoleHandler(uc *journaleditorialusecases.JournalEditorialListByRoleUseCase) *JournalEditorialListByRoleHandler {
	return &JournalEditorialListByRoleHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-editorial
// @Accept       json
// @Produce      json
// @Param        journalId  path  int  true  "Journal ID"
// @Success      200  {object}  journaleditorialusecases.JournalEditorialListByRoleResponse
// @Failure      400  {object}  response.Response
// @Router       /journal-editorial/{journalId}/list-by-role [get]
func (this *JournalEditorialListByRoleHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalIdValue, err := context.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}
	journalID := uint(journalIdValue)

	result, err := this.uc.Execute(journalID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
