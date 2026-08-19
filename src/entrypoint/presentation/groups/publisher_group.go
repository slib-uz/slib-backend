package groups

import (
	"github.com/labstack/echo/v4"
	publisher2 "slib.uz/src/entrypoint/presentation/handlers/publisher"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

//import (
//	container "slib.uz/src/di"
//	"slib.uz/src/presentation/app"
//	"slib.uz/src/presentation/interceptor/permissions"
//)
//
//var PublisherGroup = app.GetInstance().Group(
//	"/publisher",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AuthenticatedPermission,
//	//permissions.AdminPermission,
//)

//func init() {
//	PublisherGroup.GET("/list", container.InitPublisherListHandler().Handle)
//	PublisherGroup.POST("/create", container.InitPublisherCreateHandler().Handle, permissions.AdminPermission)
//	PublisherGroup.GET("/detail/:id", container.InitPublisherDetailHandler().Handle)
//	PublisherGroup.PUT("/update/:id", container.InitPublisherUpdateHandler().Handle, permissions.AdminPermission)
//	PublisherGroup.PATCH("/update/:id", container.InitPublisherUpdateHandler().Handle, permissions.AdminPermission)
//	PublisherGroup.GET("/:publisherId/admins", container.InitPublisherAdminHandler().Handle, permissions.AdminPermission)
//}

type PublisherGroup struct {
	listHandler   *publisher2.PublisherListHandler
	createHandler *publisher2.PublisherCreateHandler
	detailHandler *publisher2.PublisherDetailHandler
	updateHandler *publisher2.PublisherUpdateHandler
	adminsHandler *publisher2.PublisherAdminHandler
}

// @inject
func NewPublisherGroup(listHandler *publisher2.PublisherListHandler, createHandler *publisher2.PublisherCreateHandler, detailHandler *publisher2.PublisherDetailHandler, updateHandler *publisher2.PublisherUpdateHandler, adminsHandler *publisher2.PublisherAdminHandler) *PublisherGroup {
	return &PublisherGroup{listHandler: listHandler, createHandler: createHandler, detailHandler: detailHandler, updateHandler: updateHandler, adminsHandler: adminsHandler}
}

func (this *PublisherGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.listHandler.Handle)
	group.POST("/create", this.createHandler.Handle, permissions.AdminPermission)
	group.GET("/detail/:id", this.detailHandler.Handle)
	group.PUT("/update/:id", this.updateHandler.Handle, permissions.AdminPermission)
	group.PUT("/update/:id", this.updateHandler.Handle, permissions.AdminPermission)
	group.GET("/:publisherId/admins", this.adminsHandler.Handle, permissions.AdminPermission)
}
