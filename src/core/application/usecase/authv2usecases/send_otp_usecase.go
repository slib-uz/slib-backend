package authv2usecases

import (
	"context"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
)

type SendOtpUseCase struct {
	smsGateway gateway.SmsGateway
	service    *service.OTPService
	throttle   *service.AuthThrottle
}

// @inject
func NewSendOtpUseCase(smsGateway gateway.SmsGateway, service *service.OTPService, throttle *service.AuthThrottle) *SendOtpUseCase {
	return &SendOtpUseCase{smsGateway: smsGateway, service: service, throttle: throttle}
}

func (this *SendOtpUseCase) Execute(ctx context.Context, phoneNumber string, purpose enum.OTPPurpose) (string, error) {
	if this.throttle.CheckAndHitOTPSend(ctx, phoneNumber) {
		return "", response.TooManyRequestsError
	}

	otp, err := this.service.Make(ctx, phoneNumber, purpose)
	if err != nil {
		return "", err
	}

	if err := this.smsGateway.Send(phoneNumber, this.message(otp.Code)); err != nil {
		return "", err
	}
	return otp.SessionID, nil
}

func (this *SendOtpUseCase) message(code string) string {
	return "slib.uz tizimiga kirish uchun tasdiqlash kodi: " + code
}
