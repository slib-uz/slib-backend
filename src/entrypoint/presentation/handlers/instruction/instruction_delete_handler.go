package instruction

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/instructionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type InstructionDeleteHandler struct {
	uc *instructionusecases.InstructionDeleteUseCase
}

// @inject
func NewInstructionDeleteHandler(uc *instructionusecases.InstructionDeleteUseCase) *InstructionDeleteHandler {
	return &InstructionDeleteHandler{uc: uc}
}

// Handle godoc
// @Tags         instruction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Instruction ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /instruction/{id}/delete [delete]
func (this *InstructionDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	idValue, err := context.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}
	id := uint(idValue)

	if err := this.uc.Execute(id); err != nil {
		return err
	}

	return c.JsonResponse(200, "Instruction deleted successfully")
}
