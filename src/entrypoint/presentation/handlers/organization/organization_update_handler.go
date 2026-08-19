package organization

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/organizationusecases"
	"slib.uz/src/core/domain/entity"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type OrganizationUpdateHandler struct {
	uc *organizationusecases.OrganizationUpdateUseCase
}

// @inject
func NewOrganizationUpdateHandler(uc *organizationusecases.OrganizationUpdateUseCase) *OrganizationUpdateHandler {
	return &OrganizationUpdateHandler{uc: uc}
}

// Handle
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Param body body entity.OrganizationEntity true "Organization Body"
// @Success 200 {object} entity.OrganizationEntity
// @Failure 403 {string} string "Forbidden"
// @Router /organization/update/{id} [put]
func (this *OrganizationUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	body, err := context2.GetBody[entity.OrganizationEntity](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(uint(id), body, c.User); err != nil {
		return err
	}

	return c.JsonResponse(200, body)
}
