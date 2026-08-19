package journal_news

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journalnewsusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalNewsDetailHandler struct {
	uc *journalnewsusecases.JournalNewsDetailUseCase
}

// @inject
func NewJournalNewsDetailHandler(uc *journalnewsusecases.JournalNewsDetailUseCase) *JournalNewsDetailHandler {
	return &JournalNewsDetailHandler{uc: uc}
}

// Handle godoc
// @Tags         journal-news
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "News ID"
// @Success      200  {object}  entity.JournalNewsEntity
// @Failure      400  {object}  response.Response
// @Router       /journal-news/{id}/detail [get]
func (this *JournalNewsDetailHandler) Handle(ctx echo.Context) error {
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
