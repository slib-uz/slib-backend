# Sessiya hayot siklini bekor qilish mexanizmi (CWE-613)

**Sana:** 2026-08-12
**Holat:** Tasdiqlangan, amalga oshirishga tayyor
**Zaiflik:** Ekspertiza hisoboti 2.2.1 — "Sessiyaning saqlanib qolinishi", xavflilik darajasi **Yuqori**
**CWE:** 613 (Insufficient Session Expiration) · **OWASP Top 2025:** A07 Authentication Failures

---

## 1. Muammo

"Kiberxavfsizlik markazi" ekspertizasi aniqladi: foydalanuvchi tizimdan chiqqandan keyin ham uning
avtorizatsiya tokeni uzoq muddat aktiv qolmoqda.

Kodbazani tekshirish sababni aniqladi:

| Element | Joriy holat | Manba |
|---|---|---|
| Logout endpoint | **Umuman mavjud emas** | `groups/auth_group.go` |
| Token bekor qilish | Yo'q — JWT to'liq stateless | `service/user_auth_token_service.go:83` |
| Access token TTL | 600 daqiqa = **10 soat** | `env/.env:95` |
| Refresh token TTL | 12000 daqiqa = **8.3 kun** | `env/.env:96` |
| Refresh rotatsiyasi | Yo'q — o'sha token qaytariladi | `handlers/auth/refresh_token_handler.go:37` |
| Refresh uzatilishi | Query parametr → loglarga tushadi | `handlers/auth/refresh_token_handler.go:28` |
| `jti` | Yaratiladi, lekin dekodlashda tashlanadi | `security/jwt_token_service.go:39,78` |

Amaldagi oqibat: o'g'irlangan access token 10 soat, refresh token esa 8.3 kun ishlaydi va
foydalanuvchi "chiqish" tugmasini bosishi bunga hech qanday ta'sir qilmaydi.

### 1.1. Yo'l-yo'lakay topilgan bog'liq zaiflik

`POST /api/auth-v2/sandbox/login` — himoyasiz, autentifikatsiyasiz endpoint bo'lib, haqiqiy
access va refresh token qaytaradi (`usecase/authv2usecases/sandbox_login_usecase.go:44,58`).
OTP statik holda `sandbox_users` jadvalida saqlanadi, rate limit yo'q, route prod'da shartsiz
ro'yxatdan o'tadi (`app/app.go:318`). Foydalanuvchi topilmasa yangi hisob yaratib beradi.

Bu spec'ga kiritildi: sessiyani bekor qilish mexanizmi, agar istalgan odam yangi sessiya yasay
olsa, ma'nosini yo'qotadi.

---

## 2. Qamrov

**Kiradi:**
- `POST /api/auth/logout` endpoint
- `jti` asosidagi denylist (Redis)
- Refresh token rotatsiyasi (ikki bosqichda)
- Token TTL qiymatlarini qisqartirish
- Deploy paytida EC kalit juftini almashtirish
- Sandbox login route'ini prod'dan olib tashlash

**Kirmaydi** (alohida ishlar):
- Brute-force himoyasi / CAPTCHA (ekspertiza 2.2.7)
- Deny-by-default avtorizatsiya, IDOR (2.2.3), SQL Injection (2.2.4), fayl yuklash (2.2.9)
- Foydalanuvchi bo'yicha ommaviy bekor qilish (epoch), qurilma sessiyalari ro'yxati

---

## 3. Qabul qilingan qarorlar

