package social

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/socialusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/social/schema"
	"slib.uz/src/infrastructure/persistence/mapper"
)

type UserSocialCreateHandler struct {
	usecase *socialusecases.UserSocialCreateUseCase
}

// @inject
func NewUserSocialCreateHandler(usecase *socialusecases.UserSocialCreateUseCase) *UserSocialCreateHandler {
	return &UserSocialCreateHandler{
		usecase: usecase,
	}
}

// Handle Handler UserSocialCreateHandler
// @Summary Create User Social
// @Description Create User Social
// @Tags Social
// @Accept json
// @Produce json
// @Param data body schema.UserSocialRequest true "User Social Request"
// @Security BearerAuth
// @Router /social/create/ [post]
func (this *UserSocialCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	data, err := context2.GetBody[schema.UserSocialRequest](c)
	if err != nil {
		return err
	}

	userID := c.User.ID
	socialReqToDTO := mapper.UserSocialRequestToDTO(data)

	err = this.usecase.Execute(socialReqToDTO, userID)
	if err != nil {
		return c.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JsonResponse(http.StatusCreated, nil)
}
