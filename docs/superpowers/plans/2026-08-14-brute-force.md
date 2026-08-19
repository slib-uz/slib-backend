# Brute-force himoyasi (CWE-307) — amalga oshirish rejasi

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** OTP-asosli login'ni brute-force'dan himoyalash — OTP tekshirish urinishlarini sessionID bo'yicha va OTP yuborishni telefon bo'yicha Redis atomik counter bilan cheklash.

**Architecture:** Yangi `RateLimiter` porti atomik `Hit`/`Reset` beradi; Redis `INCR`+`EXPIRE` uni amalga oshiradi. Markazlashtirilgan `AuthThrottle` service chegaralar va fail-open mantiqini bir joyda saqlaydi. Ikki use case (send-otp, verify) uni chaqiradi.

**Tech Stack:** Go 1.25, Echo v4, Redis (go-redis/v9), google/wire. Testlar — standart `testing`, qo'lda yozilgan fake'lar.

## Global Constraints

- **Yangi bog'liqlik yo'q.** Faqat standart kutubxona va mavjud paketlar (go-redis allaqachon bor).
- **Testlar ildizdagi `test/` katalogida**, `src/` tuzilishini takrorlaydi, `package <dir>_test`.
- **Core qatlami (`src/core/**`) infratuzilmani import qilmaydi.** `AuthThrottle` faqat portlar (`ratelimit`, `conf`) va `response` ni ishlatadi.
- **Core layer `github.com/labstack/gommon/log` ishlatadi, `zap` emas.**
- **Fail-open:** har `RateLimiter` xatosi loglanadi va `blocked=false` qaytaradi — login ishlaydi.
- **Yangi env `required` bo'lmasligi kerak** — `env/.env` git-ignored, mavjud deploy buzilmasin. `envDefault` ishlating.
- **Generatsiya:** `cmd/container/container.go` (wire), `docs/*` (swaggo) — qo'lda tegilmaydi. Use case konstruktori o'zgargani uchun `make wire-build` qayta generatsiya qiladi.
- **IDE diagnostikasi bu loyihada har doim eskirgan** — faqat `go build`/`go vet`/`go test`.
- **Shox:** `feature/brute-force`, `feature/sensitive-data` ustida. Bazaviy: **152 test**, `go build` toza.
- Izohlar va commit xabarlari o'zbek tilida.

---

## Fayl tuzilishi

**Yaratiladi:**

| Fayl | Mas'uliyati |
|---|---|
| `src/core/domain/ports/ratelimit/rate_limiter.go` | `RateLimiter` interfeysi (Hit, Reset) |
| `src/core/application/service/auth_throttle_service.go` | Chegaralar + fail-open, yuragi |
| `src/infrastructure/cache/redis_rate_limiter.go` | Redis INCR+EXPIRE impl |
| `test/core/application/service/auth_throttle_service_test.go` | AuthThrottle xatti-harakati (fake limiter) |
| `test/core/application/usecase/authv2usecases/send_otp_throttle_test.go` | send-otp integratsiyasi |
| `test/core/application/usecase/authv2usecases/verify_login_throttle_test.go` | verify integratsiyasi |

**O'zgartiriladi:**

| Fayl | O'zgarish |
|---|---|
| `src/core/application/response/response.go` | `TooManyAttemptsError`, `TooManyRequestsError` (429) |
| `src/core/domain/ports/conf/config_provider.go` | 3 getter |
| `src/infrastructure/config/env.go` | 3 envDefault maydon |
| `src/infrastructure/config/adapter.go` | 3 getter impl |
| `src/core/application/usecase/authv2usecases/send_otp_usecase.go` | throttle chaqiruvi |
| `src/core/application/usecase/authv2usecases/verify_and_login_usecase.go` | throttle chaqiruvi |

---

## Boshlashdan oldin: bazaviy holat

```bash
git branch --show-current     # feature/brute-force
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -5
```

152 test o'tishi kerak. O'tmasa — to'xtang.

---

### Task 1: Poydevor — port, 429 xatolar, config chegaralari

Kod yozilishidan oldingi skeleton: interfeys, xatolar, konfiguratsiya. Test yengil
(config getterlar).