| Savol | Qaror | Sabab |
|---|---|---|
| Qamrov | Sessiya hayot sikli to'liq | Faqat logout qo'shish 10 soatlik TTL muammosini qoldirardi |
| Rollout | Bosqichma-bosqich, orqaga moslik bilan | Ishlab turgan frontend buzilmasligi kerak |
| Redis nosozligi | **Fail-open** + `ERROR` alert | Bitta replikasiz Redis butun saytni to'xtatmasligi kerak |
| Access TTL | **15 daqiqa** (darhol) | Fail-open holatida asosiy himoya shu |
| Refresh TTL | **7 kun** | Standart nisbat |
| Donadorlik | Faqat `jti`, epoch yo'q | Soddalik; "barcha qurilmalardan chiqish" keyinroq |
| Denylist joylashuvi | Alohida `TokenRevocationService` | Toza mas'uliyat ajratish, kelajakda kengaytirish qulay |
| Eski tokenlar | **Kalit almashtiriladi** | Ekspertiza davrida qo'lga kiritilgan tokenlar ham o'ladi |
| Sandbox eshigi | Shu spec ichida yopiladi | Usiz butun ishning ma'nosi qolmaydi |

---

## 4. Arxitektura

Denylist Redis'da bitta kalit turkumida yashaydi:

```
revoked:jti:{jti}  →  "<unix>"        — rotatsiya tufayli bekor qilingan
                   →  "<unix>:final"  — logout tufayli bekor qilingan
TTL                =  exp − now  (tokenning qolgan umri)
```

Qiymat sifatida `"1"` emas, aynan **vaqt belgisi** saqlanadi — bu 2-bosqichdagi rotatsiya "grace"
oynasini qo'shimcha kalitsiz beradi.

### 4.1. Bekor qilishning ikki turi

Qiymatdagi `:final` belgisi bekor qilish **sababini** ajratadi va bu grace oynasi uchun hal
qiluvchi:

| Tur | Qiymat | Kim yozadi | Grace |
|---|---|---|---|
| Rotatsiya | `"<unix>"` | `RefreshTokenUseCase` (`RevokeWithGrace`) | Qo'llaniladi |
| Logout | `"<unix>:final"` | `LogoutUseCase` (`Revoke`) | **Qo'llanilmaydi** |

Sabab saqlanmasa, `VerifyToken` denylist'dagi `jti` nima uchun u yerda ekanini bila olmaydi va
2-bosqichda (`REFRESH_ROTATION_STRICT=true`) logout qilingan refresh token grace oynasi davomida
qabul qilinaverardi. Oqibati: o'g'irlangan refresh tokeni bo'lgan hujumchi jabrlanuvchi chiqqanini
ko'rib, 60 soniya ichida `/auth/refresh-token` ni chaqirib, yangi va bekor qilinmagan juftlik olib
oladi — 60 soniyalik oyna 7 kunlik yangi sessiyaga aylanadi. Bu 12-bo'limdagi 2-mezonni bevosita
buzadi.

Belgisiz (`":final"` siz) qiymat **rotatsiya** deb o'qiladi: migratsiya yo'q va bosqichma-bosqich
deploy paytida eski instance yozgan yozuvlar ham shu tarmoqdan o'tishi kerak. Bu xavfsiz, chunki
grace umuman faqat `REFRESH_ROTATION_STRICT=true` da ishlaydi, u esa deploy paytida `false`.
Buzuq qiymat esa aksincha — yakuniy (`final`) deb hisoblanadi.

TTL qat'iy raqam emas, tokenning `exp` idan hisoblanadi. Sababi: konfiguratsiyadagi TTL keyin
o'zgarsa, qat'iy raqam jimgina noto'g'ri bo'lib qoladi; `exp` dan hisoblash o'z-o'zini tuzatadi.
Muddati o'tgan token uchun yozuv umuman yaratilmaydi.

---

## 5. Komponentlar

### 5.1. Yangi: `TokenRevocationService`

`src/core/application/service/token_revocation_service.go`

