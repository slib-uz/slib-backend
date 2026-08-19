package articles

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/articleusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type TopArticlesHandler struct {
	uc *TopArticlesUseCase
}

// @inject
func NewTopArticlesHandler(uc *TopArticlesUseCase) *TopArticlesHandler {
	return &TopArticlesHandler{uc: uc}
}

// Handle
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        page         query     int false  "Page number for pagination"
// @Param        page_size    query     int false  "Number of items per page for pagination"
// @Success      200  {array}  entity.ArticleEntity
// @Failure      400  {object}  response.Response
// @Router       /articles/top [get]
func (this *TopArticlesHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(c)

	result, err := this.uc.Execute(page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
