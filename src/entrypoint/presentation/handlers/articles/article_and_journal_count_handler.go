package articles

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/articleusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleAndJournalCountHandler struct {
	uc *articleusecases.ArticleAndJournalCountUseCase
}

// @inject
func NewArticleAndJournalCountHandler(uc *articleusecases.ArticleAndJournalCountUseCase) *ArticleAndJournalCountHandler {
	return &ArticleAndJournalCountHandler{uc: uc}
}

// Handle
// @Tags         article
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.CountStatsEntity
// @Failure      400  {object}  response.Response
// @Router       /articles/stats [get]
func (this *ArticleAndJournalCountHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	stats, err := this.uc.Execute()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JsonResponse(200, stats)
}