```go
type TokenRevocationService struct {
    cache cache.CacheProvider
}

// Revocation — denylist yozuvi: qachon va qanday bekor qilingani (4.1-bo'lim).
type Revocation struct {
    At    time.Time
    Final bool
}

// Revoke tokenni YAKUNIY ravishda denylist'ga yozadi — grace qo'llanilmaydi.
// Logout shu metodni chaqiradi.
// TTL tokenning qolgan umriga teng; muddati o'tgan token uchun hech narsa qilmaydi.
func (s *TokenRevocationService) Revoke(ctx context.Context, jti string, exp time.Time) error

// RevokeWithGrace tokenni grace oynasiga bo'ysunadigan qilib bekor qiladi.
// Faqat refresh rotatsiyasi uchun.
func (s *TokenRevocationService) RevokeWithGrace(ctx context.Context, jti string, exp time.Time) error

// RevokedAt tokenning bekor qilingan vaqti va turini qaytaradi.
// (nil, nil) — bekor qilinmagan.
// Redis nosozligida (nil, nil) qaytaradi va ERROR loglaydi (fail-open).
func (s *TokenRevocationService) RevokedAt(ctx context.Context, jti string) (*Revocation, error)
```

Yakuniy bekor qilish — sukut bo'yicha (`Revoke`), yumshoq varianti esa ataylab uzunroq nomga ega
(`RevokeWithGrace`). Kelajakda qo'shiladigan yangi chaqiruv joyi tasodifan zaifroq semantikani
tanlab qo'ymasligi uchun.

Servis `cache.CacheProvider` portiga bog'lanadi — Redis to'g'ridan-to'g'ri core qatlamiga kirmaydi.
Mavjud `infrastructure/cache/redis_cache.go` bu portni allaqachon amalga oshirgan
(`GetByKey(ctx, key)`, `Set(ctx, key, value, ttl)`).

Loglash uchun `github.com/labstack/gommon/log` ishlatiladi — bu qo'shni
`AnonymousUserTokenService` dagi mavjud uslub (`service/anonym_user_token_service.go:42`).
`zap` core qatlamida ishlatilmaydi, u faqat `infrastructure/logger` ichida qoladi.

### 5.2. Yangi: logout

- `src/core/application/usecase/authusecases/logout_usecase.go`
- `src/entrypoint/presentation/handlers/auth/logout_handler.go`
- Route: `POST /api/auth/logout`, `permissions.AuthenticatedPermission` bilan

### 5.3. O'zgaradigan komponentlar

| Fayl | O'zgarish |
|---|---|
| `core/domain/entity/token_entity.go` | `ID string` maydoni (jti) qo'shiladi |
| `infrastructure/security/jwt_token_service.go:78` | `Decode` `claims.ID` ni `TokenEntity` ga uzatadi |
| `core/application/service/user_auth_token_service.go:83` | `VerifyToken` denylist'ni so'raydi |
| `entrypoint/presentation/app/context/context_wrap.go` | `TokenID string`, `TokenExp time.Time` maydonlari |
| `interceptor/middlewares/jwt_auth_middleware.go` | Kontekstga `TokenID`/`TokenExp` yozadi |
| `handlers/auth/refresh_token_handler.go` | Refresh body'dan o'qiladi, yangi refresh qaytariladi |
| `usecase/authusecases/refresh_token_usecase.go` | Rotatsiya mantig'i |
| `groups/authv2_group.go` | `sandbox/login` faqat `PRODUCTION == false` da ro'yxatdan o'tadi |
| `infrastructure/config/env.go` | Ikkita yangi o'zgaruvchi |

### 5.4. O'zgarmaydi

Imzolash algoritmi (ES256), `AnonymousUserTokenService`, mavjud rol va permission mexanizmi.

---

## 6. Oqimlar

### 6.1. Access tokenni tekshirish

```
1. Middleware  → Authorization: Bearer ajratiladi
2. Decode      → imzo va exp tekshiriladi                (mavjud)
3. RevokedAt() → denylist so'raladi                       ← YANGI
4. DB          → foydalanuvchi o'qiladi                   (mavjud)
5. Context     → User, TokenID, TokenExp to'ldiriladi
```

Denylist tekshiruvi DB'dan **oldin** turadi: bekor qilingan token uchun ortiqcha DB so'rovi
ketmaydi. Access uchun grace oynasi yo'q — bekor qilingan bo'lsa darhol `401`.

