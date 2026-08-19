package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/usecase/authv2usecases"
	"slib.uz/src/core/domain/entity"
	apppkg "slib.uz/src/entrypoint/presentation/app"
	"slib.uz/src/entrypoint/presentation/handlers/authv2"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
)

// fakeSmsGateway — yuborilgan SMS sonini sanaydi. Malformed raqamda bu son
// nolda qolishi kerak: validatsiya SendOtpUseCase.Execute ga yetib
// bormasdan to'xtatishi kerak.
type fakeSmsGateway struct{ sends int }

func (f *fakeSmsGateway) Send(phone, message string) error { f.sends++; return nil }

// memoryLimiter — Redis o'rniga xotiradagi hisoblagich; hech qachon
// bloklamaydi (chegaradan past qiymat qaytaradi).
type memoryLimiter struct{ hits int }

func (l *memoryLimiter) Hit(ctx context.Context, key string, window time.Duration) (int64, error) {
	l.hits++
	return 1, nil
}
func (l *memoryLimiter) Reset(ctx context.Context, key string) error { return nil }

// memoryOtpRepo — DB o'rniga xotiradagi OTP repozitoriysi.
type memoryOtpRepo struct{ created int }

func (r *memoryOtpRepo) Create(ctx context.Context, otp *entity.OTPCodeEntity) (uint, error) {
	r.created++
	return 1, nil
}
func (r *memoryOtpRepo) GetByCodeAndSessionID(ctx context.Context, sessionID, code string) (*entity.OTPCodeEntity, error) {
	return nil, nil
}
func (r *memoryOtpRepo) MarkAsUsed(ctx context.Context, id uint) error { return nil }

// fakeConf — ConfigAdapter'ning shu testga kerakli qismini beradi.
// Ishlatilmaydigan metodlar panik qiladi — kutilmaganda chaqirilsa darhol
// bilinishi uchun.
type fakeConf struct{}

func (fakeConf) OtpTTLMinutes() int                              { return 10 }
func (fakeConf) OtpVerifyMaxAttempts() int                       { return 5 }
func (fakeConf) OtpSendPerMinute() int                           { return 1 }
func (fakeConf) OtpSendPerHour() int                             { return 5 }
func (fakeConf) GetReviewDeadlineDays() int                      { panic("no") }
func (fakeConf) GetFrontendURL() string                          { panic("no") }
func (fakeConf) GetROIFrontendURL() string                       { panic("no") }
func (fakeConf) GetJwtAccessTokenExpireMinutes() int             { panic("no") }
func (fakeConf) GetJwtRefreshTokenExpireMinutes() int            { panic("no") }
func (fakeConf) GetCrossRefSenderEmail() string                  { panic("no") }
func (fakeConf) GetClientBasicAuthCredentials() (string, string) { panic("no") }
func (fakeConf) IsRefreshRotationStrict() bool                   { panic("no") }
func (fakeConf) GetRefreshRotationGraceSeconds() int             { panic("no") }
func (fakeConf) IsProduction() bool                              { panic("no") }

// newSendOtpEcho send-otp yo'lini haqiqiy Echo ustida, haqiqiy validator,
// haqiqiy javob middleware'i va haqiqiy SendOtpUseCase bilan quradi. Yagona
// tashqi chegaralar — SMS gateway, rate limiter va OTP repozitoriysi —
// xotirada ishlaydigan soxta implementatsiyalar bilan almashtiriladi, shuning
// uchun DB yoki Redis shart emas.
func newSendOtpEcho(sms *fakeSmsGateway, limiter *memoryLimiter, repo *memoryOtpRepo) *echo.Echo {
	e := apppkg.NewEcho()
	e.Use((&middlewares.ResponseMiddleware{}).Call)

	throttle := service.NewAuthThrottle(limiter, fakeConf{})
	otpService := service.NewOTPService(repo, fakeConf{})
	uc := authv2usecases.NewSendOtpUseCase(sms, otpService, throttle)
	handler := authv2.NewSendOtpHandler(uc)

	e.POST("/api/auth-v2/send-otp", handler.Handle)
	return e
}

func postSendOtp(e *echo.Echo, phone string) *httptest.ResponseRecorder {
	body := strings.NewReader(`{"phone":"` + phone + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth-v2/send-otp", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// Spec markaziy da'vosi (3-bo'lim): noto'g'ri formatdagi telefon raqami
// SendOtpUseCase.Execute ishga tushmasdan rad etiladi — demak throttle
// sarflanmaydi, DB'ga OTP yozuvi tushmaydi va SMS provayderiga so'rov
// ketmaydi. Bu test Echo validator ulanishini (app.NewEcho), GetBody'dagi
// c.Validate chaqiruvini va javob middleware'ini bir butun sifatida
// sinaydi: agar ulardan biri o'chirilsa yoki almashtirilsa, bu test
// buziladi — hatto qolgan barcha birlik testlar o'tsa ham.
func TestSendOtpRejectsMalformedPhoneBeforeUseCase(t *testing.T) {
	sms := &fakeSmsGateway{}
	limiter := &memoryLimiter{}
	repo := &memoryOtpRepo{}
	e := newSendOtpEcho(sms, limiter, repo)

	rec := postSendOtp(e, "9989012345678")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status 400 kutilgandi, %d keldi (body: %s)", rec.Code, rec.Body.String())
	}
	if sms.sends != 0 {
		t.Errorf("SMS yuborilmasligi kerak edi, %d marta yuborildi", sms.sends)
	}
	if limiter.hits != 0 {
		t.Errorf("throttle sarflanmasligi kerak edi, %d marta hit bo'ldi", limiter.hits)
	}
	if repo.created != 0 {
		t.Errorf("OTP DB'ga yozilmasligi kerak edi, %d marta yozildi", repo.created)
	}
}

// To'g'ri formatdagi raqam xuddi shu simlash (wiring) bilan buzilmasligini
// tekshiradi: validator uni to'sib qo'ymaydi, so'rov handler tanasiga kirib,
// throttle, OTP yozuvi va SMS jo'natish haqiqatan ishlaydi (barchasi
// xotiradagi soxta implementatsiyalar ustida).
func TestSendOtpAcceptsWellFormedPhone(t *testing.T) {
	sms := &fakeSmsGateway{}
	limiter := &memoryLimiter{}
	repo := &memoryOtpRepo{}
	e := newSendOtpEcho(sms, limiter, repo)

	rec := postSendOtp(e, "998901234567")

	if rec.Code != http.StatusOK {
		t.Fatalf("status 200 kutilgandi, %d keldi (body: %s)", rec.Code, rec.Body.String())
	}
	if sms.sends != 1 {
		t.Errorf("SMS aynan bir marta yuborilishi kerak edi, %d marta yuborildi", sms.sends)
	}
}
