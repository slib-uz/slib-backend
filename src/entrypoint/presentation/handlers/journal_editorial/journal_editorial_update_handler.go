package journal_editorial

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journaleditorialusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/journal_editorial/schema"
)

type JournalEditorialUpdateHandler struct {
	uc *journaleditorialusecases.JournalEditorialUpdateUseCase
}

// @inject
func NewJournalEditorialUpdateHandler(uc *journaleditorialusecases.JournalEditorialUpdateUseCase) *JournalEditorialUpdateHandler {
	return &JournalEditorialUpdateHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-editorial
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  int                             true  "Editorial ID"
// @Param        request  body  schema.JournalEditorialCreateOrUpdateRequest  true  "Editorial data"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal-editorial/{id}/update [put]
func (this *JournalEditorialUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	idValue, err := context.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}
	id := uint(idValue)

	data, err := context.GetBody[schema.JournalEditorialCreateOrUpdateRequest](ctx)
	if err != nil {
		return err
	}

	editorial, err := data.ToEntity()
	if err != nil {
		return err
	}

	if err := this.uc.Execute(id, editorial, c.User); err != nil {
		return err
	}

	return c.JsonResponse(200, "Journal editorial updated successfully")
}
