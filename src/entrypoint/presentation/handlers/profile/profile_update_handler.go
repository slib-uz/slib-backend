package profile

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/mapper"
	. "slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/core/domain/entity"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/profile/schema"
)

type ProfileUpdateHandler struct {
	uc *UserUpdateUseCase
}

// @inject
func NewProfileUpdateHandler(uc *UserUpdateUseCase) *ProfileUpdateHandler {
	return &ProfileUpdateHandler{uc: uc}
}

// Handle
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param data body schema.UserUpdateRequest false "User Update Request (JSON)"
// @Success 200 {object} schema.UserUpdateResponse
// @Router /user/profile/update [patch]
func (this *ProfileUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	var updateDTO *entity.UserUpdateEntity

	request, err := context2.GetBody[schema.UserUpdateRequest](ctx)
	if err != nil {
		return err
	}

	updateDTO = mapper.UserUpdateRequestToEntity(request)

	resultDTO, err := this.uc.Execute(c.User.ID, updateDTO)
	if err != nil {
		return c.JsonResponse(400, err.Error())
	}

	response := mapper.UserUpdateEntityToResponse(resultDTO)

	return c.JsonResponse(200, response)
}
