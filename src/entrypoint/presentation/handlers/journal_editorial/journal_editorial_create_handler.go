package journal_editorial

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journaleditorialusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/journal_editorial/schema"
)

type JournalEditorialCreateHandler struct {
	uc *journaleditorialusecases.JournalEditorialCreateUseCase
}

// @inject
func NewJournalEditorialCreateHandler(uc *journaleditorialusecases.JournalEditorialCreateUseCase) *JournalEditorialCreateHandler {
	return &JournalEditorialCreateHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-editorial
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        journalId  path  int                             true  "Journal ID"
// @Param        request    body  schema.JournalEditorialCreateOrUpdateRequest  true  "Editorial data"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal-editorial/{journalId}/create [post]
func (this *JournalEditorialCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[schema.JournalEditorialCreateOrUpdateRequest](ctx)
	if err != nil {
		return err
	}

	journalIdValue, err := context.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}

	editorial, err := data.ToEntity()
	if err != nil {
		return err
	}
	editorial.JournalID = uint(journalIdValue)

	if err := this.uc.Execute(editorial, c.User); err != nil {
		return err
	}

	return c.JsonResponse(201, "Journal editorial created successfully")
}
