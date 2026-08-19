package articles

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/articleusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleListOnlyNameHandler struct {
	uc *ArticlesOnlyNameListUseCase
}

// @inject
func NewArticleListOnlyNameHandler(uc *ArticlesOnlyNameListUseCase) *ArticleListOnlyNameHandler {
	return &ArticleListOnlyNameHandler{uc: uc}
}

// Handle
// @Tags         article
// @Accept       json
// @Produce      json
// @Summary Get Article List Only Name
// @Description Get Article List Only Name
// @Param full_name query string false "Full Name"
// @Param page query int false "Page"
// @Param pageSize query int false "Page Size"
// @Success 200 {array} entity.ArticleOnlyNameEntity
// @Failure 500 {object} response.Response
// @Router /articles/list-only-name [get]
func (this *ArticleListOnlyNameHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	page, pageSize := context.GetPagingParams(c)
	fullName := ctx.QueryParam("full_name")

	result, err := this.uc.Execute(fullName, page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
