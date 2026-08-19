package institution

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/institutionusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type InstitutionDeleteHandler struct {
	uc *institutionusecases.InstitutionDeleteUseCase
}

// @inject
func NewInstitutionDeleteHandler(uc *institutionusecases.InstitutionDeleteUseCase) *InstitutionDeleteHandler {
	return &InstitutionDeleteHandler{uc: uc}
}

// Handle
// @Tags Institution
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Institution ID"
// @Success 200 {string} string "Success"
// @Failure 403 {string} string "Forbidden"
// @Router /institution/delete/{id} [delete]
func (this *InstitutionDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(uint(id)); err != nil {
		return err
	}

	return c.JsonResponse(200, "Success")
}
