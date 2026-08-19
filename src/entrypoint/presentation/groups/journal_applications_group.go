package groups

import (
	"github.com/labstack/echo/v4"
	journal_applications2 "slib.uz/src/entrypoint/presentation/handlers/journal_applications"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

//var JournalApplicationsGroup = app.GetInstance().Group(
//	"/journal-applications",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AuthenticatedPermission,
//)
//
//func init() {
//	JournalApplicationsGroup.GET("/applications", container.InitJournalApplicationListHandler().Handle)
//	JournalApplicationsGroup.GET("/applications/:id", container.InitJournalApplicationDetailHandler().Handle)
//	JournalApplicationsGroup.POST("/create", container.InitJournalApplicationCreateHandler().Handle)
//	JournalApplicationsGroup.GET("/check-issn", container.InitCheckISSNHandler().Handle)
//	JournalApplicationsGroup.POST("/review", container.InitReviewApplicationHandler().Handle, permissions.AdminPermission)
//}

type JournalApplicationsGroup struct {
	list      *journal_applications2.JournalApplicationListHandler
	detail    *journal_applications2.JournalApplicationDetailHandler
	create    *journal_applications2.JournalApplicationCreateHandler
	checkISSN *journal_applications2.CheckISSNHandler
	review    *journal_applications2.ReviewApplicationHandler
}

// @inject
func NewJournalApplicationsGroup(list *journal_applications2.JournalApplicationListHandler, detail *journal_applications2.JournalApplicationDetailHandler, create *journal_applications2.JournalApplicationCreateHandler, checkISSN *journal_applications2.CheckISSNHandler, review *journal_applications2.ReviewApplicationHandler) *JournalApplicationsGroup {
	return &JournalApplicationsGroup{list: list, detail: detail, create: create, checkISSN: checkISSN, review: review}
}

func (this *JournalApplicationsGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/applications", this.list.Handle)
	group.GET("/applications/:id", this.detail.Handle)
	group.POST("/create", this.create.Handle)
	group.GET("/check-issn", this.checkISSN.Handle)
	group.POST("/review", this.review.Handle, permissions.AdminPermission)
}