### 6.2. Logout

```
POST /api/auth/logout          (AuthenticatedPermission)
Authorization: Bearer <access>
body: { "refresh_token": "..." }        — ixtiyoriy
```

1. Joriy access token bekor qilinadi (`c.TokenID`, `c.TokenExp`)
2. Body'da refresh bo'lsa: dekodlanadi, uning `subject` i autentifikatsiyalangan foydalanuvchiga
   tegishliligi tekshiriladi, so'ng bekor qilinadi
3. Javob qaytariladi

**2-qadamdagi egalik tekshiruvi majburiy.** Usiz A foydalanuvchi B ning refresh tokenini
denylist'ga yozib, uni tizimdan uzib qo'ya olardi — bu bekor qilish orqali DoS bo'lardi.

Refresh token body'da bo'lmasa bu xato emas: faqat access bekor qilinadi va ogohlantirish
loglanadi. Bosqichli rollout davrida eski frontend refresh yubormasligi mumkin.

### 6.3. Refresh rotatsiyasi

Xatti-harakat `REFRESH_ROTATION_STRICT` bayrog'i bilan boshqariladi.

**1-bosqich (`false`) — shu reliz bilan chiqadi:**

```
1. refresh o'qiladi     — body'dan, bo'lmasa query'dan (eski frontend uchun)
2. VerifyToken(refresh) — imzo, exp va denylist tekshiriladi
3. YANGI access + YANGI refresh yaratiladi
4. eski refresh bekor QILINMAYDI
5. ikkalasi qaytariladi
```

Eski frontend yangi refresh tokenni saqlamasa ham ishlayveradi. Shunga qaramay **logout
allaqachon to'liq ishlaydi**, chunki logout refresh `jti` sini bevosita denylist'ga yozadi va
2-qadam uni ushlaydi.

**2-bosqich (`true`) — frontend yangilangach yoqiladi:**

```
3.5. eski refresh bekor qilinadi
   + REFRESH_ROTATION_GRACE_SECONDS ichida bekor qilingan refresh qabul qilinadi
     (parallel so'rovlar va qayta urinishlar uzilmasligi uchun)
   + grace'dan tashqarida qayta ishlatilsa → WARN log (o'g'irlik signali)
```

Grace aynan shu sababdan denylist qiymatida vaqt belgisi saqlanadi.

Grace **faqat rotatsiya** bilan bekor qilingan tokenlarga tegishli. Logout bilan bekor qilingan
token (`":final"`) grace ichida ham so'zsiz rad etiladi — 4.1-bo'limga qarang.

Query parametrdan voz kechish 2-bosqichda alohida kod o'zgarishi sifatida bajariladi — bayroq
bilan emas, chunki u shunchaki olib tashlanadi.

---

## 7. Konfiguratsiya

| O'zgaruvchi | Hozir | Bo'ladi |
|---|---|---|
| `JWT_ACCESS_TOKEN_EXPIRE_MINUTES` | 600 | **15** |
| `JWT_REFRESH_TOKEN_EXPIRE_MINUTES` | 12000 | **10080** |
| `REFRESH_ROTATION_STRICT` | — | **false** |
| `REFRESH_ROTATION_GRACE_SECONDS` | — | **60** |

---

## 8. Xatolarni qayta ishlash

**O'qishda (`RevokedAt`) — fail-open.** Redis yetib bo'lmasa `(nil, nil)` qaytariladi, so'rov
o'tadi, `log.Error` orqali `jti` va xato matni bilan yozuv qoldiriladi. Foydalanuvchiga xato
ko'rinmaydi.

**Yozishda (`Revoke`) — fail-closed.** Logout javoblari:

| Holat | Javob |
|---|---|
| Token allaqachon bekor qilingan | `200` |
| Refresh yaroqsiz yoki muddati o'tgan | `200` |
| Refresh boshqa foydalanuvchiniki | `403` |
| Redis'ga yozib bo'lmadi | `503` |

