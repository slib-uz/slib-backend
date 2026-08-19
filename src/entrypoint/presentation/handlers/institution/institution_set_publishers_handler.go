package institution

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/institutionusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/institution/schema"
)

type InstitutionSetPublishersHandler struct {
	uc *institutionusecases.InstitutionSetPublishersUseCase
}

// @inject
func NewInstitutionSetPublishersHandler(uc *institutionusecases.InstitutionSetPublishersUseCase) *InstitutionSetPublishersHandler {
	return &InstitutionSetPublishersHandler{uc: uc}
}

// Handle
// @Tags Institution
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Institution ID"
// @Param body body schema.InstitutionSetPublishersRequest true "Publisher IDs"
// @Success 200 {string} string "Success"
// @Failure 403 {string} string "Forbidden"
// @Router /institution/set-publishers/{id} [put]
func (this *InstitutionSetPublishersHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	body, err := context2.GetBody[schema.InstitutionSetPublishersRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(uint(id), body.PublisherIDs); err != nil {
		return err
	}

	return c.JsonResponse(200, "Success")
}
