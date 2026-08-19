package deadline

import (
	"net/http"

	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/deadlineusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/deadline/schema"
)

type ExtendDeadlineHandler struct {
	uc *ExtendDeadlineUseCase
}

// @inject
func NewExtendDeadlineHandler(uc *ExtendDeadlineUseCase) *ExtendDeadlineHandler {
	return &ExtendDeadlineHandler{uc: uc}
}

// Handle godoc
// @Tags         deadline
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param data body schema.ExtendSchemaRequest true "Extend deadline request"
// @Success      200 {object} response.Response "Deadline extended successfully"
// @Router 	/deadline/extend [post]
func (this *ExtendDeadlineHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.ExtendSchemaRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, data.ReviewStageID, data.Deadline, data.DeadlineType); err != nil {
		return err
	}

	return c.JsonResponse(http.StatusOK, "Deadline extended successfully")
}
