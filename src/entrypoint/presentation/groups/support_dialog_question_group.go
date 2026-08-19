package groups

import (
	"github.com/labstack/echo/v4"
	support_dialog2 "slib.uz/src/entrypoint/presentation/handlers/support_dialog"
)

type SupportDialogQuestionGroup struct {
	supportDialogQuestionListHandler   *support_dialog2.SupportDialogQuestionListHandler
	supportDialogCreateQuestionHandler *support_dialog2.SupportDialogCreateQuestionHandler
}

// @inject
func NewSupportDialogQuestionGroup(supportDialogQuestionListHandler *support_dialog2.SupportDialogQuestionListHandler, supportDialogCreateQuestionHandler *support_dialog2.SupportDialogCreateQuestionHandler) *SupportDialogQuestionGroup {
	return &SupportDialogQuestionGroup{supportDialogQuestionListHandler: supportDialogQuestionListHandler, supportDialogCreateQuestionHandler: supportDialogCreateQuestionHandler}
}

func (this SupportDialogQuestionGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.supportDialogQuestionListHandler.Handle)
	group.POST("/create", this.supportDialogCreateQuestionHandler.Handle)
}
