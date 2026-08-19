package institution

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/institutionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type InstitutionDetailHandler struct {
	uc *InstitutionDetailUseCase
}

// @inject
func NewInstitutionDetailHandler(uc *InstitutionDetailUseCase) *InstitutionDetailHandler {
	return &InstitutionDetailHandler{uc: uc}
}

// Handle
// @Tags Institution
// @Accept json
// @Produce json
// @Param id path int true "Institution ID"
// @Success 200 {object} entity.InstitutionEntity
// @Failure 404 {object} response.Response
// @Router /institution/detail/{id} [get]
func (this *InstitutionDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	institution, err := this.uc.Execute(uint(id))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, institution)
}
