package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/journal_config"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type JournalConfigGroup struct {
	handler                *journal_config.JournalConfigHandler
	createHandler          *journal_config.JournalConfigCreateHandler
	listHandler            *journal_config.JournalConfigListHandler
	isAllowedDomainHandler *journal_config.IsAllowedDomainHandler
}

// @inject
func NewJournalConfigGroup(
	handler *journal_config.JournalConfigHandler,
	createHandler *journal_config.JournalConfigCreateHandler,
	listHandler *journal_config.JournalConfigListHandler,
	isAllowedDomainHandler *journal_config.IsAllowedDomainHandler) *JournalConfigGroup {
	return &JournalConfigGroup{
		handler:                handler,
		createHandler:          createHandler,
		listHandler:            listHandler,
		isAllowedDomainHandler: isAllowedDomainHandler,
	}
}

func (this *JournalConfigGroup) RegisterRoutes(group *echo.Group) {
	group.GET("", this.handler.Handle)
	group.GET("/list", this.listHandler.Handle, permissions.AdminPermission)
	group.POST("/create-or-update", this.createHandler.Handle, permissions.AuthenticatedPermission)
	group.GET("/is-allowed", this.isAllowedDomainHandler.Handle)
}
