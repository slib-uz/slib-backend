package journal_editorial

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journaleditorialusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalEditorialDeleteHandler struct {
	uc *journaleditorialusecases.JournalEditorialDeleteUseCase
}

// @inject
func NewJournalEditorialDeleteHandler(uc *journaleditorialusecases.JournalEditorialDeleteUseCase) *JournalEditorialDeleteHandler {
	return &JournalEditorialDeleteHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-editorial
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Editorial ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal-editorial/{id}/delete [delete]
func (this *JournalEditorialDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	idValue, err := context.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}
	id := uint(idValue)

	if err := this.uc.Execute(id, c.User); err != nil {
		return err
	}

	return c.JsonResponse(200, "Journal editorial deleted successfully")
}
