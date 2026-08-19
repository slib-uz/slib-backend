package service

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/ratelimit"
)

// AuthThrottle OTP autentifikatsiyasini brute-force'dan himoyalaydi (CWE-307).
// Chegaralar, kalitlar va fail-open mantiqi shu yerda markazlashtirilgan.
//
// Fail-open: Redis (RateLimiter) xatosi login'ni to'xtatmaydi — xato loglanadi
// va urinishga ruxsat beriladi. Aks holda hujumchi Redis'ni ishdan chiqarib
// butun login'ni to'xtatishi mumkin edi (mavjudlik hujumi).
type AuthThrottle struct {
	limiter ratelimit.RateLimiter
	conf    conf.ConfigAdapter
}

// @inject
func NewAuthThrottle(limiter ratelimit.RateLimiter, conf conf.ConfigAdapter) *AuthThrottle {
	return &AuthThrottle{limiter: limiter, conf: conf}
}

// CheckAndHitOTPVerify urinishni sanaydi va limit oshsa true qaytaradi.
// Tekshiruvdan OLDIN chaqiriladi; kod to'g'ri chiqsa chaqiruvchi ResetOTPVerify
// bilan hisobni tozalaydi.
func (this *AuthThrottle) CheckAndHitOTPVerify(ctx context.Context, sessionID string) bool {
	key := "auth:otp:verify:" + sessionID
	window := time.Duration(this.conf.OtpTTLMinutes()) * time.Minute
	count, err := this.limiter.Hit(ctx, key, window)
	if err != nil {
		log.Warnf("[AuthThrottle] rate limiter xatosi (verify), fail-open: %v", err)
		return false
	}
	return count > int64(this.conf.OtpVerifyMaxAttempts())
}

// ResetOTPVerify muvaffaqiyatli login'dan keyin urinish hisobini tozalaydi.
func (this *AuthThrottle) ResetOTPVerify(ctx context.Context, sessionID string) {
	if err := this.limiter.Reset(ctx, "auth:otp:verify:"+sessionID); err != nil {
		log.Warnf("[AuthThrottle] rate limiter reset xatosi (verify): %v", err)
	}
}

// CheckAndHitOTPSend telefon bo'yicha ikki oynani tekshiradi: daqiqalik va soatlik.
func (this *AuthThrottle) CheckAndHitOTPSend(ctx context.Context, phone string) bool {
	minCount, err := this.limiter.Hit(ctx, "auth:otp:send:min:"+phone, time.Minute)
	if err != nil {
		log.Warnf("[AuthThrottle] rate limiter xatosi (send/min), fail-open: %v", err)
		return false
	}
	if minCount > int64(this.conf.OtpSendPerMinute()) {
		return true
	}

	hourCount, err := this.limiter.Hit(ctx, "auth:otp:send:hour:"+phone, time.Hour)
	if err != nil {
		log.Warnf("[AuthThrottle] rate limiter xatosi (send/hour), fail-open: %v", err)
		return false
	}
	return hourCount > int64(this.conf.OtpSendPerHour())
}
