package journal

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalusecases"
	"slib.uz/src/core/domain/entity"
	context "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalUpdateHandler struct {
	uc *JournalUpdateUseCase
}

// @inject
func NewJournalUpdateHandler(uc *JournalUpdateUseCase) *JournalUpdateHandler {
	return &JournalUpdateHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-manage
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        journalId  path      int  true  "Journal ID"
// @Param        request       body      entity.JournalCreateEntity  true  "Journal update data"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal-manage/update/{journalId} [put]
func (this *JournalUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[entity.JournalCreateEntity](ctx)
	if err != nil {
		return err
	}

	journalId, err := context.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(uint(journalId), data); err != nil {
		return err
	}

	return c.JsonResponse(200, "Journal updated successfully")
}
