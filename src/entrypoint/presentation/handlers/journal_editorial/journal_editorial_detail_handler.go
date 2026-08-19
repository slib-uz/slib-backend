package journal_editorial

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journaleditorialusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalEditorialDetailHandler struct {
	uc *journaleditorialusecases.JournalEditorialDetailUseCase
}

// @inject
func NewJournalEditorialDetailHandler(uc *journaleditorialusecases.JournalEditorialDetailUseCase) *JournalEditorialDetailHandler {
	return &JournalEditorialDetailHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-editorial
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Editorial ID"
// @Success      200  {object}  entity.JournalEditorialEntity
// @Failure      400  {object}  response.Response
// @Router       /journal-editorial/{id}/detail [get]
func (this *JournalEditorialDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	idValue, err := context.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}
	id := uint(idValue)

	result, err := this.uc.Execute(id)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
