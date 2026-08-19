package authv2

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authv2usecases"
)

type CheckPhoneNumberHandler struct {
	uc *authv2usecases.CheckPhoneNumberUseCase
}

// @inject
func NewCheckPhoneNumberHandler(uc *authv2usecases.CheckPhoneNumberUseCase) *CheckPhoneNumberHandler {
	return &CheckPhoneNumberHandler{uc: uc}
}

// Handle godoc
// @Tags authv2
// @Accept json
// @Produce json
// @Param phone_number query string true "Phone number"
// @Success 204 "No content"
// @Failure 400 {object} response.Response
// @Router /auth-v2/check-phone-number [get]
func (this *CheckPhoneNumberHandler) Handle(ctx echo.Context) error {
	phoneNumber := ctx.QueryParam("phone_number")

	if err := this.uc.Execute(phoneNumber); err != nil {
		return err
	}
	return ctx.NoContent(204)
}
