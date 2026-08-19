package news

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/newsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type NewsListHandler struct {
	uc *NewsListUseCase
}

// @inject
func NewNewsListHandler(uc *NewsListUseCase) *NewsListHandler {
	return &NewsListHandler{uc: uc}
}

// Handle NewsListHandler
// @Tags news
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param ordering query string false "Ordering by field" Enums(created_at,-created_at)
// @Param category_id query int false "Filter by category ID"
// @Success 200 {array} entity.NewsEntity
// @Failure 400 {object} response.Response
// @Router /news/list [get]
func (this *NewsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	ordering := ctx.QueryParam("ordering")
	categoryID := context2.GetIntQueryParam(c, "category_id", 0)
	paging, err := this.uc.Execute(page, pageSize, ordering, uint(categoryID))

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
