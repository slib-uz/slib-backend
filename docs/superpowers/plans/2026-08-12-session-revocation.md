# Sessiya bekor qilish mexanizmi — amalga oshirish rejasi

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Foydalanuvchi tizimdan chiqqanda uning access va refresh tokenlari darhol kuchdan qolsin (CWE-613).

**Architecture:** JWT stateless bo'lib qoladi, lekin har bir token o'z `jti` si bo'yicha Redis'dagi denylist'ga qarshi tekshiriladi. Yangi `TokenRevocationService` bekor qilish mantig'ini o'z ichiga oladi va `CacheProvider` porti orqali Redis'ga boradi. Tekshiruv `UserAuthTokenService.VerifyToken` ichida — bu yagona nazorat nuqtasi bo'lib, uni ham JWT middleware (access uchun), ham refresh usecase chaqiradi.

**Tech Stack:** Go 1.25, Echo v4, google/wire (DI), Redis (`go-redis/v9`), `golang-jwt/v5` (ES256), standart `testing` paketi.

**Spec:** [2026-08-12-session-revocation-design.md](../specs/2026-08-12-session-revocation-design.md)

## Global Constraints

- Modul nomi `slib.uz`; barcha import yo'llari shundan boshlanadi.
- Core qatlamida (`src/core/**`) loglash uchun **faqat** `github.com/labstack/gommon/log` ishlatiladi. `zap` faqat `src/infrastructure/logger` ichida qoladi.
- Core qatlami infratuzilmaga to'g'ridan-to'g'ri bog'lanmaydi — `src/core/domain/ports/**` portlari orqali ishlaydi.
- `@inject` izohi bo'lgan yangi konstruktor qo'shilgach **`make wire-build`** ishga tushiriladi, aks holda DI konteyner uni ko'rmaydi.
- Swagger izohi o'zgargan yoki yangi route qo'shilgan bo'lsa **`make generate-docs`** ishga tushiriladi.
- `env/.env` git'da kuzatilmaydi (`.gitignore`da `env/*`). Shuning uchun **yangi env o'zgaruvchilari hech qachon `required` bo'lmasin** — har doim `envDefault` ishlatiladi, aks holda mavjud muhitlar ishga tushmay qoladi.
- Yangi bog'liqlik (test kutubxonasi va boshqa) **qo'shilmaydi**. Testlar standart `testing` paketi va qo'lda yozilgan soxta obyektlar bilan yoziladi.
- Soxta obyektlarda katta interfeyslar uchun interfeysning o'zi struct ichiga joylanadi (embedded) — shunda faqat kerakli metodlar qayta yoziladi, qolganlari chaqirilsa panic beradi.
- Har bir task oxirida `go build ./...` va `go test ./...` toza o'tishi shart.

## File Structure

**Yaratiladigan fayllar:**

| Fayl | Mas'uliyati |
|---|---|
| `src/core/application/service/token_revocation_service.go` | Denylist'ga yozish va o'qish; fail-open siyosati |
| `src/core/application/service/token_revocation_service_test.go` | TTL hisobi, fail-open, grace testlari |
| `src/core/application/service/user_auth_token_service_test.go` | VerifyToken denylist tekshiruvi |
| `src/core/application/service/fake_cache_test.go` | Umumiy soxta `CacheProvider` |
| `src/core/application/usecase/authusecases/logout_usecase.go` | Access + refresh bekor qilish, egalik tekshiruvi |
| `src/core/application/usecase/authusecases/logout_usecase_test.go` | Egalik va xato holatlari |
| `src/core/application/usecase/authusecases/refresh_token_usecase_test.go` | Rotatsiyaning ikkala bosqichi va grace oynasi |
| `src/entrypoint/presentation/handlers/auth/logout_handler.go` | HTTP qatlami |
| `src/entrypoint/presentation/handlers/auth/schema/auth_schema.go` | Logout va refresh request tuzilmalari |
| `src/infrastructure/security/jwt_token_service_test.go` | jti saqlanishi |

**O'zgartiriladigan fayllar:**

| Fayl | O'zgarish |
|---|---|
| `src/core/domain/entity/token_entity.go` | `ID string` maydoni |
| `src/infrastructure/security/jwt_token_service.go` | `Decode` `claims.ID` ni uzatadi |
| `src/core/application/service/user_auth_token_service.go` | Denylist tekshiruvi, `ctx` va `grace` parametrlari |
| `src/entrypoint/presentation/app/context/context_wrap.go` | `TokenID`, `TokenExp` maydonlari |
| `src/entrypoint/presentation/interceptor/middlewares/jwt_auth_middleware.go` | Kontekstni to'ldiradi |
| `src/core/application/usecase/authusecases/refresh_token_usecase.go` | Rotatsiya |
| `src/entrypoint/presentation/handlers/auth/refresh_token_handler.go` | Body'dan o'qish, yangi refresh qaytarish |
| `src/entrypoint/presentation/groups/auth_group.go` | `POST /logout` |
| `src/entrypoint/presentation/groups/authv2_group.go` | Sandbox route shartli |
| `src/infrastructure/config/env.go` | 3 ta yangi o'zgaruvchi |
| `src/infrastructure/config/adapter.go` | 3 ta yangi getter |
| `src/core/domain/ports/conf/config_provider.go` | 3 ta yangi metod |
| `Makefile` | `test` target |

---

## Task 1: `TokenEntity` `jti` ni tashimasin deb yo'qotmasin

Hozir `Encode` har bir tokenga uuid `jti` yozadi, lekin `Decode` uni tashlab yuboradi. Denylist `jti` ga tayanadi, shuning uchun bu birinchi qadam.

**Files:**
- Modify: `src/core/domain/entity/token_entity.go`
- Modify: `src/infrastructure/security/jwt_token_service.go:78-82`
- Test: `src/infrastructure/security/jwt_token_service_test.go`

**Interfaces:**
- Consumes: hech narsa (birinchi task)
- Produces: `entity.TokenEntity.ID string` — dekodlangan tokenning `jti` si. Keyingi barcha tasklar shunga tayanadi. `entity.NewTokenEntity(exp time.Time, subject string, payload map[string]any) *TokenEntity` imzosi **o'zgarmaydi**; `ID` konstruktordan keyin alohida o'rnatiladi.

- [ ] **Step 1: Testni yoz**

`src/infrastructure/security/jwt_token_service_test.go`:

```go
package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"slib.uz/src/core/domain/entity"
)

// newTestService haqiqiy ECDSA kalit jufti bilan servis yasaydi.
// Fayldan o'qimaydi — test bir xil paketda bo'lgani uchun maydonlarga to'g'ridan-to'g'ri kirish mumkin.
func newTestService(t *testing.T) *JwtTokenService {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("test kaliti yaratilmadi: %v", err)
	}
	return &JwtTokenService{privateKey: key, publicKey: &key.PublicKey}
}

func newTestClaim() *entity.TokenEntity {
	return entity.NewTokenEntity(
		time.Now().Add(time.Minute),
		"42",
		map[string]any{"type": "ACCESS"},
	)
}

func TestDecodePreservesJTI(t *testing.T) {
	s := newTestService(t)

	signed := s.Encode(newTestClaim())
	if signed == "" {
		t.Fatal("Encode bo'sh token qaytardi")
	}

	decoded, err := s.Decode(signed)
	if err != nil {
		t.Fatalf("Decode xato qaytardi: %v", err)
	}
	if decoded.ID == "" {
		t.Fatal("Decode jti ni yo'qotdi: TokenEntity.ID bo'sh")
	}
}

func TestEachTokenGetsDistinctJTI(t *testing.T) {
	s := newTestService(t)

	first, err := s.Decode(s.Encode(newTestClaim()))
	if err != nil {
		t.Fatalf("birinchi Decode xato: %v", err)
	}
	second, err := s.Decode(s.Encode(newTestClaim()))
	if err != nil {
		t.Fatalf("ikkinchi Decode xato: %v", err)
	}

	if first.ID == second.ID {
		t.Fatalf("ikki xil token bir xil jti oldi: %s", first.ID)
	}
}
```

