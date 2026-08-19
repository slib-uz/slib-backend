package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/purchase"
)

//import (
//	container "slib.uz/src/di"
//	"slib.uz/src/presentation/app"
//	"slib.uz/src/presentation/interceptor/permissions"
//)
//
//var PurchaseGroup = app.GetInstance().Group(
//	"/purchase",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AuthenticatedPermission,
//)
//
//func init() {
//	PurchaseGroup.GET("/my-purchased-articles", container.InitMyPurchasedArticlesListHandler().Handle)
//}

type PurchaseGroup struct {
	userPurchasedArticlesListHandler *purchase.UserPurchasedArticlesListHandler
}

// @inject
func NewPurchaseGroup(myPurchasedArticlesListHandler *purchase.UserPurchasedArticlesListHandler) *PurchaseGroup {
	return &PurchaseGroup{userPurchasedArticlesListHandler: myPurchasedArticlesListHandler}
}

func (this *PurchaseGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/my-purchased-articles", this.userPurchasedArticlesListHandler.Handle)
}
