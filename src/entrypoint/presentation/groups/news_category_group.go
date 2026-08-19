package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/news_category"
)

type NewsCategoryGroup struct {
	list *news_category.NewsCategoryListHandler
}

// @inject
func NewNewsCategoryGroup(list *news_category.NewsCategoryListHandler) *NewsCategoryGroup {
	return &NewsCategoryGroup{list: list}
}

func (this *NewsCategoryGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.list.Handle)
}
