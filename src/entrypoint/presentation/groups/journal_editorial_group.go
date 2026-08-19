package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/entrypoint/presentation/handlers/journal_editorial"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type JournalEditorialGroup struct {
	createHandler      *journal_editorial.JournalEditorialCreateHandler
	updateHandler      *journal_editorial.JournalEditorialUpdateHandler
	deleteHandler      *journal_editorial.JournalEditorialDeleteHandler
	listHandler        *journal_editorial.JournalEditorialListHandler
	listByRoleHandler  *journal_editorial.JournalEditorialListByRoleHandler
	detailHandler      *journal_editorial.JournalEditorialDetailHandler
}

// @inject
func NewJournalEditorialGroup(
	createHandler *journal_editorial.JournalEditorialCreateHandler,
	updateHandler *journal_editorial.JournalEditorialUpdateHandler,
	deleteHandler *journal_editorial.JournalEditorialDeleteHandler,
	listHandler *journal_editorial.JournalEditorialListHandler,
	listByRoleHandler *journal_editorial.JournalEditorialListByRoleHandler,
	detailHandler *journal_editorial.JournalEditorialDetailHandler,
) *JournalEditorialGroup {
	return &JournalEditorialGroup{
		createHandler:     createHandler,
		updateHandler:     updateHandler,
		deleteHandler:     deleteHandler,
		listHandler:       listHandler,
		listByRoleHandler: listByRoleHandler,
		detailHandler:     detailHandler,
	}
}

func (this *JournalEditorialGroup) RegisterRoutes(group *echo.Group) {
	allowed := []enum.UserRole{enum.RoleAdmin, enum.RoleExpert, enum.RolePublisherAdmin}

	group.POST("/:journalId/create", permissions.RolePermission(this.createHandler.Handle, allowed...))
	group.PUT("/:id/update", permissions.RolePermission(this.updateHandler.Handle, allowed...))
	group.DELETE("/:id/delete", permissions.RolePermission(this.deleteHandler.Handle, allowed...))
	group.GET("/:journalId/list", this.listHandler.Handle)
	group.GET("/:journalId/list-by-role", this.listByRoleHandler.Handle)
	group.GET("/:id/detail", this.detailHandler.Handle)
}
