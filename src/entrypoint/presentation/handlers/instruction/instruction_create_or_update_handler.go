package instruction

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/instructionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type InstructionCreateOrUpdateHandler struct {
	uc *instructionusecases.InstructionCreateOrUpdateUseCase
}

// @inject
func NewInstructionCreateOrUpdateHandler(uc *instructionusecases.InstructionCreateOrUpdateUseCase) *InstructionCreateOrUpdateHandler {
	return &InstructionCreateOrUpdateHandler{uc: uc}
}

// Handle godoc
// @Tags         instruction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  entity.InstructionEntity  true  "Instruction data"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /instruction/save [post]
func (this *InstructionCreateOrUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[entity.InstructionEntity](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(data); err != nil {
		return err
	}

	return c.JsonResponse(200, "Instruction saved successfully")
}
