package groups

import (
	"github.com/labstack/echo/v4"
	support_dialog2 "slib.uz/src/entrypoint/presentation/handlers/support_dialog"
)

type SupportDialogAnswerGroup struct {
	supportDialogCreateAnswerHandler *support_dialog2.SupportDialogCreateAnswerHandler
	supportDialogAnswerListHandler   *support_dialog2.SupportDialogAnswerListHandler
	chatListHandler                  *support_dialog2.ChatListHandler
}

// @inject
func NewSupportDialogAnswerGroup(supportDialogCreateAnswerHandler *support_dialog2.SupportDialogCreateAnswerHandler, supportDialogAnswerListHandler *support_dialog2.SupportDialogAnswerListHandler, chatListHandler *support_dialog2.ChatListHandler) *SupportDialogAnswerGroup {
	return &SupportDialogAnswerGroup{supportDialogCreateAnswerHandler: supportDialogCreateAnswerHandler, supportDialogAnswerListHandler: supportDialogAnswerListHandler, chatListHandler: chatListHandler}
}

func (this SupportDialogAnswerGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create", this.supportDialogCreateAnswerHandler.Handle)
	group.GET("/list/:chatID", this.supportDialogAnswerListHandler.Handle)
	group.GET("/chat/list", this.chatListHandler.Handle)
}
