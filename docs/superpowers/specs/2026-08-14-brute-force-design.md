# Brute-force himoyasi (CWE-307) — dizayn spetsifikatsiyasi

**Sana:** 2026-08-14
**Audit topilmasi:** 2.2.7 To'liq ajratib bajariladigan hujum (brute-force), xavflilik darajasi **O'rta**
**CWE:** CWE-307 — Improper Restriction of Excessive Authentication Attempts
**OWASP Top 2025:** 7-o'rin (Authentication Failures)
**Shox:** `feature/brute-force`, `feature/sensitive-data` ustiga (zanjirning 6-halqasi)

---

## 1. Zaiflik

Ekspertiza kirish oynasi (`journal.sciencelib.uz`) brute-force'dan himoyalanmaganini
topgan: 10000+ parol kombinatsiyasi cheklovsiz sinovdan o'tkazilgan.

Loyihaning autentifikatsiyasi **parol emas, OTP-asosli** (`/auth-v2/*`):

- `POST /auth-v2/send-otp` — 6 xonali OTP generatsiya qiladi, SMS orqali yuboradi,
  `sessionID` (32-bayt tasodifiy hex) qaytaradi.
- `POST /auth-v2/login` — `sessionID` + `code` bilan tekshiradi
  (`OTPService.Check` → `GetByCodeAndSessionID(sessionID, code)`).

### 1.1. Ikki brute-force yuzasi

