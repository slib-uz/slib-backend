package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/authusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/auth/schema"
)

type UserLoginHandler struct {
	usecase *authusecases.UserLoginUseCase
}

// @inject
func NewUserLoginHandler(usecase *authusecases.UserLoginUseCase) *UserLoginHandler {
	return &UserLoginHandler{usecase: usecase}
}

// Handle UserLoginHandler
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body schema.UserLoginRequest true "User Login Request"
// @Router /auth/oauth/login [post]
func (this *UserLoginHandler) Handle(c echo.Context) error {

	data, err := context.GetBody[schema.UserLoginRequest](c)
	if err != nil {
		return err
	}

	token, err := this.usecase.Execute(data.Code)

	if err != nil {
		return response.NewFailResponse(400, err.Error())
	}

	return c.JSON(http.StatusOK, token)
}