- [ ] **Step 2: Test yiqilishini tasdiqla**

Run: `go test ./src/infrastructure/security/ -run TestDecodePreservesJTI -v`
Expected: FAIL — kompilyatsiya xatosi `decoded.ID undefined (type *entity.TokenEntity has no field or method ID)`

- [ ] **Step 3: `TokenEntity` ga `ID` qo'sh**

`src/core/domain/entity/token_entity.go` — to'liq yangi mazmuni:

```go
package entity

import "time"

type TokenEntity struct {
	// ID — JWT "jti" da'vosi. Encode paytida generatsiya qilinadi,
	// Decode paytida to'ldiriladi. Bekor qilish shunga tayanadi.
	ID      string
	Exp     time.Time
	Subject string
	Payload map[string]any
}

func NewTokenEntity(exp time.Time, subject string, payload map[string]any) *TokenEntity {
	return &TokenEntity{Exp: exp, Subject: subject, Payload: payload}
}
```

- [ ] **Step 4: `Decode` `jti` ni uzatsin**

`src/infrastructure/security/jwt_token_service.go` — 78-82 qatorlardagi `return` blokini almashtir:

```go
	decoded := entity.NewTokenEntity(
		claims.ExpiresAt.Time,
		claims.Subject,
		claims.Payload,
	)
	decoded.ID = claims.ID

	return decoded, nil
```

- [ ] **Step 5: Testlar o'tishini tasdiqla**

Run: `go test ./src/infrastructure/security/ -v`
Expected: PASS — `TestDecodePreservesJTI` va `TestEachTokenGetsDistinctJTI`

- [ ] **Step 6: Butun loyiha qurilishini tekshir**

Run: `go build ./...`
Expected: xatosiz (exit 0)

- [ ] **Step 7: Commit**

```bash
git add src/core/domain/entity/token_entity.go src/infrastructure/security/jwt_token_service.go src/infrastructure/security/jwt_token_service_test.go
git commit -m "feat(auth): TokenEntity jti ni saqlaydi

Decode ilgari claims.ID ni tashlab yuborardi. Denylist jti ga
tayangani uchun u endi TokenEntity.ID orqali uzatiladi."
```

---

## Task 2: `TokenRevocationService`

Denylist mantig'i — Redis kaliti, TTL hisobi va fail-open siyosati shu servisda jamlanadi.

**Files:**
- Create: `src/core/application/service/token_revocation_service.go`
- Create: `src/core/application/service/fake_cache_test.go`
- Create: `src/core/application/service/token_revocation_service_test.go`
- Modify: `Makefile`

**Interfaces:**
- Consumes: `entity.TokenEntity.ID` (Task 1); mavjud port `cache.CacheProvider` — `GetByKey(ctx context.Context, key string) (string, error)`, `Set(ctx context.Context, key string, value string, expSeconds time.Duration) error`
- Produces:
  - `service.NewTokenRevocationService(c cache.CacheProvider) *TokenRevocationService`
  - `(*TokenRevocationService).Revoke(ctx context.Context, jti string, exp time.Time) error`
  - `(*TokenRevocationService).RevokedAt(ctx context.Context, jti string) (*time.Time, error)` — `(nil, nil)` bekor qilinmaganini bildiradi
  - Test yordamchisi `newFakeCache() *fakeCache` — 3, 4 va 5-tasklarda qayta ishlatiladi

- [ ] **Step 1: Soxta `CacheProvider` ni yoz**

`src/core/application/service/fake_cache_test.go`:

```go
package service

import (
	"context"
	"time"
)

// fakeCache — CacheProvider portining xotiradagi soxta amalga oshirilishi.
// getErr / setErr o'rnatilsa, Redis nosozligini taqlid qiladi.
type fakeCache struct {
	data   map[string]string
	ttls   map[string]time.Duration
	getErr error
	setErr error
}

func newFakeCache() *fakeCache {
	return &fakeCache{
		data: make(map[string]string),
		ttls: make(map[string]time.Duration),
	}
}

func (f *fakeCache) GetByKey(ctx context.Context, key string) (string, error) {
	if f.getErr != nil {
		return "", f.getErr
	}
	return f.data[key], nil
}

func (f *fakeCache) Set(ctx context.Context, key string, value string, expSeconds time.Duration) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.data[key] = value
	f.ttls[key] = expSeconds
	return nil
}
```

- [ ] **Step 2: Servis testlarini yoz**

`src/core/application/service/token_revocation_service_test.go`:

```go
package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRevokeStoresEntryWithTTLFromExp(t *testing.T) {
	c := newFakeCache()
	s := NewTokenRevocationService(c)

	exp := time.Now().Add(10 * time.Minute)
	if err := s.Revoke(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	ttl, ok := c.ttls["revoked:jti:jti-1"]
	if !ok {
		t.Fatal("denylist yozuvi yaratilmadi")
	}
	// TTL exp dan hisoblanadi, qat'iy raqamdan emas.
	if ttl > 10*time.Minute || ttl < 9*time.Minute {
		t.Fatalf("TTL exp dan hisoblanmadi: %v", ttl)
	}
}

func TestRevokeIgnoresAlreadyExpiredToken(t *testing.T) {
	c := newFakeCache()
	s := NewTokenRevocationService(c)

	err := s.Revoke(context.Background(), "jti-old", time.Now().Add(-time.Minute))
	if err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}
	if len(c.data) != 0 {
		t.Fatal("muddati o'tgan token uchun yozuv yaratildi")
	}
}

func TestRevokeReturnsErrorWhenCacheWriteFails(t *testing.T) {
	c := newFakeCache()
	c.setErr = errors.New("redis down")
	s := NewTokenRevocationService(c)

	err := s.Revoke(context.Background(), "jti-1", time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("yozish nosozligi yashirildi — logout jimgina muvaffaqiyatsiz bo'lardi")
	}
}

func TestRevokedAtReturnsNilForUnknownToken(t *testing.T) {
	s := NewTokenRevocationService(newFakeCache())

	at, err := s.RevokedAt(context.Background(), "jti-yoq")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if at != nil {
		t.Fatal("bekor qilinmagan token bekor qilingan deb qaytdi")
	}
}

func TestRevokedAtReturnsRevocationTime(t *testing.T) {
	c := newFakeCache()
	s := NewTokenRevocationService(c)

	if err := s.Revoke(context.Background(), "jti-1", time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	at, err := s.RevokedAt(context.Background(), "jti-1")
	if err != nil {
		t.Fatalf("RevokedAt xato qaytardi: %v", err)
	}
	if at == nil {
		t.Fatal("bekor qilingan token topilmadi")
	}
	if time.Since(*at) > time.Minute {
		t.Fatalf("bekor qilingan vaqt noto'g'ri: %v", *at)
	}
}

func TestRevokedAtFailsOpenWhenCacheUnavailable(t *testing.T) {
	c := newFakeCache()
	c.getErr = errors.New("redis down")
	s := NewTokenRevocationService(c)

	at, err := s.RevokedAt(context.Background(), "jti-1")
	if err != nil {
		t.Fatalf("fail-open buzildi: xato qaytdi: %v", err)
	}
	if at != nil {
		t.Fatal("fail-open buzildi: token bekor qilingan deb qaytdi")
	}
}
```

