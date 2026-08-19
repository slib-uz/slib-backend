package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/entrypoint/presentation/handlers/authv2"
)

type AuthV2Group struct {
	sendOtpHandler          *authv2.SendOtpHandler
	verifyAndLoginHandler   *authv2.VerifyAndLoginHandler
	checkPhoneNumberHandler *authv2.CheckPhoneNumberHandler
	sandboxLoginHandler     *authv2.SandboxLoginHandler
	cfg                     conf.ConfigAdapter
}

// @inject
func NewAuthV2Group(
	sendOtpHandler *authv2.SendOtpHandler,
	verifyAndLoginHandler *authv2.VerifyAndLoginHandler,
	checkPhoneNumberHandler *authv2.CheckPhoneNumberHandler,
	sandboxLoginHandler *authv2.SandboxLoginHandler,
	cfg conf.ConfigAdapter,
) *AuthV2Group {
	return &AuthV2Group{
		sendOtpHandler:          sendOtpHandler,
		verifyAndLoginHandler:   verifyAndLoginHandler,
		checkPhoneNumberHandler: checkPhoneNumberHandler,
		sandboxLoginHandler:     sandboxLoginHandler,
		cfg:                     cfg,
	}
}

func (this *AuthV2Group) RegisterRoutes(group *echo.Group) {
	group.POST("/send-otp", this.sendOtpHandler.Handle)
	group.POST("/login", this.verifyAndLoginHandler.Handle)
	group.GET("/check-phone-number", this.checkPhoneNumberHandler.Handle)

	// Sandbox login autentifikatsiyasiz, statik OTP bilan haqiqiy sessiya beradi
	// va foydalanuvchi topilmasa yangi hisob yaratadi.
	// Prod'da route umuman ro'yxatdan o'tmaydi — 404 qaytadi.
	if !this.cfg.IsProduction() {
		group.POST("/sandbox/login", this.sandboxLoginHandler.Handle)
	}
}
