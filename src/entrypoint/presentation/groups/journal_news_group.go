package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/journal_news"
)

type JournalNewsGroup struct {
	createHandler *journal_news.JournalNewsCreateHandler
	updateHandler *journal_news.JournalNewsUpdateHandler
	deleteHandler *journal_news.JournalNewsDeleteHandler
	listHandler   *journal_news.JournalNewsListHandler
	detailHandler *journal_news.JournalNewsDetailHandler
}

// @inject
func NewJournalNewsGroup(
	createHandler *journal_news.JournalNewsCreateHandler,
	updateHandler *journal_news.JournalNewsUpdateHandler,
	deleteHandler *journal_news.JournalNewsDeleteHandler,
	listHandler *journal_news.JournalNewsListHandler,
	detailHandler *journal_news.JournalNewsDetailHandler,
) *JournalNewsGroup {
	return &JournalNewsGroup{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		listHandler:   listHandler,
		detailHandler: detailHandler,
	}
}

func (this *JournalNewsGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/:journalId/create", this.createHandler.Handle)
	group.PUT("/:id/update", this.updateHandler.Handle)
	group.DELETE("/:id/delete", this.deleteHandler.Handle)
	group.GET("/:journalId/list", this.listHandler.Handle)
	group.GET("/:id/detail", this.detailHandler.Handle)
}
