package authv2usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/usecase/authv2usecases"
	"slib.uz/src/core/domain/entity/enum"
)

// fakeSms — yuborilgan SMS sonini sanaydi.
type fakeSms struct{ sends int }

func (f *fakeSms) Send(phone, message string) error { f.sends++; return nil }

// blockingLimiter — birinchi hit'dan boshlab yuqori sanoq qaytaradi (har doim
// bloklaydi). Send throttle'ni ishga tushirish uchun.
type blockingLimiter struct{}

func (blockingLimiter) Hit(ctx context.Context, key string, window time.Duration) (int64, error) {
	return 999, nil
}
func (blockingLimiter) Reset(ctx context.Context, key string) error { return nil }

func blockingThrottle() *service.AuthThrottle {
	return service.NewAuthThrottle(blockingLimiter{}, fakeConfV2{})
}

// fakeConfV2 — ConfigAdapter throttle qismini beradi (verifyMax=5, sendMin=1).
// Boshqa metodlar panik.
type fakeConfV2 struct{}

func (fakeConfV2) OtpTTLMinutes() int                              { return 10 }
func (fakeConfV2) OtpVerifyMaxAttempts() int                       { return 5 }
func (fakeConfV2) OtpSendPerMinute() int                           { return 1 }
func (fakeConfV2) OtpSendPerHour() int                             { return 5 }
func (fakeConfV2) GetReviewDeadlineDays() int                      { panic("no") }
func (fakeConfV2) GetFrontendURL() string                          { panic("no") }
func (fakeConfV2) GetROIFrontendURL() string                       { panic("no") }
func (fakeConfV2) GetJwtAccessTokenExpireMinutes() int             { panic("no") }
func (fakeConfV2) GetJwtRefreshTokenExpireMinutes() int            { panic("no") }
func (fakeConfV2) GetCrossRefSenderEmail() string                  { panic("no") }
func (fakeConfV2) GetClientBasicAuthCredentials() (string, string) { panic("no") }
func (fakeConfV2) IsRefreshRotationStrict() bool                   { panic("no") }
func (fakeConfV2) GetRefreshRotationGraceSeconds() int             { panic("no") }
func (fakeConfV2) IsProduction() bool                              { panic("no") }

// Throttle bloklaganda SMS gateway CHAQIRILMAYDI va 429 qaytadi.
func TestSendOtpBlockedDoesNotSendSMS(t *testing.T) {
	sms := &fakeSms{}
	uc := authv2usecases.NewSendOtpUseCase(sms, nil, blockingThrottle()) // otp service nil — chaqirilmasligi kerak

	_, err := uc.Execute(context.Background(), "+998901112233", enum.OTPPurposeLogin)

	if !errors.Is(err, response.TooManyRequestsError) {
		t.Fatalf("TooManyRequestsError kutilgandi, %v keldi", err)
	}
	if sms.sends != 0 {
		t.Errorf("bloklanganda SMS yuborilmasligi kerak, %d yuborildi", sms.sends)
	}
}
