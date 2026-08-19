# Telefon raqami validatori (`phone_uz`) — implementatsiya rejasi

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `POST /api/auth-v2/send-otp` noto'g'ri formatdagi telefon raqamida 500 o'rniga tushunarli 400 qaytarsin va yaroqsiz so'rov throttle/DB/SMS zanjiriga umuman kirmasin.

**Architecture:** Mavjud, ammo hech qayerda chaqirilmagan `utils.IsValidPhoneNumber` funksiyasi `go-playground/validator` ga `phone_uz` nomli custom tag sifatida ulanadi. Tekshiruv `GetBody` → `c.Validate(data)` orqali handler tanasiga kirishdan oldin ishlaydi. `RequestValidator` `phone_uz` buzilganini alohida ushlab, xom go-playground xabari o'rniga foydalanuvchiga tushunarli matn qaytaradi.

**Tech Stack:** Go 1.25.0, Echo v4, `github.com/go-playground/validator/v10` v10.30.2. Modul nomi: `slib.uz`.

**Spetsifikatsiya:** `docs/superpowers/specs/2026-08-17-telefon-validatori-design.md`

## Global Constraints

- Qabul qilinadigan yagona format — aynan `998` + 9 ta raqam (jami 12 ta raqam). Normalizatsiya YO'Q, operator kodi tekshiruvi YO'Q.
- Mijozga ko'rsatiladigan xato matni aynan: `telefon raqami noto'g'ri formatda, kutilgan format: 998XXXXXXXXX`
- Tag nomi aynan: `phone_uz`
- Faqat `SendOtpRequest` ga qo'llanadi. `SandboxLoginRequest` va `check-phone-number` **tegilmaydi**.
- `phone_uz` dan boshqa taglar (`required`, `url`, `min`, `dive`) xom xabarini saqlab qoladi — mavjud endpointlar javobi o'zgarmasligi shart.
- Testlar mavjud uslubga mos: `<paket>_test` nomli tashqi test paketi, o'zbekcha xato xabarlari, `errors.As` bilan javob turini tekshirish.
- Testlar `test/` ostida `src/` strukturasini oynadek takrorlaydi.
- Barcha testlar: `make test` (ya'ni `go test ./... -count=1 -cover -coverpkg=./src/...`)

---

## File Structure

| Fayl | Mas'uliyati | Vazifa |
|---|---|---|
| `src/core/utils/phone_validator.go` | Format qoidasining yagona manbasi (regex) | 1 |
| `test/core/utils/phone_validator_test.go` | Regex chegaralarini qulflaydi | 1 |
| `src/entrypoint/presentation/interceptor/validator/request_validator.go` | Qoidani validator'ga ulaydi va xato xabarini shakllantiradi | 2 |
| `test/entrypoint/presentation/interceptor/validator/request_validator_test.go` | Tag ro'yxatdan o'tganini va xabar tanlash mantig'ini tekshiradi | 2 |
| `src/entrypoint/presentation/handlers/authv2/schema/send_otp_schema.go` | Qoidani aynan send-otp so'roviga qo'llaydi | 3 |
| `test/entrypoint/presentation/handlers/authv2/schema/send_otp_schema_test.go` | Haqiqiy schema himoyalanganini tekshiradi | 3 |

Vazifalar tartibi majburiy: 2-vazifa 1-dagi funksiyani chaqiradi, 3-vazifa 2-dagi tagni ishlatadi.

---

## Task 1: `IsValidPhoneNumber` ni test bilan qulflash va regex'ni ko'chirish

**Files:**
- Modify: `src/core/utils/phone_validator.go` (butun fayl, 8 qator)
- Create (test): `test/core/utils/phone_validator_test.go`

**Interfaces:**
- Consumes: hech narsa (birinchi vazifa)
- Produces: `utils.IsValidPhoneNumber(phone string) bool` — imzo o'zgarmaydi, 2-vazifa shuni chaqiradi

**Nega bu vazifa refaktoring, TDD emas.** `IsValidPhoneNumber` allaqachon mavjud va to'g'ri ishlaydi, shuning uchun bu yerdagi test **birinchi ishga tushirishdayoq o'tadi** — bu kutilgan holat, nosozlik emas. Testning maqsadi: regex'ni paket darajasiga ko'chirishdan oldin xulqni qulflab qo'yish. Ko'chirgandan keyin xuddi shu test qayta ishga tushiriladi va hamon o'tishi kerak.

- [ ] **Step 1: Xulqni qulflovchi testni yozish**

Yangi katalog kerak bo'ladi: `test/core/utils/`

`test/core/utils/phone_validator_test.go`:

```go
package utils_test

import (
	"testing"

	"slib.uz/src/core/utils"
)

// TestIsValidPhoneNumber formatning chegaralarini qulflaydi. Qoida: aynan
// "998" + 9 ta raqam. Normalizatsiya yo'q — "+", probel va defis rad etiladi.
func TestIsValidPhoneNumber(t *testing.T) {
	cases := []struct {
		name  string
		phone string
		want  bool
	}{
		{"to'g'ri raqam", "998901234567", true},
		{"boshqa operator kodi ham o'tadi", "998331234567", true},
		{"bir raqam ortiq", "9989012345678", false},
		{"bir raqam kam", "99890123456", false},
		{"998 prefiksi yo'q", "901234567", false},
		{"plus bilan", "+998901234567", false},
		{"probelli", "998 90 123 45 67", false},
		{"defisli", "998-90-123-45-67", false},
		{"boshqa mamlakat kodi", "997901234567", false},
		{"harf aralashgan", "99890123456a", false},
		{"butunlay harf", "abc", false},
		{"bo'sh satr", "", false},
		{"oldida qo'shimcha matn", "x998901234567", false},
		{"orqasida qo'shimcha matn", "998901234567x", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := utils.IsValidPhoneNumber(tc.phone); got != tc.want {
				t.Errorf("IsValidPhoneNumber(%q) = %v, %v kutilgandi", tc.phone, got, tc.want)
			}
		})
	}
}
```

Diqqat: `oldida qo'shimcha matn` va `orqasida qo'shimcha matn` holatlari muhim — ular regex'da `^` va `$` langarlari borligini tasdiqlaydi. Ularsiz ko'chirish paytida langarni tushirib qoldirish sezilmay qolardi.

- [ ] **Step 2: Testni ishga tushirish — O'TISHI kerak**

Run: `go test ./test/core/utils/ -v -count=1`

Expected: barcha 14 ta subtest PASS. Agar biror subtest yiqilsa, TO'XTANG — mavjud regex siz o'ylagandan boshqacha ishlayapti, davom etishdan oldin nomuvofiqlikni aniqlang.

- [ ] **Step 3: Regex'ni paket darajasiga ko'chirish**

`src/core/utils/phone_validator.go` faylini to'liq quyidagi bilan almashtiring:

```go
package utils

import "regexp"

// phoneNumberPattern — O'zbekiston telefon raqami: "998" va undan keyin 9 ta
// raqam. Paket darajasida, chunki funksiya har bir send-otp so'rovida
// chaqiriladi va MustCompile ni qayta-qayta ishga tushirish keraksiz.
var phoneNumberPattern = regexp.MustCompile(`^998[0-9]{9}$`)

// IsValidPhoneNumber raqam kutilgan formatda ekanini bildiradi. Normalizatsiya
// qilmaydi: "+", probel yoki defis bo'lsa rad etadi.
func IsValidPhoneNumber(phone string) bool {
	return phoneNumberPattern.MatchString(phone)
}
```

Naqsh oldingisi bilan bir xil ma'noda: eski `^998([0-9]{2})[0-9]{7}$` dagi qavs guruhi hech qayerda o'qilmagan, `[0-9]{2}[0-9]{7}` esa `[0-9]{9}` ga teng.

- [ ] **Step 4: Testni qayta ishga tushirish — hamon O'TISHI kerak**

Run: `go test ./test/core/utils/ -v -count=1`

Expected: xuddi o'sha 14 ta subtest PASS. Natija 2-qadamdagidan farq qilsa, ko'chirishda xato bor.

- [ ] **Step 5: Commit**

```bash
git add src/core/utils/phone_validator.go test/core/utils/phone_validator_test.go
git commit -m "refactor(utils): IsValidPhoneNumber uchun test va regex'ni paket darajasiga ko'chirish

Regex har chaqiruvda qayta kompilyatsiya qilinardi. Funksiya hech qayerda
chaqirilmagani uchun bu ahamiyatsiz edi, ammo u tez orada har bir send-otp
so'rovining yo'liga tushadi.

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 2: `phone_uz` qoidasini validator'ga ulash va xato xabarini shakllantirish

**Files:**
- Modify: `src/entrypoint/presentation/interceptor/validator/request_validator.go` (butun fayl, 26 qator)
- Create (test): `test/entrypoint/presentation/interceptor/validator/request_validator_test.go`

**Interfaces:**
- Consumes: `utils.IsValidPhoneNumber(phone string) bool` (1-vazifa)
- Produces:
  - `phone_uz` — schema taglarida ishlatiladigan qoida nomi (3-vazifa shunga tayanadi)
  - `validator.NewRequestValidator(v *validator.Validate) *RequestValidator` — imzo o'zgarmaydi
  - `(*RequestValidator).Validate(i interface{}) error` — imzo o'zgarmaydi; qaytaradigan xato hamon `*response.Response`, `Status: 400`

**Nega qoida `NewEcho` da emas, `NewRequestValidator` da ro'yxatdan o'tadi.** Agar `app/echo.go` da ro'yxatdan o'tsa, `NewRequestValidator` ni to'g'ridan-to'g'ri qurgan quyidagi testlar qoidasiz validator olardi. `phone_uz` noma'lum tag sifatida... aslida go-playground noma'lum tagda **panic** qiladi, ya'ni test ishlamay qolardi. Konstruktor ichida ro'yxatdan o'tkazish qoida obyekt bilan doim birga kelishini kafolatlaydi.

- [ ] **Step 1: Yiqiladigan testni yozish**

Yangi katalog kerak: `test/entrypoint/presentation/interceptor/validator/`

`test/entrypoint/presentation/interceptor/validator/request_validator_test.go`:

```go
package validator_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"slib.uz/src/core/application/response"
	myvalidator "slib.uz/src/entrypoint/presentation/interceptor/validator"
)