Redis'ga yozilmasa chiqish haqiqatan sodir bo'lmagan. `200` qaytarish mijozga yolg'on gapirish
bo'lardi: foydalanuvchi chiqdim deb o'ylaydi, token esa ishlab turaveradi. `503` mijozga qayta
urinish imkonini beradi.

Logout token holatiga nisbatan idempotent, lekin infratuzilma nosozligini yashirmaydi.

---

## 9. Testlar

Bu kodbazadagi birinchi testlar bo'ladi (hozir `_test.go` fayllari yo'q). Yangi bog'liqlik
qo'shilmaydi: standart `testing` paketi va qo'lda yozilgan soxta `CacheProvider` (ichida `map`)
yetarli.

Testlar CI'da `go test ./... -count=1` orqali yuriladi va **build'ni to'sadi**
(`.github/workflows/test.yml`, `ci.yml` dagi `test` job). Merge avtomatik deploy qilgani uchun
yurmaydigan test himoya bermaydi.

| Test | Nimani tekshiradi |
|---|---|
| `TokenRevocationService` | TTL `exp − now` dan hisoblanishi; muddati o'tgan token yozilmasligi; fail-open; `final` belgisi yozilishi va o'qilishi; belgisiz eski qiymat rotatsiya deb o'qilishi |
| `VerifyToken` | Bekor qilingan token rad etilishi; denylist DB'dan oldin so'ralishi; grace rotatsiyaga tegishli, logoutga emas |
| `LogoutUseCase` | Access va refresh ikkalasi yozilishi; begona refresh `403` bilan rad etilishi; Redis xatosida `503`; bo'sh `jti` da `ERROR` log |
| `RefreshTokenUseCase` | 1-bosqich: eski refresh yashaydi. 2-bosqich: bekor qilinadi, grace ishlaydi, logout qilingan token grace ichida ham rad etiladi |
| `JwAuthMiddleware` | `c.TokenID` va `c.TokenExp` to'ldirilishi (logout shu chokka bog'liq) |

"Begona refresh rad etilishi" testi eng muhimi — u bekor qilish orqali DoS yo'lini yopadi.

---

## 10. Deploy tartibi

`develop`/`main` ga merge qilish deployni **avtomatik** ishga tushiradi (`.github/workflows/ci.yml`).
Ya'ni "kodni deploy qilish" alohida boshqariladigan qadam emas — u merge bilan bir vaqtda sodir
bo'ladi. Shuning uchun `.env` merge'dan **oldin** yangilanishi shart.

**Merge'dan oldin (majburiy):**

