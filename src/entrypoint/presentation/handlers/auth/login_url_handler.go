package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/auth/schema"
)

type LoginURLHandler struct {
	usecase *authusecases.AuthorizeURLUseCase
}

// @inject
func NewLoginURLHandler(usecase *authusecases.AuthorizeURLUseCase) *LoginURLHandler {
	return &LoginURLHandler{usecase: usecase}
}

// Handle LoginURLHandler
// @Accept  json
// @Produce  json
// @Tags auth
// @Param request body schema.LoginURLRequest true "Login URL Request"
// @Success 200 {object} schema.LoginURLResponse
// @Router /auth/oauth/authorize/url [post]
func (this *LoginURLHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.LoginURLRequest](c)
	if err != nil {
		return err
	}

	loginURL := this.usecase.Execute(data.RedirectURL)

	return c.JsonResponse(http.StatusOK, schema.NewLoginURLResponse(data.RedirectURL, loginURL))
}
