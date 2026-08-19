package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/userfind"
)

//var UserFindGroup = app.GetInstance().Group("/users/find")
//
//func init() {
//	UserFindGroup.GET("", container.InitUserFindHandler().Handle)
//}

type UserFindGroup struct {
	find *userfind.UserFindHandler
}

// @inject
func NewUserFindGroup(find *userfind.UserFindHandler) *UserFindGroup {
	return &UserFindGroup{find: find}
}

func (this *UserFindGroup) RegisterRoutes(group *echo.Group) {
	group.GET("", this.find.Handle)
}