1. **Har bir serverda `env/.env` yangilanadi** (7-bo'lim) va ilova qayta ishga tushiriladi:
   `JWT_ACCESS_TOKEN_EXPIRE_MINUTES=15`, `JWT_REFRESH_TOKEN_EXPIRE_MINUTES=10080`,
   `REFRESH_ROTATION_STRICT=false`, `REFRESH_ROTATION_GRACE_SECONDS=60`.
   Sabab: `.env` keyin yangilansa, oradagi vaqtda access TTL 600 daqiqa bo'lib qolaveradi va
   "fail-open xavfi 15 daqiqalik TTL bilan chegaralangan" degan dalil o'sha oraliqda ishlamaydi.
2. **Dev/QA serverlarida `PRODUCTION=false` borligi tekshiriladi.** `PRODUCTION` ning kod
   ichidagi sukut qiymati `true` (`config/env.go:39`), shuning uchun `.env` da aniq
   `PRODUCTION=false` yozilmagan bo'lsa, merge bo'lishi bilan dev/QA sandbox login'ni yo'qotadi
   (route ro'yxatdan o'tmaydi → 404).
3. `make wire-build` va `make generate-docs` natijalari commit qilingani tekshiriladi — bu ishlar
   merge qilinadigan branch ichida bo'ladi. `google/wire` `@inject` izohlari bilan ishlaydi:
   `wire-build` o'tkazib yuborilsa konteyner yangi servisni ko'rmaydi.

**Merge va undan keyin:**

4. Merge qilinadi → CI testlarni yuritadi, so'ng image yig'iladi va deploy qilinadi
   (`REFRESH_ROTATION_STRICT=false` bilan). `/api/auth/logout` javob berishi tasdiqlanadi.
5. **Texnik tanaffus:** yangi EC kalit jufti qo'yiladi → barcha eski tokenlar kuchdan qoladi
6. Frontend yangilangach: `REFRESH_ROTATION_STRICT=true`
7. Keyingi relizda: query parametr orqali refresh qabul qilish olib tashlanadi

---

## 11. Qoldiq risklar

**Redis uzilganda bekor qilish ishlamaydi.** Fail-open tanlanganligi uchun uzilish davomida bekor
qilingan access token qabul qilinadi (eng ko'pi 15 daqiqa), bekor qilingan refresh esa yangi
tokenlar yasab olishi mumkin. Bu ongli savdo — sayt tirik qoladi.

Diqqat: hozir loyihada alert mexanizmi yo'q — loglar `logs/http` ga fayl sifatida yoziladi.
`log.Error` yozuvi ustiga haqiqiy ogohlantirish qo'yish alohida infratuzilma ishi bo'lib, bu
spec'ga kirmaydi, lekin fail-open qarorining qiymati aynan shunga bog'liq.

**Ommaviy bekor qilish yo'q.** Epoch mexanizmidan voz kechildi: "barcha qurilmalardan chiqish" va
rol o'zgarganda darhol kuchdan qoldirish yo'q. Rol o'zgarishi eng ko'pi 15 daqiqada tarqaladi.

**Refresh token 2-bosqichgacha query paramda qoladi**, ya'ni server va proxy loglariga tushishda
davom etadi. Bu ataylab qilingan — eski frontendni buzmaslik uchun.

**Grace oynasi keyingi logout'dan qat'i nazar rotatsiya bilan bekor qilingan tokenga qo'llanadi.**
Denylist qiymatida bekor qilish sababi kodlangan (`"<unix>"` — rotatsiya, `"<unix>:final"` — logout),
shuning uchun logout bilan bekor qilingan tokenga grace umuman tegmaydi. Ammo tartib teskari
bo'lganda oyna ochiq qoladi: foydalanuvchi `t=0` da refresh qilsa, o'sha eski token `"<unix>"`
belgisini oladi; `t=10s` da chiqsa, logout faqat *joriy* avlod tokenlarini `:final` qiladi —
oldingi avlod yozuvi rotatsiya belgisi bilan qoladi va `t=60s` gacha amal qiladi.

Ta'sir doirasi tor: faqat `REFRESH_ROTATION_STRICT=true` da mazmunga ega, grace bilan
(`REFRESH_ROTATION_GRACE_SECONDS`, standart 60s) chegaralangan, va hujumchi allaqachon
o'g'irlangan oldingi avlod refresh tokeniga ega bo'lishi shart. Toza yechim logout paytida
foydalanuvchining barcha avlodlarini bekor qilishni talab qiladi — bu esa 4.2-bo'limda
ataylab rad etilgan epoch mexanizmi. Shuning uchun ongli ravishda qabul qilingan.

---

## 12. Muvaffaqiyat mezonlari

1. Logout qilingandan keyin o'sha access token bilan yuborilgan so'rov `401` qaytaradi
2. Logout qilingandan keyin o'sha refresh token bilan `/auth/refresh-token` `401` qaytaradi
3. Yangi berilgan access token `exp` i 15 daqiqadan oshmaydi
4. Prod'da `POST /api/auth-v2/sandbox/login` `404` qaytaradi
5. Redis to'xtatilganda sayt ishlashda davom etadi va `ERROR` log yoziladi
6. Boshqa foydalanuvchining refresh tokenini bekor qilishga urinish `403` qaytaradi
