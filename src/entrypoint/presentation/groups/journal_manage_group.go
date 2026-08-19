package groups

import (
	"github.com/labstack/echo/v4"
	journal2 "slib.uz/src/entrypoint/presentation/handlers/journal"
)

//func init() {
//	JournalGroup.GET("/list", container.InitJournalListHandler().Handle)
//	JournalGroup.GET("/detail/:id", container.InitJournalDetailHandler().Handle)
//	JournalGroup.POST("/reviewers/add", container.InitAddReviewersHandler().Handle)
//	JournalGroup.GET("/:journalId/reviewers", container.InitJournalReviewersListHandler().Handle)
//	JournalGroup.POST("/reviewers/remove", container.InitRemoveReviewersHandler().Handle)
//	JournalGroup.GET("/:journalId/members", container.InitJournalMembersListHandler().Handle)
//	JournalGroup.PATCH("/update/:journalId", container.InitJournalUpdateHandler().Handle)
//	JournalGroup.GET("/:journalId/requirements", container.InitJournalRequirementsListHandler().Handle)
//}

type JournalManageGroup struct {
	addReviewers    *journal2.AddReviewersHandler
	removeReviewers *journal2.RemoveReviewersHandler
	membersList     *journal2.JournalMembersListHandler
	membersLastOnlineList *journal2.JournalMembersLastOnlineListHandler
	update          *journal2.JournalUpdateHandler
	status          *journal2.JournalStatusHandler
}

// @inject
func NewJournalManageGroup(addReviewers *journal2.AddReviewersHandler, removeReviewers *journal2.RemoveReviewersHandler, membersList *journal2.JournalMembersListHandler, membersLastOnlineList *journal2.JournalMembersLastOnlineListHandler, update *journal2.JournalUpdateHandler, status *journal2.JournalStatusHandler) *JournalManageGroup {
	return &JournalManageGroup{addReviewers: addReviewers, removeReviewers: removeReviewers, membersList: membersList, membersLastOnlineList: membersLastOnlineList, update: update, status: status}
}

func (this *JournalManageGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/reviewers/add", this.addReviewers.Handle)
	group.POST("/reviewers/remove", this.removeReviewers.Handle)
	group.GET("/:journalId/members", this.membersList.Handle)
	group.GET("/:journalId/members/last-online", this.membersLastOnlineList.Handle)
	group.PUT("/update/:journalId", this.update.Handle)
	group.PATCH("/:id/status", this.status.Handle)
}