// phoneFormatMessage — mijozga ko'rsatiladigan matn. Test uni ataylab
// nusxa qilib saqlaydi: konstantani import qilsak, matn tasodifan
// o'zgarganda test buni sezmay qolardi.
const phoneFormatMessage = "telefon raqami noto'g'ri formatda, kutilgan format: 998XXXXXXXXX"

// phoneProbe — phone_uz qoidasini tekshirish uchun lokal struktura. Haqiqiy
// schema'dan mustaqil, shuning uchun bu testlar send_otp_schema o'zgarsa ham
// o'z ma'nosini yo'qotmaydi.
type phoneProbe struct {
	Phone string `json:"phone" validate:"required,phone_uz"`
}

// plainProbe — phone_uz ishlatmaydigan struktura. phone_uz dan boshqa taglar
// xom xabarini saqlab qolishini tekshirish uchun.
type plainProbe struct {
	Title string `json:"title" validate:"required"`
}

// asFailResponse xatolikni *response.Response turiga keltiradi yoki testni
// muvaffaqiyatsiz deb belgilaydi.
func asFailResponse(t *testing.T, err error) *response.Response {
	t.Helper()

	var resp *response.Response
	if !errors.As(err, &resp) {
		t.Fatalf("*response.Response kutilgandi, %T (%v) keldi", err, err)
	}
	return resp
}

