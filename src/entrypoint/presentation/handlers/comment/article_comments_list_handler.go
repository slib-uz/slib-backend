package comment

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/commentusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleCommentsListHandler struct {
	uc *ArticleCommentsListUseCase
}

// @inject
func NewArticleCommentsListHandler(uc *ArticleCommentsListUseCase) *ArticleCommentsListHandler {
	return &ArticleCommentsListHandler{uc: uc}
}

// Handle godoc
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        articleId  path      int  true  "Article ID"
// @Param        page       query     int  false "Page number" default(1)
// @Param        page_size  query     int  false "Page size" default(10)
// @Success      200  {object}  response.Response
// @Router       /comment/article/{articleId}/list [get]
func (this *ArticleCommentsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(c)

	articleIdValue, err := context2.GetIntPathParam(ctx, "articleId")
	if err != nil {
		return err
	}
	articleID := uint(articleIdValue)

	paging, err := this.uc.Execute(articleID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
