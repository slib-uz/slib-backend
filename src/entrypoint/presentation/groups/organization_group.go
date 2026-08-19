package groups

import (
	"github.com/labstack/echo/v4"
	organization2 "slib.uz/src/entrypoint/presentation/handlers/organization"
)

type OrganizationGroup struct {
	listHandler   *organization2.OrganizationListHandler
	detailHandler *organization2.OrganizationDetailHandler
	createHandler *organization2.OrganizationCreateHandler
	updateHandler *organization2.OrganizationUpdateHandler
	deleteHandler *organization2.OrganizationDeleteHandler
}

// @inject
func NewOrganizationGroup(
	listHandler *organization2.OrganizationListHandler,
	detailHandler *organization2.OrganizationDetailHandler,
	createHandler *organization2.OrganizationCreateHandler,
	updateHandler *organization2.OrganizationUpdateHandler,
	deleteHandler *organization2.OrganizationDeleteHandler,
) *OrganizationGroup {
	return &OrganizationGroup{
		listHandler:   listHandler,
		detailHandler: detailHandler,
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
	}
}

func (this *OrganizationGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.listHandler.Handle)
	group.GET("/detail/:id", this.detailHandler.Handle)
	group.POST("/create", this.createHandler.Handle)
	group.PUT("/update/:id", this.updateHandler.Handle)
	group.DELETE("/delete/:id", this.deleteHandler.Handle)
}