func newValidator() *myvalidator.RequestValidator {
	return myvalidator.NewRequestValidator(validator.New())
}

func TestValidateAcceptsWellFormedPhone(t *testing.T) {
	if err := newValidator().Validate(&phoneProbe{Phone: "998901234567"}); err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
}

func TestValidateRejectsMalformedPhone(t *testing.T) {
	cases := []struct {
		name  string
		phone string
	}{
		{"bir raqam ortiq", "9989012345678"},
		{"998 prefiksi yo'q", "901234567"},
		{"plus bilan", "+998901234567"},
		{"harfli", "abc"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := asFailResponse(t, newValidator().Validate(&phoneProbe{Phone: tc.phone}))

			if resp.Status != 400 {
				t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
			}
			if resp.Message != phoneFormatMessage {
				t.Errorf("telefon xabari kutilgandi, %q keldi", resp.Message)
			}
		})
	}
}

// Xom go-playground xabari struct nomlarini oshkor qiladi ("Key:
// 'phoneProbe.Phone' ..."). Mijozga u yetib bormasligi kerak.
func TestValidateDoesNotLeakStructNameForPhone(t *testing.T) {
	resp := asFailResponse(t, newValidator().Validate(&phoneProbe{Phone: "abc"}))

	if strings.Contains(resp.Message, "phoneProbe") || strings.Contains(resp.Message, "Key:") {
		t.Errorf("xabar ichki tafsilotni oshkor qildi: %q", resp.Message)
	}
}

