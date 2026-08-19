package institution

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/institutionusecases"
	"slib.uz/src/core/domain/entity"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type InstitutionCreateHandler struct {
	uc *institutionusecases.InstitutionCreateUseCase
}

// @inject
func NewInstitutionCreateHandler(uc *institutionusecases.InstitutionCreateUseCase) *InstitutionCreateHandler {
	return &InstitutionCreateHandler{uc: uc}
}

// Handle
// @Tags Institution
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body entity.InstitutionEntity true "Institution Body"
// @Success 200 {object} entity.InstitutionEntity
// @Failure 403 {string} string "Forbidden"
// @Router /institution/create [post]
func (this *InstitutionCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	body, err := context2.GetBody[entity.InstitutionEntity](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(body); err != nil {
		return err
	}

	return c.JsonResponse(200, body)
}
