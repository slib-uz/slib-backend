package articles

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/articleusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleDownloadHandler struct {
	uc *articleusecases.ArticleDownloadUseCase
}

// @inject
func NewArticleDownloadHandler(uc *articleusecases.ArticleDownloadUseCase) *ArticleDownloadHandler {
	return &ArticleDownloadHandler{uc: uc}
}

// Handle
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        articleId path int true "Article ID"
// @Success      200 {object} entity.ArticleDownloadEntity
// @Router       /articles/{articleId}/download [get]
// @Security     BasicAuth
func (this *ArticleDownloadHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	articleIdValue, err := context2.GetIntPathParam(c, "articleId")
	if err != nil {
		return err
	}
	articleId := uint(articleIdValue)

	article, err := this.uc.Execute(articleId)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, article)

}