// Bo'sh maydon phone_uz ni ham buzadi, lekin sabab formatda emas — maydon
// umuman to'ldirilmagan. Foydalanuvchi "required" haqida eshitishi kerak.
func TestValidateEmptyPhoneReportsRequiredNotFormat(t *testing.T) {
	resp := asFailResponse(t, newValidator().Validate(&phoneProbe{Phone: ""}))

	if resp.Message == phoneFormatMessage {
		t.Fatal("bo'sh maydon uchun 'required' xabari kutilgandi, format xabari keldi")
	}
	if !strings.Contains(resp.Message, "required") {
		t.Errorf("'required' tagi haqida xabar kutilgandi, %q keldi", resp.Message)
	}
}

// phone_uz dan boshqa taglar hozirgi xulqini saqlaydi — bu o'zgarish mavjud
// endpointlarning javobiga tegmasligi kerak.
func TestValidateLeavesOtherTagMessagesUnchanged(t *testing.T) {
	resp := asFailResponse(t, newValidator().Validate(&plainProbe{Title: ""}))

	if resp.Status != 400 {
		t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
	}
	if !strings.Contains(resp.Message, "required") {
		t.Errorf("xom validator xabari kutilgandi, %q keldi", resp.Message)
	}
}
```

- [ ] **Step 2: Testni ishga tushirish — YIQILISHI kerak**

Run: `go test ./test/entrypoint/presentation/interceptor/validator/ -v -count=1`

Expected: FAIL. `phone_uz` hali ro'yxatdan o'tmagani uchun go-playground noma'lum tagda panic qiladi — xabar taxminan: `Undefined validation function 'phone_uz' on field 'Phone'`. `TestValidateLeavesOtherTagMessagesUnchanged` esa allaqachon o'tishi mumkin, chunki u `phone_uz` ishlatmaydi.

- [ ] **Step 3: Qoida va xabar mantig'ini yozish**

`src/entrypoint/presentation/interceptor/validator/request_validator.go` faylini to'liq quyidagi bilan almashtiring:

```go
package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/utils"
)

// phoneTag — O'zbekiston telefon raqami formatini tekshiruvchi custom qoida.
// Schema'da `validate:"required,phone_uz"` ko'rinishida ishlatiladi; tartib
// muhim, chunki go-playground birinchi buzilgan tagda to'xtaydi va bo'sh
// maydon uchun format emas, "required" xabari chiqishi kerak.
const phoneTag = "phone_uz"

// phoneFormatMessage phoneTag buzilganda mijozga boradigan matn.
// go-playground'ning xom xabari ichki struct nomlarini oshkor qiladi va
// foydalanuvchi uchun tushunarsiz.
const phoneFormatMessage = "telefon raqami noto'g'ri formatda, kutilgan format: 998XXXXXXXXX"

type RequestValidator struct {
	validator *validator.Validate
}

func NewRequestValidator(validator *validator.Validate) *RequestValidator {
	// RegisterValidation faqat tag nomi bo'sh bo'lganda xato qaytaradi, bu
	// yerda esa konstanta. Qoida shu yerda ro'yxatdan o'tadi (NewEcho da
	// emas), toki validator qurilgan har bir joyda — testlarni ham qo'shib —
	// u mavjud bo'lsin.
	_ = validator.RegisterValidation(phoneTag, isValidPhone)

	return &RequestValidator{validator: validator}
}

// isValidPhone — phoneTag qoidasining tanasi. Format qoidasining o'zi
// utils da, bitta joyda saqlanadi.
func isValidPhone(fl validator.FieldLevel) bool {
	return utils.IsValidPhoneNumber(fl.Field().String())
}

func (this *RequestValidator) Validate(i interface{}) error {
	if err := this.validator.Struct(i); err != nil {
		return response.NewFailResponse(400, validationMessage(err))
	}
	if v, ok := i.(Validatable); ok {
		if ok, err := v.Validate(); !ok && err != nil {
			return response.NewFailResponse(400, err.Error())
		}
	}
	return nil
}

