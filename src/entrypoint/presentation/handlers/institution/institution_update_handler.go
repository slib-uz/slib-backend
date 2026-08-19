package institution

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/institutionusecases"
	"slib.uz/src/core/domain/entity"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type InstitutionUpdateHandler struct {
	uc *institutionusecases.InstitutionUpdateUseCase
}

// @inject
func NewInstitutionUpdateHandler(uc *institutionusecases.InstitutionUpdateUseCase) *InstitutionUpdateHandler {
	return &InstitutionUpdateHandler{uc: uc}
}

// Handle
// @Tags Institution
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Institution ID"
// @Param body body entity.InstitutionEntity true "Institution Body"
// @Success 200 {object} entity.InstitutionEntity
// @Failure 403 {string} string "Forbidden"
// @Router /institution/update/{id} [put]
func (this *InstitutionUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	body, err := context2.GetBody[entity.InstitutionEntity](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(uint(id), body); err != nil {
		return err
	}

	return c.JsonResponse(200, body)
}
