package social

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/socialusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/social/schema"
	"slib.uz/src/infrastructure/persistence/mapper"
)

type UserSocialUpdateHandler struct {
	usecase *socialusecases.UserSocialUpdateUseCase
}

// @inject
func NewUserSocialUpdateHandler(usecase *socialusecases.UserSocialUpdateUseCase) *UserSocialUpdateHandler {
	return &UserSocialUpdateHandler{
		usecase: usecase,
	}
}

// Handle Handler UserSocialHandler
// @Summary Update User Social
// @Description Update User Social
// @Tags Social
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param data body schema.UserSocialRequest true "User Social Request"
// @Security BearerAuth
// @Router /social/update/{id} [put]
func (this *UserSocialUpdateHandler) Handle(c echo.Context) error {
	ctx := c.(*context2.Context)
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ID format")
	}

	reqBody, err := context2.GetBody[schema.UserSocialRequest](ctx)
	if err != nil {
		return err
	}

	userID := ctx.User.ID

	dto := mapper.UserSocialRequestToDTO(reqBody)
	err = this.usecase.Execute(uint(id), dto, userID)
	if err != nil {
		return ctx.JsonResponse(http.StatusBadRequest, err.Error())
	}

	return ctx.JsonResponse(http.StatusOK, nil)
}