// validationMessage phoneTag buzilgan bo'lsa tushunarli matn qaytaradi, aks
// holda validator xatosini o'zgarishsiz uzatadi — qolgan endpointlarning
// javobi o'zgarmasligi uchun.
func validationMessage(err error) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, fieldErr := range errs {
			if fieldErr.Tag() == phoneTag {
				return phoneFormatMessage
			}
		}
	}
	return err.Error()
}
```

Diqqat: `NewRequestValidator` parametri `validator` deb nomlangan va paket nomini soyalaydi — bu mavjud koddan meros. Funksiya ichida `validator.RegisterValidation(...)` parametrning metodini chaqiradi (paketni emas), bu to'g'ri. `isValidPhone` va `validationMessage` funksiyalari paket darajasida bo'lgani uchun ularda `validator.FieldLevel` va `validator.ValidationErrors` paketga ishora qiladi — soyalash u yerga yetib bormaydi.

- [ ] **Step 4: Testni ishga tushirish — O'TISHI kerak**

Run: `go test ./test/entrypoint/presentation/interceptor/validator/ -v -count=1`

Expected: barcha testlar PASS.

Agar `TestValidateEmptyPhoneReportsRequiredNotFormat` yiqilsa, bu go-playground taglarni e'lon tartibida qayta ishlashi haqidagi taxmin noto'g'ri ekanini bildiradi. U holda `validationMessage` ni tuzating: `phoneTag` ni qaytarishdan oldin xuddi shu maydonda `required` buzilgan-buzilmaganini tekshiring va agar buzilgan bo'lsa, xom xabarni qaytaring.

- [ ] **Step 5: Butun test to'plamini ishga tushirish**

Run: `make test`

Expected: hammasi PASS. Bu qadam `Validate` o'zgarishi boshqa handler testlarini buzmaganini tasdiqlaydi.

- [ ] **Step 6: Commit**

```bash
git add src/entrypoint/presentation/interceptor/validator/request_validator.go \
        test/entrypoint/presentation/interceptor/validator/request_validator_test.go
git commit -m "feat(validator): phone_uz custom qoidasi va tushunarli xato xabari

utils.IsValidPhoneNumber go-playground validator'ga tag sifatida ulandi.
phone_uz buzilganda mijoz struct nomlarini oshkor qiluvchi xom xabar
o'rniga tushunarli matn oladi; qolgan taglar hozirgicha qoladi.

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 3: Qoidani `send-otp` so'roviga qo'llash

**Files:**
- Modify: `src/entrypoint/presentation/handlers/authv2/schema/send_otp_schema.go:4`
- Create (test): `test/entrypoint/presentation/handlers/authv2/schema/send_otp_schema_test.go`

**Interfaces:**
- Consumes: `phone_uz` tagi va `myvalidator.NewRequestValidator` (2-vazifa)
- Produces: yakuniy xulq — hech qanday keyingi vazifa bunga tayanmaydi

- [ ] **Step 1: Yiqiladigan testni yozish**

Yangi kataloglar kerak: `test/entrypoint/presentation/handlers/authv2/schema/`

`test/entrypoint/presentation/handlers/authv2/schema/send_otp_schema_test.go`:

```go
package schema_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"slib.uz/src/core/application/response"
	"slib.uz/src/entrypoint/presentation/handlers/authv2/schema"
	myvalidator "slib.uz/src/entrypoint/presentation/interceptor/validator"
)

const phoneFormatMessage = "telefon raqami noto'g'ri formatda, kutilgan format: 998XXXXXXXXX"

// SendOtpRequest tekshiruvdan o'tmasa, yaroqsiz raqam throttle limitini
// sarflab, DB'ga OTP yozib, SMS provayderiga so'rov jo'natardi — o'sha
// provayder 503 qaytarib, mijozga 500 bo'lib yetardi.
func TestSendOtpRequestRejectsMalformedPhone(t *testing.T) {
	v := myvalidator.NewRequestValidator(validator.New())

	err := v.Validate(&schema.SendOtpRequest{Phone: "9989012345678"})

	var resp *response.Response
	if !errors.As(err, &resp) {
		t.Fatalf("*response.Response kutilgandi, %T (%v) keldi", err, err)
	}
	if resp.Status != 400 {
		t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
	}
	if resp.Message != phoneFormatMessage {
		t.Errorf("telefon xabari kutilgandi, %q keldi", resp.Message)
	}
}

func TestSendOtpRequestAcceptsWellFormedPhone(t *testing.T) {
	v := myvalidator.NewRequestValidator(validator.New())

	if err := v.Validate(&schema.SendOtpRequest{Phone: "998901234567"}); err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
}
```