**Files:**
- Create: `src/core/domain/ports/ratelimit/rate_limiter.go`
- Modify: `src/core/application/response/response.go`
- Modify: `src/core/domain/ports/conf/config_provider.go`
- Modify: `src/infrastructure/config/env.go`
- Modify: `src/infrastructure/config/adapter.go`

**Interfaces:**
- Produces:
  - `ratelimit.RateLimiter` (`Hit(ctx, key string, window time.Duration) (int64, error)`, `Reset(ctx, key string) error`)
  - `response.TooManyAttemptsError`, `response.TooManyRequestsError`
  - `conf.ConfigAdapter.OtpVerifyMaxAttempts() int`, `OtpSendPerMinute() int`, `OtpSendPerHour() int`

- [ ] **Step 1: `RateLimiter` portini yozish**

Yangi fayl `src/core/domain/ports/ratelimit/rate_limiter.go`:

```go
// Package ratelimit auth urinishlarini cheklash uchun atomik hisoblagich
// abstraksiyasini beradi. Implementatsiya (Redis) infratuzilma qatlamida.
package ratelimit

import (
	"context"
	"time"
)

type RateLimiter interface {
	// Hit kalitni atomik oshiradi va oshirilgandan keyingi joriy hisobni
	// qaytaradi. Kalit yangi bo'lsa (birinchi hit), window muddatiga TTL
	// o'rnatiladi.
	Hit(ctx context.Context, key string, window time.Duration) (int64, error)
	// Reset kalitni o'chiradi (muvaffaqiyatli login'dan keyin).
	Reset(ctx context.Context, key string) error
}
```

- [ ] **Step 2: 429 xatolarni qo'shish**

`src/core/application/response/response.go` — `EntityToLargeError` qatoridan keyin
(mavjud xatolar bloki ichida):

```go
	// Brute-force himoyasi (CWE-307): urinish/so'rov chegarasi oshdi.
	TooManyAttemptsError = NewFailResponse(429, "juda ko'p urinish, yangi kod so'rang")
	TooManyRequestsError = NewFailResponse(429, "juda ko'p so'rov, biroz kuting")
```

- [ ] **Step 3: Config portiga getterlar qo'shish**

`src/core/domain/ports/conf/config_provider.go` — interfeysga qo'shing (mavjud
`OtpTTLMinutes()` yonida):

```go
	OtpVerifyMaxAttempts() int
	OtpSendPerMinute() int
	OtpSendPerHour() int
```

- [ ] **Step 4: env maydonlarini qo'shish**

`src/infrastructure/config/env.go` — `OtpTtlMinutes` maydoni yonida (envDefault
naqshi bilan, `RefreshRotationStrict` kabi):

```go
	OtpVerifyMaxAttempts int `env:"OTP_VERIFY_MAX_ATTEMPTS" envDefault:"5"`
	OtpSendPerMinute     int `env:"OTP_SEND_PER_MINUTE" envDefault:"1"`
	OtpSendPerHour       int `env:"OTP_SEND_PER_HOUR" envDefault:"5"`
```

- [ ] **Step 5: adapter getterlarini qo'shish**

`src/infrastructure/config/adapter.go` — `OtpTTLMinutes` getteri yonida:

```go
func (this *ConfigAdapterImpl) OtpVerifyMaxAttempts() int {
	return this.env.OtpVerifyMaxAttempts
}

func (this *ConfigAdapterImpl) OtpSendPerMinute() int {
	return this.env.OtpSendPerMinute
}

func (this *ConfigAdapterImpl) OtpSendPerHour() int {
	return this.env.OtpSendPerHour
}
```

- [ ] **Step 6: Qurish va testlar**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -5
```

Kutilgan: qurish toza, 152 test (yangi test yo'q — bu skeleton). `ConfigAdapterImpl`
endi interfeysning barcha metodlarini qondiradi, aks holda `go build` yiqiladi.

- [ ] **Step 7: Commit**

```bash
git add src/core/domain/ports/ratelimit/rate_limiter.go \
        src/core/application/response/response.go \
        src/core/domain/ports/conf/config_provider.go \
        src/infrastructure/config/env.go \
        src/infrastructure/config/adapter.go
git commit -m "feat(cwe-307): RateLimiter porti, 429 xatolar, throttle chegaralari

