package instruction

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/instructionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type InstructionListHandler struct {
	uc *instructionusecases.InstructionListUseCase
}

// @inject
func NewInstructionListHandler(uc *instructionusecases.InstructionListUseCase) *InstructionListHandler {
	return &InstructionListHandler{uc: uc}
}

// Handle godoc
// @Tags         instruction
// @Accept       json
// @Produce      json
// @Success      200  {array}  entity.InstructionEntity
// @Failure      400  {object}  response.Response
// @Router       /instruction/list [get]
func (this *InstructionListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	result, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
