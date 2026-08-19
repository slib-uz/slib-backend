package comment

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/commentusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type CommentStatsHandler struct {
	uc *CommentStatsUsecase
}

// @inject
func NewCommentStatsHandler(uc *CommentStatsUsecase) *CommentStatsHandler {
	return &CommentStatsHandler{uc: uc}
}

// Handle Comment Stats
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        articleId  path      int  true  "Article ID"
// @Success      200  {object}  entity.CommentStatsEntity
// @Failure      400  {object}  response.Response
// @Router       /comment/article/{articleId}/stats [get]
func (this *CommentStatsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	articleID, err := context2.GetIntPathParam(ctx, "articleId")
	if err != nil {
		return err
	}

	stats, err := this.uc.Execute(uint(articleID))
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return c.JsonResponse(200, stats)
}