- [ ] **Step 3: Testlar yiqilishini tasdiqla**

Run: `go test ./src/core/application/service/ -v`
Expected: FAIL — kompilyatsiya xatosi `undefined: NewTokenRevocationService`

- [ ] **Step 4: Servisni yoz**

`src/core/application/service/token_revocation_service.go`:

```go
package service

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/domain/ports/cache"
)

// revokedKeyPrefix — denylist kalitlari uchun Redis prefiksi.
const revokedKeyPrefix = "revoked:jti:"

// TokenRevocationService bekor qilingan tokenlarning jti ro'yxatini boshqaradi.
//
// Kalit qiymati sifatida bekor qilingan vaqt (unix soniya) saqlanadi.
// Bu refresh rotatsiyasidagi "grace" oynasini qo'shimcha kalitsiz beradi.
type TokenRevocationService struct {
	cache cache.CacheProvider
}

// @inject
func NewTokenRevocationService(c cache.CacheProvider) *TokenRevocationService {
	return &TokenRevocationService{cache: c}
}

// Revoke tokenni denylist'ga yozadi. TTL tokenning qolgan umriga teng:
// token o'z-o'zidan o'lgach yozuvni saqlashning ma'nosi yo'q.
// Muddati o'tgan token uchun hech narsa qilmaydi.
func (this *TokenRevocationService) Revoke(ctx context.Context, jti string, exp time.Time) error {
	if jti == "" {
		return nil
	}

	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil
	}

	revokedAt := strconv.FormatInt(time.Now().Unix(), 10)

	return this.cache.Set(ctx, revokedKeyPrefix+jti, revokedAt, ttl)
}

// RevokedAt tokenning bekor qilingan vaqtini qaytaradi.
// (nil, nil) — token bekor qilinmagan.
//
// Redis yetib bo'lmasa ham (nil, nil) qaytaradi va ERROR loglaydi: bu ongli
// fail-open qarori — Redis yiqilishi butun saytni to'xtatmasligi kerak.
// Xavf oynasi access token TTL si bilan chegaralanadi.
func (this *TokenRevocationService) RevokedAt(ctx context.Context, jti string) (*time.Time, error) {
	if jti == "" {
		return nil, nil
	}

	raw, err := this.cache.GetByKey(ctx, revokedKeyPrefix+jti)
	if err != nil {
		log.Error("TokenRevocationService.RevokedAt: denylist o'qib bo'lmadi, fail-open. jti=", jti, " err=", err.Error())
		return nil, nil
	}

	if raw == "" {
		return nil, nil
	}

	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		// Buzuq qiymat — xavfsiz tomonga og'amiz va tokenni bekor qilingan deb hisoblaymiz.
		log.Error("TokenRevocationService.RevokedAt: yaroqsiz qiymat, token bekor qilingan deb hisoblanadi. jti=", jti, " value=", raw)
		zero := time.Unix(0, 0)
		return &zero, nil
	}

	revokedAt := time.Unix(seconds, 0)

	return &revokedAt, nil
}
```

- [ ] **Step 5: Testlar o'tishini tasdiqla**

Run: `go test ./src/core/application/service/ -v`
Expected: PASS — 6 ta test

- [ ] **Step 6: `Makefile` ga `test` target qo'sh**

`Makefile` oxiriga qo'sh:

```makefile
test:
	@echo "Running tests..."
	@go test ./... -count=1
```

- [ ] **Step 7: DI konteynerini qayta generatsiya qil**

Run: `make wire-build`
Expected: xatosiz tugaydi; `cmd/container/container.go` yangilanadi

Run: `go build ./...`
Expected: xatosiz

- [ ] **Step 8: Commit**

```bash
git add src/core/application/service/token_revocation_service.go src/core/application/service/token_revocation_service_test.go src/core/application/service/fake_cache_test.go Makefile cmd/container/container.go
git commit -m "feat(auth): TokenRevocationService qo'shildi

jti asosidagi denylist. TTL tokenning qolgan umridan hisoblanadi.
O'qishda fail-open: Redis yiqilsa sayt ishlashda davom etadi."
```

---

## Task 3: `VerifyToken` denylist'ni tekshirsin, kontekst token identifikatorini bilsin

Yagona nazorat nuqtasi. Bu metodni ham JWT middleware, ham refresh usecase chaqiradi — shuning uchun bitta o'zgarish ikkala yo'lni yopadi. Ayni paytda middleware kontekstga `TokenID`/`TokenExp` yozadi, keyingi taskdagi logout shularga tayanadi.

**Files:**
- Modify: `src/core/application/service/user_auth_token_service.go`
- Modify: `src/entrypoint/presentation/app/context/context_wrap.go:9-14`
- Modify: `src/entrypoint/presentation/interceptor/middlewares/jwt_auth_middleware.go`
- Modify: `src/core/application/usecase/authusecases/refresh_token_usecase.go`
- Test: `src/core/application/service/user_auth_token_service_test.go`

**Interfaces:**
- Consumes: `TokenRevocationService.RevokedAt` (Task 2), `entity.TokenEntity.ID` (Task 1), `newFakeCache()` (Task 2)
- Produces:
  - `(*UserAuthTokenService).VerifyToken(ctx context.Context, tokenString string, tokenType enum.TokenType, grace time.Duration) (*entity.UserBasicEntity, *entity.TokenEntity, error)` — **imzo o'zgardi**: `ctx` va `grace` qo'shildi, dekodlangan token ham qaytariladi
  - `context.Context` struct maydonlari: `TokenID string`, `TokenExp time.Time`

- [ ] **Step 1: Testni yoz**

`src/core/application/service/user_auth_token_service_test.go`:

