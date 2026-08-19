package comment

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/commentusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/comment/schema"
)

type CommentCreateHandler struct {
	uc *CommentCreateUseCase
}

// @inject
func NewCommentCreateHandler(uc *CommentCreateUseCase) *CommentCreateHandler {
	return &CommentCreateHandler{uc: uc}
}

// Handle godoc
// @Tags         comment
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        comment  body      schema.CommentCreateRequest  true  "Comment Create Request"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /comment/create [post]
func (this *CommentCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.CommentCreateRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, data.ToEntity()); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return c.JsonResponse(200, "Comment created successfully")

}
