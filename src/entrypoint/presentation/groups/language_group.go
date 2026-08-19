package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/language"
)

//var LanguageGroup = app.GetInstance().Group("/language")

//func init() {
//	LanguageGroup.GET("/list", container.InitLanguageListHandler().Handle)
//}

type LanguageGroup struct {
	list *language.LanguageListHandler
}

// @inject
func NewLanguageGroup(list *language.LanguageListHandler) *LanguageGroup {
	return &LanguageGroup{list: list}
}

func (this *LanguageGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.list.Handle)
}