```go
package service

import (
	"context"
	"testing"
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeTokenService — TokenService portining soxta amalga oshirilishi.
type fakeTokenService struct {
	token *entity.TokenEntity
	err   error
}

func (f *fakeTokenService) Encode(claim *entity.TokenEntity) string { return "signed" }

func (f *fakeTokenService) Decode(token string) (*entity.TokenEntity, error) {
	return f.token, f.err
}

// fakeUserRepo — UserRepository interfeysi struct ichiga joylangan (embedded),
// shuning uchun faqat GetById qayta yoziladi; qolgan 15 metod chaqirilsa panic beradi.
type fakeUserRepo struct {
	repository.UserRepository
	user *entity.UserEntity
}

func (f *fakeUserRepo) GetById(id uint) (*entity.UserEntity, error) {
	return f.user, nil
}

func newAuthService(revocation *TokenRevocationService, token *entity.TokenEntity) *UserAuthTokenService {
	return &UserAuthTokenService{
		tokenService:   &fakeTokenService{token: token},
		userRepository: &fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		revocation:     revocation,
	}
}

func accessToken(jti string, exp time.Time) *entity.TokenEntity {
	t := entity.NewTokenEntity(exp, "42", map[string]any{"type": string(enum.TokenTypeAccess)})
	t.ID = jti
	return t
}

func TestVerifyTokenAcceptsLiveToken(t *testing.T) {
	revocation := NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)
	svc := newAuthService(revocation, accessToken("jti-1", exp))

	user, token, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0)
	if err != nil {
		t.Fatalf("yaroqli token rad etildi: %v", err)
	}
	if user == nil || user.ID != 42 {
		t.Fatal("foydalanuvchi qaytarilmadi")
	}
	if token == nil || token.ID != "jti-1" {
		t.Fatal("dekodlangan token qaytarilmadi")
	}
}

func TestVerifyTokenRejectsRevokedToken(t *testing.T) {
	revocation := NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)

	if err := revocation.Revoke(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	svc := newAuthService(revocation, accessToken("jti-1", exp))

	_, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0)
	if err == nil {
		t.Fatal("bekor qilingan token qabul qilindi")
	}
}

func TestVerifyTokenAcceptsRevokedTokenInsideGrace(t *testing.T) {
	revocation := NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)

	if err := revocation.Revoke(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	svc := newAuthService(revocation, accessToken("jti-1", exp))

	// Endigina bekor qilindi, grace 60 soniya — hali qabul qilinishi kerak.
	_, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, time.Minute)
	if err != nil {
		t.Fatalf("grace oynasi ichidagi token rad etildi: %v", err)
	}
}

func TestVerifyTokenRejectsWrongTokenType(t *testing.T) {
	revocation := NewTokenRevocationService(newFakeCache())
	svc := newAuthService(revocation, accessToken("jti-1", time.Now().Add(time.Minute)))

	_, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeRefresh, 0)
	if err == nil {
		t.Fatal("access token refresh sifatida qabul qilindi")
	}
}

// countingUserRepo GetById necha marta chaqirilganini sanaydi.
type countingUserRepo struct {
	repository.UserRepository
	user  *entity.UserEntity
	calls int
}

func (f *countingUserRepo) GetById(id uint) (*entity.UserEntity, error) {
	f.calls++
	return f.user, nil
}

// Denylist tekshiruvi DB'dan oldin turishi kerak: bekor qilingan token uchun
// ortiqcha DB so'rovi ketmasligi lozim.
func TestVerifyTokenSkipsDatabaseWhenTokenRevoked(t *testing.T) {
	revocation := NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)

	if err := revocation.Revoke(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	repo := &countingUserRepo{user: &entity.UserEntity{ID: 42}}
	svc := &UserAuthTokenService{
		tokenService:   &fakeTokenService{token: accessToken("jti-1", exp)},
		userRepository: repo,
		revocation:     revocation,
	}

	if _, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0); err == nil {
		t.Fatal("bekor qilingan token qabul qilindi")
	}
	if repo.calls != 0 {
		t.Fatalf("bekor qilingan token uchun DB'ga %d marta borildi, 0 kutilgandi", repo.calls)
	}
}
```

- [ ] **Step 2: Test yiqilishini tasdiqla**

Run: `go test ./src/core/application/service/ -run TestVerifyToken -v`
Expected: FAIL — kompilyatsiya xatosi: `too many arguments in call to svc.VerifyToken` va `unknown field revocation`

- [ ] **Step 3: `UserAuthTokenService` ni yangila**

`src/core/application/service/user_auth_token_service.go` — import blokiga `"context"` qo'sh, struct, konstruktor va `VerifyToken` ni almashtir:

```go
type UserAuthTokenService struct {
	tokenService   security.TokenService
	userRepository repository.UserRepository
	cfg            conf.ConfigAdapter
	revocation     *TokenRevocationService
}

// @inject
func NewUserAuthTokenService(
	tokenService security.TokenService,
	userRepository repository.UserRepository,
	cfg conf.ConfigAdapter,
	revocation *TokenRevocationService,
) *UserAuthTokenService {
	return &UserAuthTokenService{
		tokenService:   tokenService,
		userRepository: userRepository,
		cfg:            cfg,
		revocation:     revocation,
	}
}

// VerifyToken tokenni tekshiradi va foydalanuvchi bilan birga dekodlangan tokenni qaytaradi.
//
// grace > 0 bo'lsa, shu muddat ichida bekor qilingan token hali ham qabul qilinadi.
// Bu refresh rotatsiyasida parallel so'rovlar uzilmasligi uchun kerak.
// Access yo'lida grace har doim 0 bo'ladi.
func (this *UserAuthTokenService) VerifyToken(
	ctx context.Context,
	tokenString string,
	tokenType enum.TokenType,
	grace time.Duration,
) (*entity.UserBasicEntity, *entity.TokenEntity, error) {

	token, errDecode := this.tokenService.Decode(tokenString)
	if errDecode != nil {
		return nil, nil, errDecode
	}

	if token.Payload["type"] != string(tokenType) {
		return nil, nil, response.InvalidTokenError
	}

	revokedAt, err := this.revocation.RevokedAt(ctx, token.ID)
	if err != nil {
		return nil, nil, err
	}
	if revokedAt != nil && time.Since(*revokedAt) > grace {
		log.Warn("UserAuthTokenService.VerifyToken: bekor qilingan token ishlatildi. jti=", token.ID, " subject=", token.Subject)
		return nil, nil, response.InvalidTokenError
	}

	userID, err := strconv.Atoi(token.Subject)
	if err != nil {
		return nil, nil, response.InvalidTokenError
	}

	user, err := this.userRepository.GetById(uint(userID))
	if err != nil {
		return nil, nil, err
	}

	return mapper.UserEntityToBasic(user), token, nil
}
```

Import blokiga `"github.com/labstack/gommon/log"` ni ham qo'sh.

- [ ] **Step 4: Testlar o'tishini tasdiqla**

Run: `go test ./src/core/application/service/ -v`
Expected: PASS — Task 2 dagi 6 ta va bu taskdagi 5 ta test

- [ ] **Step 5: Kontekstga token maydonlarini qo'sh**

`src/entrypoint/presentation/app/context/context_wrap.go` — `Context` structini almashtir va `"time"` importini qo'sh:

```go
type Context struct {
	echo.Context
	User        *entity.UserBasicEntity
	Client      *entity.ClientEntity
	AnonymousID string

	// TokenID va TokenExp joriy access tokenning jti si va muddati.
	// JWT middleware to'ldiradi; logout ularsiz tokenni bekor qila olmaydi.
	TokenID  string
	TokenExp time.Time
}
```

- [ ] **Step 6: Middleware'ni yangila**

`src/entrypoint/presentation/interceptor/middlewares/jwt_auth_middleware.go` — `userAuth` metodini almashtir:

```go
func (this *JwAuthMiddleware) userAuth(c *context.Context) error {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)

	if authHeader == "" || !this.isBearer(authHeader) {
		return nil
	}
	tokenStr := this.getBearerToken(authHeader)

	// Access yo'lida grace yo'q: bekor qilingan token darhol rad etiladi.
	userDto, token, err := this.authService.VerifyToken(c.Request().Context(), tokenStr, enum.TokenTypeAccess, 0)
	if err != nil {
		return err
	}

	c.User = userDto
	c.TokenID = token.ID
	c.TokenExp = token.Exp

	this.updateLastOnlineAsync(userDto.ID)
	return nil
}
```

- [ ] **Step 7: Refresh usecase'ni yangi imzoga moslashtir**

