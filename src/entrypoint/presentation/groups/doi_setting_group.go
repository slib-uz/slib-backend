package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/doi"
)

type DoiSettingGroup struct {
	getHandler         *doi.DoiSettingGetHandler
	createHandler      *doi.DoiSettingCreateHandler
	updateHandler      *doi.DoiSettingUpdateHandler
	depositHandler     *doi.CrossRefDepositHandler
	depositListHandler *doi.CrossRefDepositListHandler
	checkAuthHandler   *doi.CrossRefCheckAuthHandler
	prefixInfoHandler  *doi.CrossRefPrefixInfoHandler
}

// @inject
func NewDoiSettingGroup(
	getHandler *doi.DoiSettingGetHandler,
	createHandler *doi.DoiSettingCreateHandler,
	updateHandler *doi.DoiSettingUpdateHandler,
	depositHandler *doi.CrossRefDepositHandler,
	depositListHandler *doi.CrossRefDepositListHandler,
	checkAuthHandler *doi.CrossRefCheckAuthHandler,
	prefixInfoHandler *doi.CrossRefPrefixInfoHandler,
) *DoiSettingGroup {
	return &DoiSettingGroup{
		getHandler:         getHandler,
		createHandler:      createHandler,
		updateHandler:      updateHandler,
		depositHandler:     depositHandler,
		depositListHandler: depositListHandler,
		checkAuthHandler:   checkAuthHandler,
		prefixInfoHandler:  prefixInfoHandler,
	}
}

func (this *DoiSettingGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/deposits", this.depositListHandler.Handle)
	group.POST("/deposit/:articleId", this.depositHandler.Handle)
	group.POST("/check-auth", this.checkAuthHandler.Handle)
	group.GET("/prefix/:prefix", this.prefixInfoHandler.Handle)

	group.GET("/:journalId", this.getHandler.Handle)
	group.POST("/:journalId", this.createHandler.Handle)
	group.PUT("/:journalId", this.updateHandler.Handle)
}
