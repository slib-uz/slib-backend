package support_dialog

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/support_dialog_usecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ChatListHandler struct {
	uc *ChatListUseCase
}

// @inject
func NewChatListHandler(uc *ChatListUseCase) *ChatListHandler {
	return &ChatListHandler{uc: uc}
}

// Handle ChatListHandler
// @Summary List chats
// @Description List chats
// @Security BearerAuth
// @Tags Support Dialog
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param ordering query string false "Ordering" Enums(created_at,-created_at)
// @Success 200 {array} entity.SupportDialogEntity
// @Failure 400 {object} response.Response
// @Router /support-dialog-answer/chat/list [get]
func (this ChatListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	ordering := ctx.QueryParam("ordering")
	paging, err := this.uc.Execute(page, pageSize, ordering)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
