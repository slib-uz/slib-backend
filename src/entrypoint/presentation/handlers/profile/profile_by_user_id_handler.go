package profile

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ProfileByIDHandler struct {
	uc *userusecases.UserProfileUseCase
}

// @inject
func NewProfileByIDHandler(uc *userusecases.UserProfileUseCase) *ProfileByIDHandler {
	return &ProfileByIDHandler{uc: uc}
}

// Handle ProfileByIDHandler
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.UserProfileEntity
// @Param id path string true "ID"
// @Security BearerAuth
// @Failure 404 {object} response.Response
// @Router /user/profile/{id} [get]
func (this *ProfileByIDHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JsonResponse(400, "Invalid ID")
	}

	data, err := this.uc.Execute(c.User, uint(userID))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, data)
}
