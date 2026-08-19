package groups

import (
	"github.com/labstack/echo/v4"
	role2 "slib.uz/src/entrypoint/presentation/handlers/role"
)

//var RoleGroup = app.GetInstance().Group(
//	"/role",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AuthenticatedPermission,
//	//permissions.AdminPermission,
//)

//func init() {
//	RoleGroup.POST("/create", container.InitUserRoleCreateHandler().Handle)
//	RoleGroup.GET("/user-roles", container.InitUserRoleListHandler().Handle)
//	RoleGroup.DELETE("/delete/:id", container.InitUserRoleDeleteHandler().Handle)
//}

type RoleGroup struct {
	createHandler *role2.UserRoleCreateHandler
	listHandler   *role2.UserRoleListHandler
	deleteHandler *role2.RoleDeleteHandler
}

// @inject
func NewRoleGroup(createHandler *role2.UserRoleCreateHandler, listHandler *role2.UserRoleListHandler, deleteHandler *role2.RoleDeleteHandler) *RoleGroup {
	return &RoleGroup{createHandler: createHandler, listHandler: listHandler, deleteHandler: deleteHandler}
}

func (this *RoleGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create", this.createHandler.Handle)
	group.GET("/user-roles", this.listHandler.Handle)
	group.DELETE("/delete/:id", this.deleteHandler.Handle)
}
