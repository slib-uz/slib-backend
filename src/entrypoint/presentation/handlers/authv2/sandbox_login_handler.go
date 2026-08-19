package authv2

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authv2usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/authv2/schema"
)

type SandboxLoginHandler struct {
	uc *authv2usecases.SandboxLoginUseCase
}

// @inject
func NewSandboxLoginHandler(uc *authv2usecases.SandboxLoginUseCase) *SandboxLoginHandler {
	return &SandboxLoginHandler{uc: uc}
}

// Handle godoc
// @Tags authv2
// @Accept  json
// @Produce  json
// @Param data body schema.SandboxLoginRequest true "Login Request"
// @Success 200 {object} schema.LoginResponse
// @Router /auth-v2/sandbox/login [post]
func (this *SandboxLoginHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	req, err := context.GetBody[schema.SandboxLoginRequest](ctx)
	if err != nil {
		return err
	}

	t, err := this.uc.Execute(req.PhoneNumber, req.Otp)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, schema.LoginResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	})
}
