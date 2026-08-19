package articles

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/articleusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/articles/schema"
)

type ArticleUpdateHandler struct {
	uc *ArticleUpdateUseCase
}

// @inject
func NewArticleUpdateHandler(uc *ArticleUpdateUseCase) *ArticleUpdateHandler {
	return &ArticleUpdateHandler{uc: uc}
}

// Handle godoc
// @Tags         article
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        articleId  path  int  true  "Article ID"
// @Param        request  body  schema.ArticleUpdateRequestSchema  true  "Article update data"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /articles/update/{articleId} [put]
func (this *ArticleUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[schema.ArticleUpdateRequestSchema](ctx)
	if err != nil {
		return err
	}

	articleId, err := context.GetIntPathParam(ctx, "articleId")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, uint(articleId), data.ToEntity()); err != nil {
		return err
	}

	return c.JsonResponse(200, "Article updated successfully")
}
