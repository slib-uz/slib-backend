package job

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/jobusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JobDeleteHandler struct {
	uc *JobDeleteUseCase
}

// @inject
func NewJobDeleteHandler(uc *JobDeleteUseCase) *JobDeleteHandler {
	return &JobDeleteHandler{uc: uc}
}

// Handle
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /user/job [delete]
func (this *JobDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	err := this.uc.Execute(c.User.ID)
	if err != nil {
		return c.JsonResponse(400, err.Error())
	}

	return c.JsonResponse(200, map[string]string{"message": "Job deleted successfully"})
}
