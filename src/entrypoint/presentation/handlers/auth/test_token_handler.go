package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type TestTokenHandler struct {
	uc *authusecases.TestTokenUseCase
}

// @inject
func NewTestTokenHandler(uc *authusecases.TestTokenUseCase) *TestTokenHandler {
	return &TestTokenHandler{uc: uc}
}

// Handle
// @Accept  json
// @Produce  json
// @Tags auth
// @Security BasicAuth
// @Param id path int true "User ID"
// @Success 200 {object} entity.AuthTokenEntity
// @Failure 404 {object} response.Response
// @Router /auth/test/token/{id} [get]
func (this *TestTokenHandler) Handle(ctx echo.Context) error {
	id, err := context.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}

	token, err := this.uc.Execute(uint(id))

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, token)
}
