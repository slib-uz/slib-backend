package institution

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/institutionusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/institution/schema"
)

type InstitutionDetachPublisherHandler struct {
	uc *institutionusecases.InstitutionDetachPublisherUseCase
}

// @inject
func NewInstitutionDetachPublisherHandler(uc *institutionusecases.InstitutionDetachPublisherUseCase) *InstitutionDetachPublisherHandler {
	return &InstitutionDetachPublisherHandler{uc: uc}
}

// Handle
// @Tags Institution
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Institution ID"
// @Param body body schema.InstitutionDetachPublisherRequest true "Publisher IDs"
// @Success 200 {object} map[string]int64
// @Failure 403 {string} string "Forbidden"
// @Router /institution/{id}/detach-publisher [post]
func (this *InstitutionDetachPublisherHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	body, err := context2.GetBody[schema.InstitutionDetachPublisherRequest](c)
	if err != nil {
		return err
	}

	affected, err := this.uc.Execute(uint(id), body.PublisherIDs)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]int64{"detached_count": affected})
}
