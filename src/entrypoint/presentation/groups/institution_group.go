package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/domain/entity/enum"
	institution2 "slib.uz/src/entrypoint/presentation/handlers/institution"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type InstitutionGroup struct {
	listHandler          *institution2.InstitutionListHandler
	detailHandler        *institution2.InstitutionDetailHandler
	createHandler        *institution2.InstitutionCreateHandler
	updateHandler        *institution2.InstitutionUpdateHandler
	deleteHandler        *institution2.InstitutionDeleteHandler
	setPublishersHandler    *institution2.InstitutionSetPublishersHandler
	detachPublisherHandler  *institution2.InstitutionDetachPublisherHandler
	adminsHandler           *institution2.InstitutionAdminHandler
}

// @inject
func NewInstitutionGroup(
	listHandler *institution2.InstitutionListHandler,
	detailHandler *institution2.InstitutionDetailHandler,
	createHandler *institution2.InstitutionCreateHandler,
	updateHandler *institution2.InstitutionUpdateHandler,
	deleteHandler *institution2.InstitutionDeleteHandler,
	setPublishersHandler *institution2.InstitutionSetPublishersHandler,
	detachPublisherHandler *institution2.InstitutionDetachPublisherHandler,
	adminsHandler *institution2.InstitutionAdminHandler,
) *InstitutionGroup {
	return &InstitutionGroup{
		listHandler:            listHandler,
		detailHandler:          detailHandler,
		createHandler:          createHandler,
		updateHandler:          updateHandler,
		deleteHandler:          deleteHandler,
		setPublishersHandler:   setPublishersHandler,
		detachPublisherHandler: detachPublisherHandler,
		adminsHandler:          adminsHandler,
	}
}

func (this *InstitutionGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.listHandler.Handle)
	group.GET("/detail/:id", this.detailHandler.Handle)
	group.POST("/create", permissions.RolePermission(this.createHandler.Handle, enum.RoleAdmin))
	group.PUT("/update/:id", permissions.RolePermission(this.updateHandler.Handle, enum.RoleAdmin))
	group.PUT("/set-publishers/:id", permissions.RolePermission(this.setPublishersHandler.Handle, enum.RoleAdmin))
	group.POST("/:id/detach-publisher", permissions.RolePermission(this.detachPublisherHandler.Handle, enum.RoleAdmin))
	group.DELETE("/delete/:id", permissions.RolePermission(this.deleteHandler.Handle, enum.RoleAdmin))
	group.GET("/:institutionId/admins", this.adminsHandler.Handle, permissions.AdminPermission)
}
