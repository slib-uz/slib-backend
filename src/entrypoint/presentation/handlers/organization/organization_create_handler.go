package organization

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/organizationusecases"
	"slib.uz/src/core/domain/entity"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type OrganizationCreateHandler struct {
	uc *organizationusecases.OrganizationCreateUseCase
}

// @inject
func NewOrganizationCreateHandler(uc *organizationusecases.OrganizationCreateUseCase) *OrganizationCreateHandler {
	return &OrganizationCreateHandler{uc: uc}
}

// Handle
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body entity.OrganizationEntity true "Organization Body"
// @Success 200 {object} entity.OrganizationEntity
// @Failure 400,403,409 {object} response.Response
// @Router /organization/create [post]
func (this *OrganizationCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	body, err := context2.GetBody[entity.OrganizationEntity](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, body); err != nil {
		return err
	}

	return c.JsonResponse(200, body)
}
