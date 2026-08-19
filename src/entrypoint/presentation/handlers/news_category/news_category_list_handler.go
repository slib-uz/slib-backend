package news_category

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/newscategoryusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type NewsCategoryListHandler struct {
	uc *NewsCategoryListUseCase
}

// @inject
func NewNewsCategoryListHandler(uc *NewsCategoryListUseCase) *NewsCategoryListHandler {
	return &NewsCategoryListHandler{uc: uc}
}

// Handle
// @Tags news-category
// @Accept json
// @Produce json
// @Success 200 {array} entity.NewsCategoryEntity
// @Router /news-category/list [get]
func (this *NewsCategoryListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	categories, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, categories)
}
