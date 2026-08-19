package authv2

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authv2usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/authv2/schema"
)

type VerifyAndLoginHandler struct {
	uc *authv2usecases.VerifyAndLoginUseCase
}

// @inject
func NewVerifyAndLoginHandler(uc *authv2usecases.VerifyAndLoginUseCase) *VerifyAndLoginHandler {
	return &VerifyAndLoginHandler{uc: uc}
}

// Handle godoc
// @Tags authv2
// @Accept  json
// @Produce  json
// @Param data body schema.LoginRequest true "Login Request"
// @Success 200 {object} schema.LoginResponse
// @Router /auth-v2/login [post]
func (this *VerifyAndLoginHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[schema.LoginRequest](ctx)
	if err != nil {
		return err
	}

	token, scope, err := this.uc.Execute(ctx.Request().Context(), data.SessionID, data.Code, data.ScienceID, data.FullName)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, schema.LoginResponse{
		Scope:        scope,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}
