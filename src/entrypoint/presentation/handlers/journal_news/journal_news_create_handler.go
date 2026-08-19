package journal_news

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journalnewsusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalNewsCreateHandler struct {
	uc *journalnewsusecases.JournalNewsCreateUseCase
}

// @inject
func NewJournalNewsCreateHandler(uc *journalnewsusecases.JournalNewsCreateUseCase) *JournalNewsCreateHandler {
	return &JournalNewsCreateHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-news
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        journalId  path  int                        true  "Journal ID"
// @Param        request    body  entity.JournalNewsEntity   true  "News data"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal-news/{journalId}/create [post]
func (this *JournalNewsCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[entity.JournalNewsEntity](ctx)
	if err != nil {
		return err
	}

	journalIdValue, err := context.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}
	data.JournalID = uint(journalIdValue)

	if err := this.uc.Execute(data, c.User); err != nil {
		return err
	}

	return c.JsonResponse(201, "Journal news created successfully")
}
