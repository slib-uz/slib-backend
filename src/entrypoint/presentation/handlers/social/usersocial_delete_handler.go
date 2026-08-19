package social

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/socialusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UserSocialDeleteHandler struct {
	usecase *socialusecases.UserSocialDestroyUseCase
}

// @inject
func NewUserSocialDeleteHandler(usecase *socialusecases.UserSocialDestroyUseCase) *UserSocialDeleteHandler {
	return &UserSocialDeleteHandler{
		usecase: usecase,
	}
}

// Handle UserSocialDeleteHandler
// @Tags Social
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /social/delete/{id} [delete]
func (this *UserSocialDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JsonResponse(http.StatusBadRequest, map[string]string{"error": "invalid ID format"})
	}
	userID := c.User.ID
	err = this.usecase.Execute(uint(id), userID)
	if err != nil {
		return c.JsonResponse(http.StatusBadRequest, err.Error())
	}

	return c.JsonResponse(http.StatusNoContent, nil)
}
