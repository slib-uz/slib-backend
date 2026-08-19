package groups

import (
	"github.com/labstack/echo/v4"
	etaqriz2 "slib.uz/src/entrypoint/presentation/handlers/etaqriz"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

//import (
//	container "slib.uz/src/di"
//	"slib.uz/src/presentation/app"
//	"slib.uz/src/presentation/handlers/etaqriz"
//)
//
//var EtaqrizGroup = app.GetInstance().Group("/etaqriz")
//
//func init() {
//	EtaqrizGroup.POST("/callback", container.InitPeerReviewAnswerHandler().Handle)
//	EtaqrizGroup.GET("/reviewer/find", container.InitFindReviewerHandler().Handle)
//}

type EtaqrizGroup struct {
	peerReviewAnswerHandler *etaqriz2.PeerReviewAnswerHandler
	findReviewerHandler     *etaqriz2.FindReviewerHandler
}

// @inject
func NewEtaqrizGroup(peerReviewAnswerHandler *etaqriz2.PeerReviewAnswerHandler, findReviewerHandler *etaqriz2.FindReviewerHandler) *EtaqrizGroup {
	return &EtaqrizGroup{peerReviewAnswerHandler: peerReviewAnswerHandler, findReviewerHandler: findReviewerHandler}
}

func (this *EtaqrizGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/callback", this.peerReviewAnswerHandler.Handle)
	group.GET("/reviewer/find", this.findReviewerHandler.Handle, permissions.AuthenticatedPermission)
}
