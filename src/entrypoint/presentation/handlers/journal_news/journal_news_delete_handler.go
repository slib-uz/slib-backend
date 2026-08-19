package journal_news

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journalnewsusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalNewsDeleteHandler struct {
	uc *journalnewsusecases.JournalNewsDeleteUseCase
}

// @inject
func NewJournalNewsDeleteHandler(uc *journalnewsusecases.JournalNewsDeleteUseCase) *JournalNewsDeleteHandler {
	return &JournalNewsDeleteHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-news
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "News ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal-news/{id}/delete [delete]
func (this *JournalNewsDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	idValue, err := context.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}
	id := uint(idValue)

	if err := this.uc.Execute(id, c.User); err != nil {
		return err
	}

	return c.JsonResponse(200, "Journal news deleted successfully")
}
