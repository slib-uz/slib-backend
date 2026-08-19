package institution

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	. "slib.uz/src/core/application/usecase/institutionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type InstitutionAdminHandler struct {
	uc *InstitutionAdminUseCase
}

// @inject
func NewInstitutionAdminHandler(uc *InstitutionAdminUseCase) *InstitutionAdminHandler {
	return &InstitutionAdminHandler{uc: uc}
}

// Handle
// @Tags Institution
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param institutionId path int true "Institution ID"
// @Success 200 {array} entity.UserRoleWithBasicUserEntity
// @Failure 400 {object} response.Response
// @Router /institution/{institutionId}/admins [get]
func (this *InstitutionAdminHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	result, err := this.uc.Execute(this.GetInstitutionID(ctx))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, result)
}

func (this *InstitutionAdminHandler) GetInstitutionID(ctx echo.Context) uint {
	institutionId := ctx.Param("institutionId")
	if institutionId == "" {
		panic(response.NewFailResponse(400, "Institution ID is required"))
	}
	id, err := strconv.Atoi(institutionId)
	if err != nil {
		panic(response.NewFailResponse(400, "Invalid Institution ID"))
	}
	return uint(id)
}
