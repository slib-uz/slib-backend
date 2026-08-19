package authv2

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authv2usecases"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/authv2/schema"
)

type SendOtpHandler struct {
	uc *authv2usecases.SendOtpUseCase
}

// @inject
func NewSendOtpHandler(uc *authv2usecases.SendOtpUseCase) *SendOtpHandler {
	return &SendOtpHandler{uc: uc}
}

// Handle godoc
// @Tags authv2
// @Accept json
// @Produce json
// @Param request body schema.SendOtpRequest true "Send OTP Request"
// @Success 200 {object} schema.SendOtpResponse "OTP sent successfully"
// @Failure 400 {object} response.Response
// @Router /auth-v2/send-otp [post]
func (this *SendOtpHandler) Handle(ctx echo.Context) error {
	data, err := context.GetBody[schema.SendOtpRequest](ctx)
	if err != nil {
		return err
	}

	sessionID, err := this.uc.Execute(ctx.Request().Context(), data.Phone, enum.OTPPurposeLogin)
	if err != nil {
		return err
	}

	response := schema.SendOtpResponse{SessionID: sessionID}
	return ctx.JSON(200, map[string]any{
		"status": 200,
		"data":   response,
	})
}