Poydevor: atomik Hit/Reset porti, TooManyAttempts/TooManyRequests (429),
va OTP urinish/yuborish chegaralari (envDefault, mavjud deploy buzilmaydi)."
```

---

### Task 2: `AuthThrottle` service

Butun ishning yuragi: chegaralar, kalitlar, fail-open bir joyda. Fake `RateLimiter`
bilan to'liq testlanadi.

**Files:**
- Create: `src/core/application/service/auth_throttle_service.go`
- Test: `test/core/application/service/auth_throttle_service_test.go`

**Interfaces:**
- Consumes: `ratelimit.RateLimiter`, `conf.ConfigAdapter` (Task 1)
- Produces:
  - `service.NewAuthThrottle(limiter ratelimit.RateLimiter, conf conf.ConfigAdapter) *service.AuthThrottle`
  - `(*AuthThrottle).CheckAndHitOTPVerify(ctx, sessionID string) bool`
  - `(*AuthThrottle).ResetOTPVerify(ctx, sessionID string)`
  - `(*AuthThrottle).CheckAndHitOTPSend(ctx, phone string) bool`

- [ ] **Step 1: Yiqiladigan testni yozish**

Yangi fayl `test/core/application/service/auth_throttle_service_test.go`:

```go
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
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/core/application/service/... -run "Verify|Send|FailOpen" -count=1
```

Kutilgan: kompilyatsiya xatosi — `service.NewAuthThrottle` va `AuthThrottle`
aniqlanmagan.

- [ ] **Step 3: `AuthThrottle` service'ni yozish**

Yangi fayl `src/core/application/service/auth_throttle_service.go`:

```go
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
```

- [ ] **Step 4: Testlarni ishga tushirish**

```bash
go build ./... && go test ./test/core/application/service/... -run "Verify|Send|FailOpen" -count=1 -v
```

Kutilgan: 5 ta test PASS.

- [ ] **Step 5: Mutatsiya tekshiruvi**

Vaqtincha `CheckAndHitOTPVerify` da `count > int64(...)` ni `count >= int64(...)` ga
o'zgartiring:

```bash
go test ./test/core/application/service/... -run Verify -count=1
```

Kutilgan: `TestVerifyAllowsFiveBlocksSixth` yiqiladi (5-urinish endi bloklanadi).
**Yiqilmasa — to'xtang.** Keyin qaytaring.

Yana: fail-open mutatsiyasi — `CheckAndHitOTPVerify` dagi `return false` ni
`return true` ga o'zgartiring:

```bash
go test ./test/core/application/service/... -run FailOpen -count=1
```

Kutilgan: `TestFailOpenOnLimiterError` yiqiladi. Keyin qaytaring.

- [ ] **Step 6: Commit**

```bash
git add src/core/application/service/auth_throttle_service.go \
        test/core/application/service/auth_throttle_service_test.go
git commit -m "feat(cwe-307): AuthThrottle service — chegaralar va fail-open

OTP tekshirish (sessionID, 5 urinish) va yuborish (telefon, daqiqa/soat)
cheklovlari bir joyda. Redis xatosida fail-open: login ishlaydi, xato
loglanadi. Chegara va fail-open shoxlari mutatsiya bilan tasdiqlangan."
```

---

### Task 3: Redis implementatsiya

`RateLimiter` portining Redis amalga oshirilishi. Redis kerak bo'lgani uchun
testsiz (butun mantiq Task 2 da fake bilan qamralgan).

**Files:**
- Create: `src/infrastructure/cache/redis_rate_limiter.go`

**Interfaces:**
- Consumes: `ratelimit.RateLimiter` (Task 1)
- Produces: `cache.NewRedisRateLimiter(client *redis.Client) ratelimit.RateLimiter`

- [ ] **Step 1: Redis limiterni yozish**

Yangi fayl `src/infrastructure/cache/redis_rate_limiter.go` (`RedisCache` naqshiga
mos):

```go
package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"slib.uz/src/core/domain/ports/ratelimit"
)

type RedisRateLimiter struct {
	client *redis.Client
}

// @inject
func NewRedisRateLimiter(client *redis.Client) ratelimit.RateLimiter {
	return &RedisRateLimiter{client: client}
}

