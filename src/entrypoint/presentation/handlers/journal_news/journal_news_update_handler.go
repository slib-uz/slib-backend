package journal_news

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journalnewsusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalNewsUpdateHandler struct {
	uc *journalnewsusecases.JournalNewsUpdateUseCase
}

// @inject
func NewJournalNewsUpdateHandler(uc *journalnewsusecases.JournalNewsUpdateUseCase) *JournalNewsUpdateHandler {
	return &JournalNewsUpdateHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-news
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  int                        true  "News ID"
// @Param        request  body  entity.JournalNewsEntity   true  "News data"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal-news/{id}/update [put]
func (this *JournalNewsUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	idValue, err := context.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}
	id := uint(idValue)

	data, err := context.GetBody[entity.JournalNewsEntity](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(id, data, c.User); err != nil {
		return err
	}

	return c.JsonResponse(200, "Journal news updated successfully")
}
