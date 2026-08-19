package auth

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/auth/schema"
)

type LogoutHandler struct {
	uc *authusecases.LogoutUseCase
}

// @inject
func NewLogoutHandler(uc *authusecases.LogoutUseCase) *LogoutHandler {
	return &LogoutHandler{uc: uc}
}

// Handle godoc
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data  body      schema.LogoutRequest  false  "Refresh token (ixtiyoriy)"
// @Success      200   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      503   {object}  response.Response
// @Router       /auth/logout [post]
func (this *LogoutHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	// Body bo'sh bo'lishi mumkin — bu xato emas, faqat access bekor qilinadi.
	req, err := context2.GetBody[schema.LogoutRequest](ctx)
	if err != nil {
		req = &schema.LogoutRequest{}
	}

	if err := this.uc.Execute(
		c.Request().Context(),
		c.User.ID,
		c.TokenID,
		c.TokenExp,
		req.RefreshToken,
	); err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]any{"success": true})
}