`src/core/application/usecase/authusecases/refresh_token_usecase.go` — `Execute` ni almashtir (rotatsiya keyingi taskda qo'shiladi, hozir faqat kompilyatsiya uchun):

```go
func (this *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (string, error) {
	user, _, err := this.service.VerifyToken(ctx, refreshToken, enum.TokenTypeRefresh, 0)
	if err != nil {
		return "", err
	}
	return this.service.GenerateAccessToken(user.ID), nil
}
```

Import blokiga `"context"` qo'sh.

- [ ] **Step 8: Refresh handler'ni moslashtir**

`src/entrypoint/presentation/handlers/auth/refresh_token_handler.go` — `Handle` ichidagi chaqiruvni almashtir:

```go
	accessToken, err := this.uc.Execute(ctx.Request().Context(), refreshToken)
```

- [ ] **Step 9: Qurilish va testlarni tekshir**

Run: `make wire-build && go build ./... && go test ./... -count=1`
Expected: hammasi xatosiz

- [ ] **Step 10: Commit**

```bash
git add src/core/application/service/ src/entrypoint/presentation/app/context/context_wrap.go src/entrypoint/presentation/interceptor/middlewares/jwt_auth_middleware.go src/core/application/usecase/authusecases/refresh_token_usecase.go src/entrypoint/presentation/handlers/auth/refresh_token_handler.go cmd/container/container.go
git commit -m "feat(auth): VerifyToken denylist'ni tekshiradi

Access va refresh yo'llari yagona nazorat nuqtasidan o'tadi.
Middleware kontekstga jti va exp yozadi — logout shularga tayanadi."
```

---

## Task 4: `POST /auth/logout`

**Files:**
- Create: `src/core/application/usecase/authusecases/logout_usecase.go`
- Create: `src/core/application/usecase/authusecases/logout_usecase_test.go`
- Create: `src/entrypoint/presentation/handlers/auth/schema/auth_schema.go`
- Create: `src/entrypoint/presentation/handlers/auth/logout_handler.go`
- Modify: `src/entrypoint/presentation/groups/auth_group.go`

**Interfaces:**
- Consumes: `TokenRevocationService.Revoke` (Task 2), `context.Context.TokenID` / `.TokenExp` (Task 3), `security.TokenService.Decode` (Task 1)
- Produces:
  - `authusecases.NewLogoutUseCase(revocation *service.TokenRevocationService, tokenService security.TokenService) *LogoutUseCase`
  - `(*LogoutUseCase).Execute(ctx context.Context, userID uint, accessJTI string, accessExp time.Time, refreshToken string) error`
  - `schema.LogoutRequest{ RefreshToken string }`

- [ ] **Step 1: Usecase testini yoz**

`src/core/application/usecase/authusecases/logout_usecase_test.go`:

```go
package authusecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type fakeCache struct {
	data   map[string]string
	setErr error
}

func newFakeCache() *fakeCache {
	return &fakeCache{data: make(map[string]string)}
}

func (f *fakeCache) GetByKey(ctx context.Context, key string) (string, error) {
	return f.data[key], nil
}

func (f *fakeCache) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.data[key] = value
	return nil
}

type fakeTokenService struct {
	token *entity.TokenEntity
	err   error
}

func (f *fakeTokenService) Encode(claim *entity.TokenEntity) string { return "signed" }

func (f *fakeTokenService) Decode(token string) (*entity.TokenEntity, error) {
	return f.token, f.err
}

func refreshTokenOf(subject string, jti string) *entity.TokenEntity {
	t := entity.NewTokenEntity(time.Now().Add(time.Hour), subject, map[string]any{
		"type": string(enum.TokenTypeRefresh),
	})
	t.ID = jti
	return t
}

func TestLogoutRevokesAccessAndRefresh(t *testing.T) {
	c := newFakeCache()
	uc := NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{token: refreshTokenOf("42", "refresh-jti")})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "raw-refresh")
	if err != nil {
		t.Fatalf("logout xato qaytardi: %v", err)
	}

	if _, ok := c.data["revoked:jti:access-jti"]; !ok {
		t.Error("access token bekor qilinmadi")
	}
	if _, ok := c.data["revoked:jti:refresh-jti"]; !ok {
		t.Error("refresh token bekor qilinmadi")
	}
}

func TestLogoutRejectsRefreshTokenOfAnotherUser(t *testing.T) {
	c := newFakeCache()
	// Refresh token 99-foydalanuvchiniki, lekin so'rovni 42 yubormoqda.
	uc := NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{token: refreshTokenOf("99", "refresh-jti")})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "raw-refresh")
	if err == nil {
		t.Fatal("begona refresh token qabul qilindi — bekor qilish orqali DoS mumkin")
	}

	if _, ok := c.data["revoked:jti:refresh-jti"]; ok {
		t.Error("begona refresh token bekor qilindi")
	}
}

func TestLogoutSucceedsWithoutRefreshToken(t *testing.T) {
	c := newFakeCache()
	uc := NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "")
	if err != nil {
		t.Fatalf("refreshsiz logout xato qaytardi: %v", err)
	}
	if _, ok := c.data["revoked:jti:access-jti"]; !ok {
		t.Error("access token bekor qilinmadi")
	}
}

func TestLogoutSucceedsWhenRefreshTokenIsInvalid(t *testing.T) {
	c := newFakeCache()
	uc := NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{err: errors.New("invalid token")})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "buzuq")
	if err != nil {
		t.Fatalf("yaroqsiz refresh logoutni yiqitdi: %v", err)
	}
}

func TestLogoutFailsWhenCacheWriteFails(t *testing.T) {
	c := newFakeCache()
	c.setErr = errors.New("redis down")
	uc := NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "")
	if err == nil {
		t.Fatal("Redis nosozligi yashirildi — foydalanuvchi chiqdim deb o'ylaydi, token esa ishlaydi")
	}
}
```

- [ ] **Step 2: Test yiqilishini tasdiqla**

Run: `go test ./src/core/application/usecase/authusecases/ -v`
Expected: FAIL — `undefined: NewLogoutUseCase`

- [ ] **Step 3: Usecase'ni yoz**

`src/core/application/usecase/authusecases/logout_usecase.go`:

```go
package authusecases

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/security"
)

type LogoutUseCase struct {
	revocation   *service.TokenRevocationService
	tokenService security.TokenService
}

// @inject
func NewLogoutUseCase(revocation *service.TokenRevocationService, tokenService security.TokenService) *LogoutUseCase {
	return &LogoutUseCase{revocation: revocation, tokenService: tokenService}
}

// Execute joriy access tokenni va berilgan bo'lsa refresh tokenni bekor qiladi.
//
// Faqat access tokenni bekor qilish yetarli emas: refresh qo'lda qolsa,
// mijoz darhol yangi access token oladi va logout ma'nosiz bo'ladi.
func (this *LogoutUseCase) Execute(
	ctx context.Context,
	userID uint,
	accessJTI string,
	accessExp time.Time,
	refreshToken string,
) error {

	if err := this.revocation.Revoke(ctx, accessJTI, accessExp); err != nil {
		log.Error("LogoutUseCase: access token bekor qilinmadi. user_id=", userID, " err=", err.Error())
		return response.NewFailResponse(503, "chiqish amalga oshmadi, qaytadan urinib ko'ring")
	}

	if refreshToken == "" {
		log.Warn("LogoutUseCase: refresh token yuborilmadi, faqat access bekor qilindi. user_id=", userID)
		return nil
	}

	token, err := this.tokenService.Decode(refreshToken)
	if err != nil {
		// Yaroqsiz yoki muddati o'tgan refresh — chiqish baribir muvaffaqiyatli hisoblanadi.
		return nil
	}

	if token.Payload["type"] != string(enum.TokenTypeRefresh) {
		return nil
	}

	// Egalik tekshiruvi: usiz istalgan foydalanuvchi boshqasining refresh
	// tokenini bekor qilib, uni tizimdan uzib qo'ya olardi.
	if token.Subject != strconv.Itoa(int(userID)) {
		log.Warn("LogoutUseCase: begona refresh tokenni bekor qilishga urinish. user_id=", userID, " subject=", token.Subject)
		return response.PermissionDeniedError
	}

	if err := this.revocation.Revoke(ctx, token.ID, token.Exp); err != nil {
		log.Error("LogoutUseCase: refresh token bekor qilinmadi. user_id=", userID, " err=", err.Error())
		return response.NewFailResponse(503, "chiqish amalga oshmadi, qaytadan urinib ko'ring")
	}

	return nil
}
```

- [ ] **Step 4: Testlar o'tishini tasdiqla**

Run: `go test ./src/core/application/usecase/authusecases/ -v`
Expected: PASS — 5 ta test

- [ ] **Step 5: Request schema'sini yoz**

`src/entrypoint/presentation/handlers/auth/schema/auth_schema.go`:

```go
package schema

// LogoutRequest — refresh token ixtiyoriy: bosqichli rollout davrida
// eski frontend uni yubormasligi mumkin.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequest — refresh tokenni body orqali qabul qilish uchun.
// Query parametr 1-bosqichda orqaga moslik uchun qoladi.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
```

- [ ] **Step 6: Handler'ni yoz**

`src/entrypoint/presentation/handlers/auth/logout_handler.go`:

```go
package auth

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/auth/schema"
)

type LogoutHandler struct {
	uc *authusecases.LogoutUseCase
}

// @inject
func NewLogoutHandler(uc *authusecases.LogoutUseCase) *LogoutHandler {
	return &LogoutHandler{uc: uc}
}

// Handle godoc
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data  body      schema.LogoutRequest  false  "Refresh token (ixtiyoriy)"
// @Success      200   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      503   {object}  response.Response
// @Router       /auth/logout [post]
func (this *LogoutHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	// Body bo'sh bo'lishi mumkin — bu xato emas, faqat access bekor qilinadi.
	req, err := context2.GetBody[schema.LogoutRequest](ctx)
	if err != nil {
		req = &schema.LogoutRequest{}
	}

	if err := this.uc.Execute(
		c.Request().Context(),
		c.User.ID,
		c.TokenID,
		c.TokenExp,
		req.RefreshToken,
	); err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]any{"success": true})
}
```

- [ ] **Step 7: Route'ni ro'yxatdan o'tkaz**

`src/entrypoint/presentation/groups/auth_group.go` da uchta o'zgarish:

1. `AuthGroup` structiga maydon qo'sh:
```go
	logoutHandler           *auth.LogoutHandler
```

2. `NewAuthGroup` parametrlariga `logoutHandler *auth.LogoutHandler` qo'sh va qaytariladigan structga `logoutHandler: logoutHandler,` yoz.

3. `RegisterRoutes` ga qator qo'sh:
```go
	group.POST("/logout", this.logoutHandler.Handle, permissions.AuthenticatedPermission)
```

- [ ] **Step 8: Qurilish, DI va hujjatlar**

Run: `make wire-build && go build ./... && go test ./... -count=1 && make generate-docs`
Expected: hammasi xatosiz

- [ ] **Step 9: Commit**

```bash
git add src/core/application/usecase/authusecases/ src/entrypoint/presentation/handlers/auth/ src/entrypoint/presentation/groups/auth_group.go cmd/container/container.go src/entrypoint/presentation/docs/
git commit -m "feat(auth): POST /auth/logout qo'shildi

Access va refresh tokenlar birga bekor qilinadi. Refresh tokenning
egaligi tekshiriladi — usiz begona sessiyani uzish mumkin edi.
Redis'ga yozib bo'lmasa 503 qaytariladi, 200 emas."
```

---

## Task 5: Refresh rotatsiyasi va konfiguratsiya bayroqlari

**Files:**
- Modify: `src/infrastructure/config/env.go`
- Modify: `src/infrastructure/config/adapter.go`
- Modify: `src/core/domain/ports/conf/config_provider.go`
- Modify: `src/core/application/usecase/authusecases/refresh_token_usecase.go`
- Modify: `src/entrypoint/presentation/handlers/auth/refresh_token_handler.go`
- Test: `src/core/application/usecase/authusecases/refresh_token_usecase_test.go`

**Interfaces:**
- Consumes: `VerifyToken(ctx, token, type, grace)` (Task 3), `TokenRevocationService.Revoke` (Task 2), `schema.RefreshTokenRequest` (Task 4)
- Produces:
  - `conf.ConfigAdapter` yangi metodlari: `IsRefreshRotationStrict() bool`, `GetRefreshRotationGraceSeconds() int`, `IsProduction() bool`
  - `authusecases.RefreshTokenResult{ AccessToken string; RefreshToken string }`
  - `(*RefreshTokenUseCase).Execute(ctx context.Context, refreshToken string) (*RefreshTokenResult, error)`

- [ ] **Step 1: Config o'zgaruvchilarini qo'sh**

`src/infrastructure/config/env.go` — `Config` structidagi `// Auth` bo'limiga qo'sh:

```go
	// Auth
	JwtPrivateKeyPath            string `env:"JWT_PRIVATE_KEY_PATH,required"`
	JwtPublicKeyPath             string `env:"JWT_PUBLIC_KEY_PATH,required"`
	JwtAccessTokenExpireMinutes  int    `env:"JWT_ACCESS_TOKEN_EXPIRE_MINUTES,required"`
	JwtRefreshTokenExpireMinutes int    `env:"JWT_REFRESH_TOKEN_EXPIRE_MINUTES,required"`

	// Sessiya. Bu uchtasi required EMAS: env/.env git'da kuzatilmaydi,
	// required qilinsa mavjud muhitlar ishga tushmay qoladi.
	RefreshRotationStrict       bool `env:"REFRESH_ROTATION_STRICT" envDefault:"false"`
	RefreshRotationGraceSeconds int  `env:"REFRESH_ROTATION_GRACE_SECONDS" envDefault:"60"`

	// Production default true — xavfsiz tomonga og'ish.
	// O'rnatilmagan bo'lsa sandbox login o'chiq bo'ladi.
	Production bool `env:"PRODUCTION" envDefault:"true"`
```

- [ ] **Step 2: Portga metodlarni qo'sh**

`src/core/domain/ports/conf/config_provider.go` — interfeysga qo'sh:

```go
	IsRefreshRotationStrict() bool
	GetRefreshRotationGraceSeconds() int
	IsProduction() bool
```

- [ ] **Step 3: Adapterda amalga oshir**

`src/infrastructure/config/adapter.go` oxiriga qo'sh:

```go
func (this *ConfigAdapterImpl) IsRefreshRotationStrict() bool {
	return this.env.RefreshRotationStrict
}

func (this *ConfigAdapterImpl) GetRefreshRotationGraceSeconds() int {
	return this.env.RefreshRotationGraceSeconds
}

func (this *ConfigAdapterImpl) IsProduction() bool {
	return this.env.Production
}
```

- [ ] **Step 4: Qurilishni tekshir**

Run: `go build ./...`
Expected: xatosiz

- [ ] **Step 5: Rotatsiya testlarini yoz**

`src/core/application/usecase/authusecases/refresh_token_usecase_test.go`.
`fakeCache`, `fakeTokenService` va `refreshTokenOf` Task 4 dagi test faylida
allaqachon e'lon qilingan — bir xil paket bo'lgani uchun qayta ishlatiladi.

```go
package authusecases

import (
	"context"
	"testing"

	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeUserRepo — UserRepository interfeysi struct ichiga joylangan (embedded),
// shuning uchun faqat GetById qayta yoziladi.
type fakeUserRepo struct {
	repository.UserRepository
	user *entity.UserEntity
}

func (f *fakeUserRepo) GetById(id uint) (*entity.UserEntity, error) {
	return f.user, nil
}

// fakeConfig — ConfigAdapter interfeysi embedded; faqat kerakli getterlar yoziladi.
type fakeConfig struct {
	conf.ConfigAdapter
	strict bool
	grace  int
}

func (f *fakeConfig) IsRefreshRotationStrict() bool        { return f.strict }
func (f *fakeConfig) GetRefreshRotationGraceSeconds() int  { return f.grace }
func (f *fakeConfig) GetJwtAccessTokenExpireMinutes() int  { return 15 }
func (f *fakeConfig) GetJwtRefreshTokenExpireMinutes() int { return 10080 }

func newRefreshUseCase(c *fakeCache, strict bool, token *entity.TokenEntity) *RefreshTokenUseCase {
	revocation := service.NewTokenRevocationService(c)
	cfg := &fakeConfig{strict: strict, grace: 60}
	auth := service.NewUserAuthTokenService(
		&fakeTokenService{token: token},
		&fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		cfg,
		revocation,
	)
	return NewRefreshTokenUseCase(auth, revocation, cfg)
}

func TestRefreshKeepsOldTokenAliveInPhaseOne(t *testing.T) {
	c := newFakeCache()
	uc := newRefreshUseCase(c, false, refreshTokenOf("42", "old-jti"))

	result, err := uc.Execute(context.Background(), "raw-refresh")
	if err != nil {
		t.Fatalf("refresh xato qaytardi: %v", err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatal("yangi juftlik qaytarilmadi")
	}
	if _, revoked := c.data["revoked:jti:old-jti"]; revoked {
		t.Fatal("1-bosqichda eski refresh bekor qilindi — eski frontend buziladi")
	}
}

func TestRefreshRevokesOldTokenInPhaseTwo(t *testing.T) {
	c := newFakeCache()
	uc := newRefreshUseCase(c, true, refreshTokenOf("42", "old-jti"))

	if _, err := uc.Execute(context.Background(), "raw-refresh"); err != nil {
		t.Fatalf("refresh xato qaytardi: %v", err)
	}
	if _, revoked := c.data["revoked:jti:old-jti"]; !revoked {
		t.Fatal("2-bosqichda eski refresh bekor qilinmadi")
	}
}

func TestRefreshAcceptsRecentlyRevokedTokenInsideGrace(t *testing.T) {
	c := newFakeCache()
	token := refreshTokenOf("42", "old-jti")

	revocation := service.NewTokenRevocationService(c)
	// Token endigina bekor qilindi — grace oynasi ichida hali qabul qilinishi kerak.
	if err := revocation.Revoke(context.Background(), "old-jti", token.Exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	cfg := &fakeConfig{strict: true, grace: 60}
	auth := service.NewUserAuthTokenService(
		&fakeTokenService{token: token},
		&fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		cfg,
		revocation,
	)
	uc := NewRefreshTokenUseCase(auth, revocation, cfg)

	if _, err := uc.Execute(context.Background(), "raw-refresh"); err != nil {
		t.Fatalf("grace oynasi ichidagi refresh rad etildi: %v", err)
	}
}
```

- [ ] **Step 6: Testlar yiqilishini tasdiqla**

Run: `go test ./src/core/application/usecase/authusecases/ -run TestRefresh -v`
Expected: FAIL — kompilyatsiya xatosi: `not enough arguments in call to NewRefreshTokenUseCase` va `result.AccessToken undefined`

- [ ] **Step 7: Rotatsiya usecase'ini yoz**

`src/core/application/usecase/authusecases/refresh_token_usecase.go` — to'liq yangi mazmuni:

```go
package authusecases

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/conf"
)

// RefreshTokenResult — rotatsiyadan keyingi yangi juftlik.
type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenUseCase struct {
	service    *service.UserAuthTokenService
	revocation *service.TokenRevocationService
	cfg        conf.ConfigAdapter
}

// @inject
func NewRefreshTokenUseCase(
	svc *service.UserAuthTokenService,
	revocation *service.TokenRevocationService,
	cfg conf.ConfigAdapter,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{service: svc, revocation: revocation, cfg: cfg}
}

// Execute refresh token asosida yangi access va refresh juftligini qaytaradi.
//
// 1-bosqich (REFRESH_ROTATION_STRICT=false): eski refresh bekor qilinmaydi,
// shuning uchun yangi juftlikni saqlamaydigan eski frontend ishlashda davom etadi.
// Logout baribir ishlaydi, chunki u refresh jti sini bevosita denylist'ga yozadi.
//
// 2-bosqich (true): eski refresh bekor qilinadi, grace oynasi parallel
// so'rovlarni uzilishdan saqlaydi.
func (this *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (*RefreshTokenResult, error) {
	strict := this.cfg.IsRefreshRotationStrict()

	grace := time.Duration(0)
	if strict {
		grace = time.Duration(this.cfg.GetRefreshRotationGraceSeconds()) * time.Second
	}

	user, token, err := this.service.VerifyToken(ctx, refreshToken, enum.TokenTypeRefresh, grace)
	if err != nil {
		return nil, err
	}

	result := &RefreshTokenResult{
		AccessToken:  this.service.GenerateAccessToken(user.ID),
		RefreshToken: this.service.GenerateRefreshToken(user.ID),
	}

	if strict {
		if err := this.revocation.Revoke(ctx, token.ID, token.Exp); err != nil {
			// Yangi juftlik allaqachon berildi — so'rovni yiqitmaymiz, lekin belgilab qo'yamiz.
			log.Error("RefreshTokenUseCase: eski refresh bekor qilinmadi. jti=", token.ID, " err=", err.Error())
		}
	}

	return result, nil
}
```

- [ ] **Step 8: Testlar o'tishini tasdiqla**

Run: `go test ./src/core/application/usecase/authusecases/ -v`
Expected: PASS — Task 4 dagi 5 ta va bu taskdagi 3 ta test

- [ ] **Step 9: Handler'ni yangila**

`src/entrypoint/presentation/handlers/auth/refresh_token_handler.go` — to'liq yangi mazmuni.

**Javob shakli yassi qoladi.** Bu endpoint har doim `{"access_token":...,"refresh_token":...}`
qaytargan va ishlab turgan frontend shunga bog'liq. `c.JsonResponse` javobni
`{"status":...,"data":{...}}` ichiga o'raydi — bu mijoz uchun buzilish bo'lardi, ya'ni aynan
bosqichli rollout oldini olmoqchi bo'lgan narsa. Shuning uchun `ctx.JSON` ishlatiladi.
(Kodbaza bu jihatdan bir xil emas: `/auth/oauth/login` yassi qaytaradi, `/auth-v2/login` esa
o'ralgan. Bu tarixiy holat; refresh yassi bo'lgan va yassi qoladi.)

```go
package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/authusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/auth/schema"
)

type RefreshTokenHandler struct {
	uc *authusecases.RefreshTokenUseCase
}

// @inject
func NewRefreshTokenHandler(usecase *authusecases.RefreshTokenUseCase) *RefreshTokenHandler {
	return &RefreshTokenHandler{uc: usecase}
}

// Handle godoc
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data           body   schema.RefreshTokenRequest  false  "Refresh token"
// @Param        refresh_token  query  string                      false  "Refresh token (eskirgan, body'dan foydalaning)"
// @Success      200  {object}  schema.RefreshTokenResponse
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/refresh-token [post]
func (this *RefreshTokenHandler) Handle(ctx echo.Context) error {
	refreshToken := ""
	if req, err := context2.GetBody[schema.RefreshTokenRequest](ctx); err == nil {
		refreshToken = req.RefreshToken
	}
	// Orqaga moslik: 1-bosqichda eski frontend query parametr yuboradi.
	// 2-bosqichda bu tarmoq olib tashlanadi.
	if refreshToken == "" {
		refreshToken = ctx.QueryParam("refresh_token")
	}
	if refreshToken == "" {
		return response.NewFailResponse(400, "refresh_token talab qilinadi")
	}

	result, err := this.uc.Execute(ctx.Request().Context(), refreshToken)
	if err != nil {
		return err
	}

	// Javob ataylab yassi: o'rash mijoz uchun buzilish bo'lardi.
	return ctx.JSON(http.StatusOK, map[string]any{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}
```

`schema.RefreshTokenResponse` `handlers/auth/schema/auth_schema.go` ga qo'shiladi — u faqat
swagger uchun, javob shaklini hujjatda mahkamlaydi:

```go
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
```

- [ ] **Step 10: Qurilish, DI, testlar va hujjatlar**

Run: `make wire-build && go build ./... && go test ./... -count=1 && make generate-docs`
Expected: hammasi xatosiz

- [ ] **Step 11: Commit**

```bash
git add src/infrastructure/config/ src/core/domain/ports/conf/ src/core/application/usecase/authusecases/ src/entrypoint/presentation/handlers/auth/refresh_token_handler.go cmd/container/container.go src/entrypoint/presentation/docs/
git commit -m "feat(auth): refresh token rotatsiyasi (1-bosqich)

Har bir refresh yangi juftlik qaytaradi. REFRESH_ROTATION_STRICT=false
bo'lgani uchun eski refresh hozircha yashaydi — eski frontend buzilmaydi.
Refresh endi body orqali ham qabul qilinadi."
```

---

## Task 6: Sandbox login prod'da o'chirilsin

`POST /api/auth-v2/sandbox/login` autentifikatsiyasiz, statik OTP bilan haqiqiy sessiya beradi va foydalanuvchi topilmasa yangi hisob yaratadi. Prod'da mavjud bo'lsa, sessiyani bekor qilish ishining ma'nosi qolmaydi.

**Files:**
- Modify: `src/entrypoint/presentation/groups/authv2_group.go`

**Interfaces:**
- Consumes: `conf.ConfigAdapter.IsProduction()` (Task 5)
- Produces: hech narsa (oxirgi task)

- [ ] **Step 1: Guruhga config'ni ulash va route'ni shartli qilish**

`src/entrypoint/presentation/groups/authv2_group.go` — to'liq yangi mazmuni:

```go
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
```

- [ ] **Step 2: Qurilish va DI**

Run: `make wire-build && go build ./... && go test ./... -count=1`
Expected: hammasi xatosiz

- [ ] **Step 3: Prod rejimda route yo'qligini qo'lda tekshir**

Run: `PRODUCTION=true go run ./cmd/http/main.go` (alohida terminalda), so'ng:

```bash
curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/api/auth-v2/sandbox/login \
  -H 'Content-Type: application/json' -d '{"phone":"998900000001","otp":"111111"}'
```

Expected: `404`

So'ng `PRODUCTION=false` bilan qayta ishga tushirib, o'sha so'rov `404` **emas**ligini tasdiqla (400 yoki 200 — route mavjudligini bildiradi).

- [ ] **Step 4: Hujjatlarni yangila va commit**

Run: `make generate-docs`

```bash
git add src/entrypoint/presentation/groups/authv2_group.go cmd/container/container.go src/entrypoint/presentation/docs/
git commit -m "fix(auth): sandbox login prod'da ro'yxatdan o'tmaydi

Endpoint autentifikatsiyasiz, statik OTP bilan haqiqiy sessiya berardi.
PRODUCTION default qiymati true — o'rnatilmagan muhitda ham o'chiq."
```

---

## Deploy (kod emas, operatsion qadamlar)

Bu qadamlar `env/.env` git'da kuzatilmagani uchun rejadagi commitlarga kirmaydi.

Tartib spec 10-bo'limi bilan bir xil: **`.env` merge'dan oldin**. `develop`/`main` ga merge
qilish deployni avtomatik ishga tushiradi (`.github/workflows/ci.yml`), ya'ni "deploy qilish"
alohida boshqariladigan qadam emas — merge bilan birga sodir bo'ladi.

### Merge'dan OLDIN (majburiy)

- [ ] **1. `env/.env` ni har bir serverda yangila** (prod, dev/QA — hammasida) va ilovani qayta
      ishga tushir:

```
JWT_ACCESS_TOKEN_EXPIRE_MINUTES=15
JWT_REFRESH_TOKEN_EXPIRE_MINUTES=10080
REFRESH_ROTATION_STRICT=false
REFRESH_ROTATION_GRACE_SECONDS=60
PRODUCTION=true      # dev/QA da: PRODUCTION=false
```

Merge'dan keyin qilinsa, oradagi vaqtda access TTL 600 daqiqa bo'lib qolaveradi va "fail-open
xavfi 15 daqiqalik TTL bilan chegaralangan" degan dalil o'sha oraliqda ishlamaydi.

- [ ] **2. Dev/QA da `PRODUCTION=false` borligini tekshir.** `PRODUCTION` ning kod ichidagi
      sukut qiymati `true` (`src/infrastructure/config/env.go:39`). `.env` da aniq
      `PRODUCTION=false` yozilmagan bo'lsa, merge bo'lishi bilan dev/QA sandbox login'ni
      yo'qotadi (`POST /api/auth-v2/sandbox/login` → 404) va QA oqimi uziladi.

- [ ] **3. `make wire-build` va `make generate-docs` natijalari branch ichida commit qilinganini
      tekshir**, `go build ./...` va `go test ./... -count=1` toza o'tishini tasdiqla.

### Merge va undan keyin

- [ ] **4. Merge qil.** CI avval testlarni yuritadi, so'ng image yig'ib deploy qiladi.
      Deploydan keyin `/api/auth/logout` javob berishini tasdiqla.

- [ ] **5. Texnik tanaffusda EC kalit juftini almashtir:**

```bash
openssl ecparam -genkey -name prime256v1 -noout -out assets/keys/private_key.pem
openssl ec -in assets/keys/private_key.pem -pubout -out assets/keys/public_key.pem
```

Ilovani qayta ishga tushir. Barcha eski tokenlar — jumladan ekspertiza davrida qo'lga kiritilganlari — imzo tekshiruvidan o'tmay qoladi. Barcha foydalanuvchilar qayta kirishi kerak bo'ladi.

- [ ] **6. Muvaffaqiyat mezonlarini tekshir** (spec 12-bo'limi):

```bash
# Logout qilingandan keyin access token 401 qaytarishi kerak
curl -s -o /dev/null -w "%{http_code}\n" -H "Authorization: Bearer $ACCESS" http://localhost:8080/api/user/list

# Logout qilingandan keyin refresh token 401 qaytarishi kerak
curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/api/auth/refresh-token \
  -H 'Content-Type: application/json' -d "{\"refresh_token\":\"$REFRESH\"}"

# Redis to'xtatilganda sayt ishlashda davom etishi kerak
docker stop slib-redis-queue
curl -s -o /dev/null -w "%{http_code}\n" -H "Authorization: Bearer $ACCESS" http://localhost:8080/api/user/list
docker start slib-redis-queue
```

- [ ] **7. Frontend yangilangach:** `REFRESH_ROTATION_STRICT=true` qilib qayta ishga tushir.

- [ ] **8. Keyingi relizda:** `refresh_token_handler.go` dan query parametr tarmog'ini olib tashla.
