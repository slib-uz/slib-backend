package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/admin"
)

//var AdminGroup = app.GetInstance().Group(
//	"/admin",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AdminPermission,
//)
//
//func init() {
//	AdminGroup.GET("/antiplag/balance", container.InitAntiPlagBalanceHandler().Handle)
//}

type AdminGroup struct {
	balanceHandler *admin.AntiPlagBalanceHandler
}

// @inject
func NewAdminGroup(balanceHandler *admin.AntiPlagBalanceHandler) *AdminGroup {
	return &AdminGroup{balanceHandler: balanceHandler}
}

func (this *AdminGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/antiplag/balance", this.balanceHandler.Handle)
}
