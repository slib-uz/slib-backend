package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/edition"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type EditionGroup struct {
	createHandler         *edition.EditionCreateHandler
	detailHandler         *edition.EditionDetailHandler
	listHandler           *edition.EditionListHandler
	listByYearHandler     *edition.EditionByYearListHandler
	attachArticlesHandler *edition.EditionAttachArticlesHandler
	articlesListHandler   *edition.EditionArticlesListHandler
	updateHandler         *edition.EditionUpdateHandler
	detachArticlesHandler *edition.EditionDetachArticlesHandler
	deleteHandler         *edition.EditionDeleteHandler
}

// @inject
func NewEditionGroup(
	createHandler *edition.EditionCreateHandler,
	detailHandler *edition.EditionDetailHandler,
	listHandler *edition.EditionListHandler,
	listByYearHandler *edition.EditionByYearListHandler,
	attachArticlesHandler *edition.EditionAttachArticlesHandler,
	articlesListHandler *edition.EditionArticlesListHandler,
	updateHandler *edition.EditionUpdateHandler,
	detachArticlesHandler *edition.EditionDetachArticlesHandler,
	deleteHandler *edition.EditionDeleteHandler,
) *EditionGroup {
	return &EditionGroup{
		createHandler:         createHandler,
		detailHandler:         detailHandler,
		listHandler:           listHandler,
		listByYearHandler:     listByYearHandler,
		attachArticlesHandler: attachArticlesHandler,
		articlesListHandler:   articlesListHandler,
		updateHandler:         updateHandler,
		detachArticlesHandler: detachArticlesHandler,
		deleteHandler:         deleteHandler,
	}
}

func (this *EditionGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create", this.createHandler.Handle, permissions.AuthenticatedPermission)
	group.GET("/:id/detail", this.detailHandler.Handle)
	group.GET("/:journalId/list", this.listHandler.Handle)
	group.GET("/:journalId/list-by-year", this.listByYearHandler.Handle)
	group.POST("/:editionId/attach-articles", this.attachArticlesHandler.Handle, permissions.AuthenticatedPermission)
	group.GET("/articles", this.articlesListHandler.Handle)
	group.PUT("/:editionId/update", this.updateHandler.Handle, permissions.AuthenticatedPermission)
	group.DELETE("/:editionId/delete", this.deleteHandler.Handle, permissions.AuthenticatedPermission)
	group.POST("/:editionId/detach-articles", this.detachArticlesHandler.Handle, permissions.AuthenticatedPermission)
}
