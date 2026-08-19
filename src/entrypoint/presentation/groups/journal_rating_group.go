package groups

import (
	"github.com/labstack/echo/v4"
	journal_rating2 "slib.uz/src/entrypoint/presentation/handlers/journal_rating"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type JournalRatingGroup struct {
	createHandler *journal_rating2.JournalRatingCreateHandler
	listHandler   *journal_rating2.JournalRatingListHandler
	statsHandler  *journal_rating2.JournalRatingStatsHandler
	deleteHandler *journal_rating2.JournalRatingDeleteHandler
}

// @inject
func NewJournalRatingGroup(createHandler *journal_rating2.JournalRatingCreateHandler, listHandler *journal_rating2.JournalRatingListHandler, statsHandler *journal_rating2.JournalRatingStatsHandler, deleteHandler *journal_rating2.JournalRatingDeleteHandler) *JournalRatingGroup {
	return &JournalRatingGroup{createHandler: createHandler, listHandler: listHandler, statsHandler: statsHandler, deleteHandler: deleteHandler}
}

func (this *JournalRatingGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list/:journalId", this.listHandler.Handle)
	group.GET("/stats/:journalId", this.statsHandler.Handle)
	authGroup := group.Group("", permissions.AuthenticatedPermission)
	authGroup.POST("/create", this.createHandler.Handle)
	authGroup.DELETE("/:id", this.deleteHandler.Handle, permissions.AuthenticatedPermission)
}
