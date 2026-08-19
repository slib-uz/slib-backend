package profile

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/activityusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type DoctoratesHandler struct {
	uc *UserDoctorateUseCase
}

// @inject
func NewDoctoratesHandler(uc *UserDoctorateUseCase) *DoctoratesHandler {
	return &DoctoratesHandler{uc: uc}
}

// Handle JobsHandler
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} []entity.JobEntity
// @Security BearerAuth
// @Failure 404 {object} response.Response
// @Router /user/doctorate [get]
func (this *DoctoratesHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	degree, err := this.uc.Execute(c.User.ID)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, degree)
}