// Hit kalitni atomik oshiradi. Natija 1 bo'lsa (yangi kalit), window muddatiga
// TTL o'rnatiladi. INCR atomik, shuning uchun ikki parallel urinish ham to'g'ri
// sanoq oladi.
func (this *RedisRateLimiter) Hit(ctx context.Context, key string, window time.Duration) (int64, error) {
	count, err := this.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		if err := this.client.Expire(ctx, key, window).Err(); err != nil {
			return count, err
		}
	}
	return count, nil
}

func (this *RedisRateLimiter) Reset(ctx context.Context, key string) error {
	return this.client.Del(ctx, key).Err()
}
```

- [ ] **Step 2: Qurish**

```bash
go build ./... && go vet ./...
```

Kutilgan: toza. `NewRedisRateLimiter` `ratelimit.RateLimiter` ni qondiradi.

- [ ] **Step 3: Commit**

```bash
git add src/infrastructure/cache/redis_rate_limiter.go
git commit -m "feat(cwe-307): RateLimiter Redis implementatsiyasi

INCR atomik oshiradi, birinchi hit'da EXPIRE o'rnatadi. Reset kalitni
o'chiradi. Mavjud Redis client (DB=1) qayta ishlatiladi."
```

---

### Task 4: Use case integratsiyasi

Ikki use case throttle'ni chaqiradi. Fake throttle bilan integratsiya testlanadi.

**Files:**
- Modify: `src/core/application/usecase/authv2usecases/send_otp_usecase.go`
- Modify: `src/core/application/usecase/authv2usecases/verify_and_login_usecase.go`
- Test: `test/core/application/usecase/authv2usecases/send_otp_throttle_test.go` (create)
- Test: `test/core/application/usecase/authv2usecases/verify_login_throttle_test.go` (create)

**Interfaces:**
- Consumes: `service.AuthThrottle` (Task 2), `response.TooManyAttemptsError`, `response.TooManyRequestsError` (Task 1)
- Produces: o'zgargan konstruktorlar (wire qayta generatsiya qiladi)

**MUHIM:** `AuthThrottle` — konkret `*service.AuthThrottle` (interfeys emas), va
use case'lar `service` paketida joylashgan. Test uchun fake yozib bo'lmaydi
(konkret tur). Shuning uchun test **haqiqiy `AuthThrottle` + fake `RateLimiter`**
bilan quriladi — bu integratsiyani ham qamraydi. `send_otp` va `verify` ning
throttle'dan boshqa bog'liqliklari (sms gateway, otp service, repository) fake.

- [ ] **Step 1: send-otp integratsiya testini yozish**

Avval `send_otp_usecase` va `verify_and_login_usecase` ning mavjud bog'liqliklarini
o'qing (`gateway.SmsGateway`, `*service.OTPService`, `repository.UserRepository` va
h.k.) va ularning fake'larini yozing. `OTPService` konkret tur — u
`repository.OTPCodeRepository` va `conf.ConfigAdapter` oladi, shuning uchun uni
haqiqiy `OTPService` + fake OTP repo bilan quring.

Yangi fayl `test/core/application/usecase/authv2usecases/send_otp_throttle_test.go`:

```go
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

func (fakeConfV2) OtpTTLMinutes() int        { return 10 }
func (fakeConfV2) OtpVerifyMaxAttempts() int { return 5 }
func (fakeConfV2) OtpSendPerMinute() int     { return 1 }
func (fakeConfV2) OtpSendPerHour() int       { return 5 }
func (fakeConfV2) GetReviewDeadlineDays() int              { panic("no") }
func (fakeConfV2) GetFrontendURL() string                  { panic("no") }
func (fakeConfV2) GetROIFrontendURL() string               { panic("no") }
func (fakeConfV2) GetJwtAccessTokenExpireMinutes() int     { panic("no") }
func (fakeConfV2) GetJwtRefreshTokenExpireMinutes() int    { panic("no") }
func (fakeConfV2) GetCrossRefSenderEmail() string          { panic("no") }
func (fakeConfV2) GetClientBasicAuthCredentials() (string, string) { panic("no") }
func (fakeConfV2) IsRefreshRotationStrict() bool           { panic("no") }
func (fakeConfV2) GetRefreshRotationGraceSeconds() int     { panic("no") }
func (fakeConfV2) IsProduction() bool                      { panic("no") }

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
```

**Diqqat:** `NewSendOtpUseCase` imzosi 3-qadamda o'zgaradi (`throttle` qo'shiladi).
Test shu yangi imzoga tayanadi. `otp service nil` — throttle avval bloklagani
uchun `Make` chaqirilmaydi (bloklangan yo'lda nil xavfsiz).

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/core/application/usecase/authv2usecases/... -run SendOtpBlocked -count=1
```

