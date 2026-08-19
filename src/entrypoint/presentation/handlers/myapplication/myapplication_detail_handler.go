package myapplication

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/myapplicationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type MyApplicationDetailHandler struct {
	uc *MyApplicationDetailUseCase
}

// @inject
func NewMyApplicationDetailHandler(uc *MyApplicationDetailUseCase) *MyApplicationDetailHandler {
	return &MyApplicationDetailHandler{uc: uc}
}

// Handle
// @Tags my-application
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} entity.ApplicationResponseEntity
// @Failure 404 {object} response.Response "Application not found"
// @Router /my-application/detail/{id} [get]
func (this *MyApplicationDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	application, err := this.uc.Execute(c.User.ID, uint(id))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, application)
}
