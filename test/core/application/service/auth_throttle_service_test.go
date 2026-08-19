package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"slib.uz/src/core/application/service"
)

// fakeLimiter — RateLimiter ning qo'lda yozilgan soxtasi. Har kalit uchun
// hisobni saqlaydi; failKeys da bo'lgan kalit uchun xato qaytaradi.
type fakeLimiter struct {
	counts   map[string]int64
	resets   map[string]int
	failHit  bool
	failKeys map[string]bool
}

func newFakeLimiter() *fakeLimiter {
	return &fakeLimiter{counts: map[string]int64{}, resets: map[string]int{}, failKeys: map[string]bool{}}
}

func (f *fakeLimiter) Hit(ctx context.Context, key string, window time.Duration) (int64, error) {
	if f.failHit || f.failKeys[key] {
		return 0, errors.New("redis down")
	}
	f.counts[key]++
	return f.counts[key], nil
}

func (f *fakeLimiter) Reset(ctx context.Context, key string) error {
	f.resets[key]++
	delete(f.counts, key)
	return nil
}

// fakeConf — ConfigAdapter ning throttle uchun kerakli qismini beradi.
// Boshqa metodlar chaqirilmaydi (chaqirilsa panik).
type fakeConf struct {
	ttl, verifyMax, sendMin, sendHour int
}

func (c fakeConf) OtpTTLMinutes() int         { return c.ttl }
func (c fakeConf) OtpVerifyMaxAttempts() int  { return c.verifyMax }
func (c fakeConf) OtpSendPerMinute() int      { return c.sendMin }
func (c fakeConf) OtpSendPerHour() int        { return c.sendHour }
func (c fakeConf) GetReviewDeadlineDays() int { panic("chaqirilmasligi kerak") }
func (c fakeConf) GetFrontendURL() string     { panic("chaqirilmasligi kerak") }
func (c fakeConf) GetROIFrontendURL() string  { panic("chaqirilmasligi kerak") }
func (c fakeConf) GetJwtAccessTokenExpireMinutes() int  { panic("chaqirilmasligi kerak") }
func (c fakeConf) GetJwtRefreshTokenExpireMinutes() int { panic("chaqirilmasligi kerak") }
func (c fakeConf) GetCrossRefSenderEmail() string       { panic("chaqirilmasligi kerak") }
func (c fakeConf) GetClientBasicAuthCredentials() (string, string) { panic("chaqirilmasligi kerak") }
func (c fakeConf) IsRefreshRotationStrict() bool        { panic("chaqirilmasligi kerak") }
func (c fakeConf) GetRefreshRotationGraceSeconds() int  { panic("chaqirilmasligi kerak") }
func (c fakeConf) IsProduction() bool                   { panic("chaqirilmasligi kerak") }

func newThrottle(l *fakeLimiter) *service.AuthThrottle {
	return service.NewAuthThrottle(l, fakeConf{ttl: 10, verifyMax: 5, sendMin: 1, sendHour: 5})
}

// verify: 5 urinishga ruxsat, 6-chidan bloklash.
func TestVerifyAllowsFiveBlocksSixth(t *testing.T) {
	th := newThrottle(newFakeLimiter())
	for i := 1; i <= 5; i++ {
		if th.CheckAndHitOTPVerify(context.Background(), "sess") {
			t.Fatalf("%d-urinish bloklanmasligi kerak edi", i)
		}
	}
	if !th.CheckAndHitOTPVerify(context.Background(), "sess") {
		t.Error("6-urinish bloklanishi kerak edi")
	}
}

// verify reset: muvaffaqiyatdan keyin hisob tozalanadi.
func TestVerifyResetClearsCount(t *testing.T) {
	l := newFakeLimiter()
	th := newThrottle(l)
	th.CheckAndHitOTPVerify(context.Background(), "sess")
	th.ResetOTPVerify(context.Background(), "sess")
	if l.resets["auth:otp:verify:sess"] != 1 {
		t.Errorf("reset chaqirilishi kerak edi, %d marta chaqirildi", l.resets["auth:otp:verify:sess"])
	}
	if l.counts["auth:otp:verify:sess"] != 0 {
		t.Error("reset'dan keyin hisob 0 bo'lishi kerak edi")
	}
}

// send: bir daqiqada 1 so'rovga ruxsat, 2-si bloklanadi.
func TestSendBlocksSecondInMinute(t *testing.T) {
	th := newThrottle(newFakeLimiter())
	if th.CheckAndHitOTPSend(context.Background(), "+998901112233") {
		t.Error("1-so'rov o'tishi kerak edi")
	}
	if !th.CheckAndHitOTPSend(context.Background(), "+998901112233") {
		t.Error("2-so'rov (bir daqiqada) bloklanishi kerak edi")
	}
}

// fail-open: limiter xatosi bo'lsa hech narsa bloklanmaydi.
func TestFailOpenOnLimiterError(t *testing.T) {
	l := newFakeLimiter()
	l.failHit = true
	th := newThrottle(l)
	if th.CheckAndHitOTPVerify(context.Background(), "sess") {
		t.Error("limiter xatosida verify bloklanmasligi kerak (fail-open)")
	}
	if th.CheckAndHitOTPSend(context.Background(), "+998901112233") {
		t.Error("limiter xatosida send bloklanmasligi kerak (fail-open)")
	}
}

// fail-open: faqat soatlik limiter xato bersa (daqiqalik normal ishlasa),
// daqiqalik hit o'tadi va soatlik hit xato beradi — fail-open shoxi ishga tushadi.
func TestFailOpenOnHourlyLimiterError(t *testing.T) {
	l := newFakeLimiter()
	phone := "+998901112233"
	l.failKeys["auth:otp:send:hour:"+phone] = true
	th := newThrottle(l)
	if th.CheckAndHitOTPSend(context.Background(), phone) {
		t.Error("soatlik limiter xatosida send bloklanmasligi kerak (fail-open)")
	}
}

// send: daqiqalik oyna o'tsa ham soatlik chegara bloklaydi.
// Daqiqalik kalitni har safar reset qilib, soatlik hisob 5 dan oshishini sinaymiz.
func TestSendHourlyLimitBlocks(t *testing.T) {
	l := newFakeLimiter()
	th := newThrottle(l)
	phone := "+998901112233"
	for i := 1; i <= 5; i++ {
		if th.CheckAndHitOTPSend(context.Background(), phone) {
			t.Fatalf("%d-so'rov (soatlik chegara ichida) o'tishi kerak edi", i)
		}
		// daqiqalik kalitni tozalab, keyingi daqiqani simulyatsiya qilamiz
		l.Reset(context.Background(), "auth:otp:send:min:"+phone)
	}
	if !th.CheckAndHitOTPSend(context.Background(), phone) {
		t.Error("6-so'rov soatlik chegara bilan bloklanishi kerak edi")
	}
}