**A. OTP kodni brute-force qilish** — auditning aynan topilmasi.
Hujumchi `sessionID` ni send-otp javobidan oladi (o'z yoki qurbon telefoni bilan),
so'ng shu session uchun barcha **1,000,000** olti xonali kodni sinaydi. Har noto'g'ri
urinish shunchaki `invalid code` qaytaradi — **hech qanday cheklov yo'q**
(`otp_service.go:44-71`). 6 xonali kod + cheklovsiz urinish = amaliy jihatdan
buziladigan.

**B. OTP yuborishni spam qilish.**
`POST /auth-v2/send-otp` cheksiz chaqirilishi mumkin — har chaqiruv SMS yuboradi
(`smsetc_gateway`). Bu qurbonni bezovta qiladi va SMS byudjetini yoqib yuboradi.
Cheklov yo'q.

### 1.2. Hozirgi holat

- **Rate limiting yo'q** — na Echo middleware, na Redis counter, na DB urinish jadvali.
- Redis loyihada bor (`redis_client.go`, DB=1), lekin faqat views count uchun.
- `c.RealIP()` faqat logging'da ishlatiladi, auth'da emas.
- `CacheProvider` porti faqat `Get`/`Set` — atomik `INCR` yo'q.

---

## 2. Qamrov

**Ichida:**
- OTP tekshirish (`/auth-v2/login`) — sessionID bo'yicha urinishlarni cheklash
- OTP yuborish (`/auth-v2/send-otp`) — telefon bo'yicha SMS spam cheklovi

**Qamrovdan tashqarida** (foydalanuvchi qarori):
- `GET /auth-v2/check-phone-number` telefon enumeratsiyasi — CWE-307 emas (CWE-200
  sinfi), ataylab ommaviy bo'lishi mumkin, alohida ish.
- `/auth/oauth/*` (ScienceID OAuth) — parol tekshiruvi ScienceID tomonida, bu yerda
  brute-force yuzasi yo'q.

---

## 3. Yechim

### 3.1. `RateLimiter` porti

Yangi port: `src/core/domain/ports/ratelimit/rate_limiter.go`.

Hozirgi `CacheProvider` faqat Get/Set beradi; rate limiting atomik oshirishni talab
qiladi (Get→parse→Set poyga beradi va ikki parallel urinish bir xil hisobni
ko'radi). Alohida port semantikani toza ifodalaydi:

```go
package ratelimit

import (
    "context"
    "time"
)

type RateLimiter interface {
    // Hit kalitni atomik oshiradi va oshirilgandan keyingi joriy hisobni qaytaradi.
    // Kalit yangi bo'lsa (birinchi hit), window muddatiga TTL o'rnatiladi.
    Hit(ctx context.Context, key string, window time.Duration) (int64, error)
    // Reset kalitni o'chiradi (muvaffaqiyatli login'dan keyin hisobni tozalash).
    Reset(ctx context.Context, key string) error
}
```

### 3.2. Redis implementatsiya

`src/infrastructure/cache/redis_rate_limiter.go`. Mavjud Redis client qayta
ishlatiladi. `Hit` atomik: `INCR`, va natija `1` bo'lsa (yangi kalit) `EXPIRE`
o'rnatiladi — Redis `pipeline` yoki `INCR` keyin shartli `EXPIRE` bilan.

```go
func (r *RedisRateLimiter) Hit(ctx context.Context, key string, window time.Duration) (int64, error) {
    count, err := r.client.Incr(ctx, key).Result()
    if err != nil {
        return 0, err
    }
    if count == 1 {
        // Yangi kalit — TTL o'rnatamiz. Xatosi kritik emas: kalit TTL'siz qolsa
        // ham keyingi oynada Reset yoki qo'lda tozalash ishlaydi; lekin xato
        // qaytariladi, throttle service fail-open qiladi.
        if err := r.client.Expire(ctx, key, window).Err(); err != nil {
            return count, err
        }
    }
    return count, nil
}
```

### 3.3. `AuthThrottle` service

`src/core/application/service/auth_throttle_service.go`. Chegaralar, kalitlar va
**fail-open mantiqi bir joyda**.

```go
type AuthThrottle struct {
    limiter ratelimit.RateLimiter
    conf    conf.ConfigAdapter
}

// CheckAndHitOTPVerify urinishni sanaydi va limit oshsa bloklaydi.
// Har chaqiruvda hit qiladi (tekshiruvdan OLDIN chaqiriladi); kod to'g'ri chiqsa
// chaqiruvchi ResetOTPVerify bilan hisobni tozalaydi.
func (t *AuthThrottle) CheckAndHitOTPVerify(ctx context.Context, sessionID string) (blocked bool) {
    key := "auth:otp:verify:" + sessionID
    window := time.Duration(t.conf.OtpTTLMinutes()) * time.Minute
    count, err := t.limiter.Hit(ctx, key, window)
    if err != nil {
        log.Warnf("rate limiter xatosi (verify), fail-open: %v", err)
        return false // fail-open: login ishlaydi
    }
    return count > int64(t.conf.OtpVerifyMaxAttempts())
}

func (t *AuthThrottle) ResetOTPVerify(ctx context.Context, sessionID string) {
    if err := t.limiter.Reset(ctx, "auth:otp:verify:"+sessionID); err != nil {
        log.Warnf("rate limiter reset xatosi (verify): %v", err)
    }
}

// CheckAndHitOTPSend telefon bo'yicha ikki oynani tekshiradi: daqiqalik va soatlik.
func (t *AuthThrottle) CheckAndHitOTPSend(ctx context.Context, phone string) (blocked bool) {
    minKey := "auth:otp:send:min:" + phone
    minCount, err := t.limiter.Hit(ctx, minKey, time.Minute)
    if err != nil {
        log.Warnf("rate limiter xatosi (send/min), fail-open: %v", err)
        return false
    }
    if minCount > int64(t.conf.OtpSendPerMinute()) {
        return true
    }

    hourKey := "auth:otp:send:hour:" + phone
    hourCount, err := t.limiter.Hit(ctx, hourKey, time.Hour)
    if err != nil {
        log.Warnf("rate limiter xatosi (send/hour), fail-open: %v", err)
        return false
    }
    return hourCount > int64(t.conf.OtpSendPerHour())
}
```

**Fail-open:** har `RateLimiter` xatosi `log.Warnf` (gommon/log — core layer `zap`
ishlatmaydi) bilan qayd etiladi va `blocked=false` qaytaradi. Redis uzilishi
login'ni to'xtatmaydi.

**Muhim — hit tartibi va chegara qat'iy belgilangan:** `CheckAndHitOTPVerify`
tekshiruvdan **oldin** hit qiladi va `count > OtpVerifyMaxAttempts` bo'lsa bloklaydi.
`OtpVerifyMaxAttempts=5` bilan: 5-urinish `count=5` (`5 > 5` yolg'on → o'tadi),
6-urinish `count=6` (`6 > 5` rost → bloklanadi). Ya'ni **5 urinishga ruxsat,
6-chidan bloklash**. Bu qat'iy semantika — testlar aynan shunga tayanadi.

### 3.4. Chegaralar (config)

`ConfigAdapter` portiga (`src/core/domain/ports/conf/config_provider.go`) qo'shiladi:

```go
OtpVerifyMaxAttempts() int
OtpSendPerMinute() int
OtpSendPerHour() int
```

`env.go` da `envDefault` bilan (yangi env `required` bo'lmasligi kerak — `env/.env`
git-ignored, mavjud deploy buzilmasligi uchun):

```go
OtpVerifyMaxAttempts int `env:"OTP_VERIFY_MAX_ATTEMPTS" envDefault:"5"`
OtpSendPerMinute     int `env:"OTP_SEND_PER_MINUTE" envDefault:"1"`
OtpSendPerHour       int `env:"OTP_SEND_PER_HOUR" envDefault:"5"`
```

`adapter.go` da uch getter qo'shiladi.

### 3.5. Use case integratsiyasi

**`send_otp_usecase.Execute(ctx, phoneNumber, purpose)`** — boshida:

```go
if this.throttle.CheckAndHitOTPSend(ctx, phoneNumber) {
    return "", response.TooManyRequestsError
}
// ... mavjud OTP generatsiya + SMS ...
```

Bloklangan bo'lsa SMS **yuborilmaydi** (429 qaytadi).

**`verify_and_login_usecase.Execute(ctx, sessionID, code, ...)`** — `OTPService.Check`
dan **oldin**:

```go
if this.throttle.CheckAndHitOTPVerify(ctx, sessionID) {
    return nil, 0, response.TooManyAttemptsError
}

result, err := this.otpService.Check(ctx, code, sessionID)
if err != nil {
    return nil, 0, err   // noto'g'ri kod — hit qoladi
}
this.throttle.ResetOTPVerify(ctx, sessionID)  // to'g'ri kod — hisob tozalanadi
// ... mavjud login ...
```

Ikkala use case konstruktori `*AuthThrottle` oladi (`@inject`, wire qayta
generatsiya qiladi).

### 3.6. Xatolik javoblari

`src/core/application/response/response.go` ga qo'shiladi:

```go
TooManyAttemptsError = NewFailResponse(429, "juda ko'p urinish, yangi kod so'rang")
TooManyRequestsError = NewFailResponse(429, "juda ko'p so'rov, biroz kuting")
```

(Hozirgi tekshiruvda `response` paketida 429 topilmadi — shuning uchun ikkalasi
yangi qo'shiladi.)

---

## 4. Testlar

Loyihaning mavjud tartibi: testlar ildizdagi `test/` katalogida, `src/` tuzilishini
takrorlaydi, `package <dir>_test`. Yangi bog'liqlik yo'q — standart `testing` va
qo'lda yozilgan fake'lar.

**`AuthThrottle` service** (fake `RateLimiter` bilan — Redis kerak emas):

- verify: 1–5-urinish o'tadi, 6-urinishда bloklanadi
- verify reset: `ResetOTPVerify` fake'da `Reset` ni chaqiradi (kalit o'chiriladi)
- send/daqiqa: 1-so'rov o'tadi, 2-si (bir daqiqada) bloklanadi
- send/soat: daqiqalik oyna ostida, soatlik chegara alohida bloklaydi
- **fail-open:** fake `RateLimiter` xato qaytarsa → `blocked=false` (uchala metod uchun)

**Use case integratsiyasi** (fake throttle + fake OTP service/gateway):

- `send_otp`: throttle bloklaganда SMS gateway **chaqirilmaydi**, 429 qaytadi
- `verify_and_login`: throttle bloklaganда `OTPService.Check` **chaqirilmaydi**;
  kod to'g'ri bo'lganda `ResetOTPVerify` chaqiriladi

**`redis_rate_limiter`** — Redis talab qiladi, loyihada Redis bilan test yo'q,
shuning uchun testlanmaydi. Butun qaror `AuthThrottle` da va u fake bilan to'liq
qamraladi.

**Mutatsiya tekshiruvi:** `AuthThrottle` ning chegara taqqoslashi (`count > max`) va
fail-open shoxi mutatsiya bilan tasdiqlanadi — testlar haqiqatan tishlashini
ko'rsatish uchun.

---

## 5. Qamrovdan tashqarida

- **Telefon enumeratsiyasi** (`/auth-v2/check-phone-number`) — alohida sinf.
- **IP bo'yicha rate limit** — sessionID/telefon kaliti yetarli va aniqroq; IP
  qatlami NAT muammolari va IP almashtirish tufayli qo'shimcha qiymat kam.
- **OAuth (ScienceID)** — parol tekshiruvi tashqi tomonda.
- **Redis integratsiya testlari** — infratuzilma ishi; `AuthThrottle` fake bilan
  to'liq qamralgan.
- **CAPTCHA / progressive delay** — hozir oddiy hard limit; kelajakda kerak bo'lsa.

---

## 6. Shox va zanjir

```
develop ← PR #18  (feature/security-hardening)  CWE-613 + CWE-639
        ← PR #21  (feature/upload-hardening)    CWE-434 + CWE-79
        ← PR #22  (feature/sql-injection)       CWE-89
        ← PR #24  (feature/sensitive-data)      CWE-200
        ← yangi   (feature/brute-force)         CWE-307
```

Merge tartibi majburiy: #18 → #21 → #22 → #24 → yangi. `feature/sensitive-data`
ustiga quriladi, chunki `test/` katalogi shu zanjirda keladi.
