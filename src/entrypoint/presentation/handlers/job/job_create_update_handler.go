package job

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/mapper"
	. "slib.uz/src/core/application/usecase/jobusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/job/schema"
)

type JobCreateUpdateHandler struct {
	uc *JobCreateUpdateUseCase
}

// @inject
func NewJobCreateUpdateHandler(uc *JobCreateUpdateUseCase) *JobCreateUpdateHandler {
	return &JobCreateUpdateHandler{uc: uc}
}

// Handle
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body schema.JobCreateUpdateRequest true "Job Create/Update Request"
// @Success 200 {object} schema.JobCreateUpdateResponse
// @Router /user/job [post]
func (this *JobCreateUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	request, err := context2.GetBody[schema.JobCreateUpdateRequest](ctx)
	if err != nil {
		return err
	}

	createUpdateDTO := mapper.JobCreateUpdateRequestToDTO(request)

	resultDTO, err := this.uc.Execute(c.User.ID, createUpdateDTO)
	if err != nil {
		return c.JsonResponse(400, err.Error())
	}

	response := mapper.JobDTOToCreateUpdateResponse(resultDTO)

	return c.JsonResponse(200, response)
}
