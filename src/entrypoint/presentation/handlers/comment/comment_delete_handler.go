package comment

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/commentusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type CommentDeleteHandler struct {
	uc *CommentDeleteUseCase
}

// @inject
func NewCommentDeleteHandler(uc *CommentDeleteUseCase) *CommentDeleteHandler {
	return &CommentDeleteHandler{uc: uc}
}

// Handle CommentDeleteHandler
// @Tags comment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /comment/{id} [delete]
func (this *CommentDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}
	if err := this.uc.Execute(c.User.ID, uint(id)); err != nil {
		return c.JsonResponse(400, err.Error())
	}
	return c.JsonResponse(200, "Comment deleted successfully")
}