Kutilgan: kompilyatsiya xatosi — `NewSendOtpUseCase` hali 2 argument oladi.

- [ ] **Step 3: send-otp use case'ni o'zgartirish**

`send_otp_usecase.go`:

```go
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
```

Import blokiga `"slib.uz/src/core/application/response"` qo'shing.

- [ ] **Step 4: send testini o'tkazish**

```bash
go build ./... && go test ./test/core/application/usecase/authv2usecases/... -run SendOtpBlocked -count=1
```

Kutilgan: PASS.

- [ ] **Step 5: verify integratsiya testini yozish**

Yangi fayl `test/core/application/usecase/authv2usecases/verify_login_throttle_test.go`.
`VerifyAndLoginUseCase` bog'liqliklari (`repository.UserRepository`,
`repository.UserProfileRepository`, `*service.OTPService`, `*service.UserAuthTokenService`,
`session.Atomic`) — fake yoki nil. Throttle bloklaganда `OTPService.Check`
chaqirilmasligini tekshiramiz, shuning uchun `OTPService` nil bo'lishi mumkin (blok
yo'lida chaqirilmaydi):

```go
package authv2usecases_test

import (
	"context"
	"errors"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/authv2usecases"
)

// Throttle bloklaganda OTP Check chaqirilmaydi va 429 qaytadi.
func TestVerifyBlockedReturnsError(t *testing.T) {
	// blockingThrottle() send_otp_throttle_test.go dagi yordamchi.
	// Boshqa bog'liqliklar nil — blok yo'lida chaqirilmaydi.
	uc := authv2usecases.NewVerifyAndLoginUseCase(nil, nil, nil, nil, nil, blockingThrottle())

	_, _, err := uc.Execute(context.Background(), "sess", "123456", "", "")

	if !errors.Is(err, response.TooManyAttemptsError) {
		t.Fatalf("TooManyAttemptsError kutilgandi, %v keldi", err)
	}
}
```

**Diqqat:** `blockingThrottle()` va `fakeConfV2` allaqachon `send_otp_throttle_test.go`
da e'lon qilingan (bir xil `authv2usecases_test` paketi) — qayta e'lon qilmang.

- [ ] **Step 6: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/core/application/usecase/authv2usecases/... -run VerifyBlocked -count=1
```

Kutilgan: kompilyatsiya xatosi — `NewVerifyAndLoginUseCase` hali 5 argument oladi.

- [ ] **Step 7: verify use case'ni o'zgartirish**

`verify_and_login_usecase.go`:

```go
type VerifyAndLoginUseCase struct {
	repository        repository.UserRepository
	profileRepository repository.UserProfileRepository
	service           *service.OTPService
	authTokenService  *service.UserAuthTokenService
	atomic            session.Atomic
	throttle          *service.AuthThrottle
}

