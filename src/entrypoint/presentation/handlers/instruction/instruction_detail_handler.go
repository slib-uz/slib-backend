package instruction

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/instructionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type InstructionDetailHandler struct {
	uc *instructionusecases.InstructionDetailUseCase
}

// @inject
func NewInstructionDetailHandler(uc *instructionusecases.InstructionDetailUseCase) *InstructionDetailHandler {
	return &InstructionDetailHandler{uc: uc}
}

// Handle godoc
// @Tags         instruction
// @Accept       json
// @Produce      json
// @Param        key  path  string  true  "Instruction key"
// @Success      200  {object}  entity.InstructionEntity
// @Failure      404  {object}  response.Response
// @Router       /instruction/{key} [get]
func (this *InstructionDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	key := ctx.Param("key")

	result, err := this.uc.Execute(key)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