- [ ] **Step 2: Testni ishga tushirish — YIQILISHI kerak**

Run: `go test ./test/entrypoint/presentation/handlers/authv2/schema/ -v -count=1`

Expected: `TestSendOtpRequestRejectsMalformedPhone` FAIL — `xatolik kutilmagandi` emas, balki `*response.Response kutilgandi, <nil> (<nil>) keldi`, chunki schema hozir `9989012345678` ni qabul qiladi. `TestSendOtpRequestAcceptsWellFormedPhone` PASS bo'ladi.

- [ ] **Step 3: Schema'ga tagni qo'shish**

`src/entrypoint/presentation/handlers/authv2/schema/send_otp_schema.go` faylini to'liq quyidagi bilan almashtiring:

```go
package schema

type SendOtpRequest struct {
	Phone string `json:"phone" validate:"required,phone_uz"`
}

type SendOtpResponse struct {
	SessionID string `json:"session_id"`
}
```

Faqat `Phone` qatori o'zgaradi. `SandboxLoginRequest` va `check-phone-number` ataylab tegilmaydi.

- [ ] **Step 4: Testni ishga tushirish — O'TISHI kerak**

Run: `go test ./test/entrypoint/presentation/handlers/authv2/schema/ -v -count=1`

Expected: ikkala test ham PASS.

- [ ] **Step 5: Butun test to'plamini ishga tushirish**

Run: `make test`

Expected: hammasi PASS.

- [ ] **Step 6: Swagger hujjatini yangilash**

Run: `make generate-docs`

Keyin `git status` bilan `src/entrypoint/presentation/docs/` ostida o'zgarish bor-yo'qligini tekshiring. Agar bor bo'lsa, uni ham commit'ga qo'shing; yo'q bo'lsa (validate tagi swagger'ga chiqmasligi mumkin), bu normal — davom eting.

- [ ] **Step 7: Commit**

```bash
git add src/entrypoint/presentation/handlers/authv2/schema/send_otp_schema.go \
        test/entrypoint/presentation/handlers/authv2/schema/send_otp_schema_test.go
git status --short src/entrypoint/presentation/docs/
# agar docs o'zgargan bo'lsa: git add src/entrypoint/presentation/docs/
git commit -m "fix(authv2): send-otp noto'g'ri telefon raqamida 400 qaytaradi

SendOtpRequest.Phone endi phone_uz qoidasidan o'tadi. Yaroqsiz raqam
handler tanasiga umuman kirmaydi, ya'ni throttle limiti sarflanmaydi,
DB'ga OTP yozuvi tushmaydi va SMS provayderiga so'rov ketmaydi.
Ilgari bu zanjir provayderdan 503 olib, mijozga 500 qaytarardi.

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Yakuniy tekshiruv

Uchala vazifa tugagach:

- [ ] `make test` — barcha testlar o'tadi
- [ ] `go vet ./...` — ogohlantirish yo'q
- [ ] `git log --oneline -3` — uchta commit ko'rinadi
- [ ] Qo'lda tekshiruv (ixtiyoriy). Server `:8080` da ishlaydi (`src/entrypoint/presentation/app/app.go:359`), `make run` bilan ko'tariladi:
  ```bash
  curl -s -o /dev/null -w '%{http_code}\n' -X POST localhost:8080/api/auth-v2/send-otp \
    -H 'Content-Type: application/json' -d '{"phone":"9989012345678"}'
  ```
  Expected: `400` (ilgari `500`). Log'da `[SmsEtcGateway]` qatori **bo'lmasligi** kerak — bu yaroqsiz so'rov tashqi xizmatga yetib bormaganini tasdiqlaydi.

## Ushbu rejadan tashqarida

Spetsifikatsiyaning 5-bo'limida qayd etilgan, ataylab qilinmaydigan ishlar:

- `sandbox/login` va `check-phone-number` endpointlarida xuddi shu tekshiruv
- Telefon raqamini normalizatsiya qilish
- Operator kodi bo'yicha tekshiruv
- SMS gateway haqiqatan nosoz bo'lganda 500 emas, 502/503 qaytarish
