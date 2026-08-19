package groups

import (
	"github.com/labstack/echo/v4"
	myapplication2 "slib.uz/src/entrypoint/presentation/handlers/myapplication"
)

//
//var MyApplicationGroup = app.GetInstance().Group(
//	"/my-application",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AuthenticatedPermission,
//)

//func init() {
//	MyApplicationGroup.GET("/list", container.InitMyApplicationsListHandler().Handle)
//	MyApplicationGroup.GET("/detail/:id", container.InitMyApplicationDetailHandler().Handle)
//	MyApplicationGroup.POST("/resubmit-for-technical-review/:applicationId", container.InitApplicationResubmitHandler().Handle)
//	MyApplicationGroup.POST("/resubmit-for-peer-review", container.InitResubmitForPeerReviewHandler().Handle)
//}

type MyApplicationGroup struct {
	list                       *myapplication2.MyApplicationsListHandler
	detail                     *myapplication2.MyApplicationDetailHandler
	resubmitForTechnicalReview *myapplication2.ResubmitForTechnicalReviewHandler
	resubmitForPeerReview      *myapplication2.ResubmitForPeerReviewHandler
}

// @inject
func NewMyApplicationGroup(list *myapplication2.MyApplicationsListHandler, detail *myapplication2.MyApplicationDetailHandler, resubmitForTechnicalReview *myapplication2.ResubmitForTechnicalReviewHandler, resubmitForPeerReview *myapplication2.ResubmitForPeerReviewHandler) *MyApplicationGroup {
	return &MyApplicationGroup{list: list, detail: detail, resubmitForTechnicalReview: resubmitForTechnicalReview, resubmitForPeerReview: resubmitForPeerReview}
}

func (this *MyApplicationGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.list.Handle)
	group.GET("/detail/:id", this.detail.Handle)
	group.POST("/resubmit-for-technical-review/:applicationId", this.resubmitForTechnicalReview.Handle)
	group.POST("/resubmit-for-peer-review", this.resubmitForPeerReview.Handle)
}
