package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/journal"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type JournalGroup struct {
	list              *journal.JournalListHandler
	detail            *journal.JournalDetailHandler
	reviewersList     *journal.JournalReviewersListHandler
	topList           *journal.JournalTopListHandler
	statisticsList    *journal.JournalStatisticsListHandler
	statisticsExcel   *journal.JournalStatisticsExcelHandler
	completionPercent *journal.JournalCompletionPercentHandler

	// middlewares
	anonymAuthMiddleware *middlewares.JwAnonymAuthMiddleware
}

// @inject
func NewJournalGroup(
	list *journal.JournalListHandler,
	detail *journal.JournalDetailHandler,
	reviewersList *journal.JournalReviewersListHandler,
	topList *journal.JournalTopListHandler,
	statisticsList *journal.JournalStatisticsListHandler,
	statisticsExcel *journal.JournalStatisticsExcelHandler,
	anonymAuthMiddleware *middlewares.JwAnonymAuthMiddleware,
	completionPercent *journal.JournalCompletionPercentHandler,
) *JournalGroup {
	return &JournalGroup{
		list:                 list,
		detail:               detail,
		reviewersList:        reviewersList,
		topList:              topList,
		statisticsList:       statisticsList,
		statisticsExcel:      statisticsExcel,
		anonymAuthMiddleware: anonymAuthMiddleware,
		completionPercent:    completionPercent,
	}
}

func (this *JournalGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.list.Handle)
	group.GET("/detail/:id", this.detail.Handle, this.anonymAuthMiddleware.Call)
	group.GET("/:journalId/reviewers", this.reviewersList.Handle, permissions.AuthenticatedPermission)
	group.GET("/top", this.topList.Handle)
	group.GET("/statistics", this.statisticsList.Handle)
	group.GET("/statistics/excel", this.statisticsExcel.Handle)
	group.GET("/:id/completion-percent", this.completionPercent.Handle)
}