// @inject
func NewVerifyAndLoginUseCase(repository repository.UserRepository, profileRepository repository.UserProfileRepository, service *service.OTPService, authTokenService *service.UserAuthTokenService, atomic session.Atomic, throttle *service.AuthThrottle) *VerifyAndLoginUseCase {
	return &VerifyAndLoginUseCase{repository: repository, profileRepository: profileRepository, service: service, authTokenService: authTokenService, atomic: atomic, throttle: throttle}
}
```

`Execute` boshini o'zgartiring (`Check`dan oldin throttle, muvaffaqiyatda reset):

```go
func (this *VerifyAndLoginUseCase) Execute(ctx context.Context, sessionID string, code string, scienceID, fullName string) (*entity.AuthTokenEntity, enum.AuthScope, error) {
	if this.throttle.CheckAndHitOTPVerify(ctx, sessionID) {
		return nil, "", response.TooManyAttemptsError
	}

	otp, err := this.service.Check(ctx, code, sessionID)
	if err != nil {
		return nil, "", err
	}

	this.throttle.ResetOTPVerify(ctx, sessionID)

	user, err := this.repository.GetByPhoneNumber(otp.Phone)
	// ... qolgani o'zgarmaydi ...
```

`response` import allaqachon bor (faylda ishlatiladi).

- [ ] **Step 8: Qurish va testlar**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -20
```

Kutilgan: qurish toza, 159 test (157 + 2 yangi). Wire hali qo'lda emas — Task 5 da
generatsiya.

- [ ] **Step 9: Commit**

```bash
git add src/core/application/usecase/authv2usecases/send_otp_usecase.go \
        src/core/application/usecase/authv2usecases/verify_and_login_usecase.go \
        test/core/application/usecase/authv2usecases/send_otp_throttle_test.go \
        test/core/application/usecase/authv2usecases/verify_login_throttle_test.go
git commit -m "fix(cwe-307): send-otp va verify use case'lari throttle bilan himoyalandi

send-otp: bloklanganda SMS yuborilmaydi (429). verify: Check'dan oldin
urinish sanaladi, kod to'g'ri bo'lsa hisob tozalanadi. Auditning aynan
topilmasi (OTP brute-force) yopildi."
```

---

### Task 5: Wire va yakuniy tekshiruv

**Files:** faqat generatsiya (`cmd/container/container.go`) va tekshirish

- [ ] **Step 1: Wire qayta generatsiya**

Use case konstruktorlari va yangi `AuthThrottle`/`RedisRateLimiter` `@inject` bilan
belgilangan — wire ularni ulashi kerak.

```bash
make wire-build
go build ./...
```

Kutilgan: qurish toza. `container.go` o'zgaradi (yangi provayderlar). Agar wire
`AuthThrottle` yoki `RateLimiter` ni ulay olmasa (`@inject` yo'q yoki bog'lanmagan),
xatoni o'qing — `NewAuthThrottle`, `NewRedisRateLimiter` `// @inject` izohiga ega
ekanini tasdiqlang.

- [ ] **Step 2: Generatorlar idempotent**

```bash
git add -A
make wire-build && make generate-docs
git diff --stat
```

Kutilgan: `git diff --stat` bo'sh (birinchi generatsiyadan keyin qayta generatsiya
o'zgarish bermaydi).

- [ ] **Step 3: To'liq test to'plami**

```bash
go build ./... && go vet ./...
go test ./... -count=1
```

Kutilgan: 159 test, 0 yiqilish, ishchi daraxt toza (wire commit qilingandan keyin).

- [ ] **Step 4: Wire natijasini commit qilish**

```bash
git add cmd/container/container.go
git commit -m "chore(cwe-307): wire — throttle va rate limiter ulandi"
```

- [ ] **Step 5: Spec bilan solishtirish**

`docs/superpowers/specs/2026-08-14-brute-force-design.md` ni ochib har bir bo'limni
bajarilgan ish bilan solishtiring. Farq topilsa — kodni emas, avval farqni xabar
qiling.

- [ ] **Step 6: Shoxni yakunlash**

**REQUIRED SUB-SKILL:** `superpowers:finishing-a-development-branch`

Bazaviy shox: `feature/sensitive-data` (zanjir #18 → #21 → #22 → #24 → shu ish).

---

## Kutilayotgan yakuniy holat

| Ko'rsatkich | Boshlanish | Yakun |
|---|---|---|
| Testlar | 152 | 159 |
| Cheklovsiz OTP brute-force | bor | yopilgan (5 urinish) |
| Cheksiz SMS spam | bor | yopilgan (daqiqa/soat) |
| Yangi bog'liqlik | — | 0 |
| Yangi env (required) | — | 0 (envDefault) |

## Deploy oldidan (kod ishi emas)

- Redis DB=1 mavjud va ishlayotganini tasdiqlang (fail-open bo'lsa ham, himoya
  faqat Redis bilan ishlaydi).
- Ixtiyoriy: `.env` ga `OTP_VERIFY_MAX_ATTEMPTS`, `OTP_SEND_PER_MINUTE`,
  `OTP_SEND_PER_HOUR` qo'shib chegaralarni sozlash (standart 5/1/5).
- Frontend jamoasi 429 javoblarini (`TooManyAttemptsError`, `TooManyRequestsError`)
  to'g'ri ko'rsatishini tekshiring.
