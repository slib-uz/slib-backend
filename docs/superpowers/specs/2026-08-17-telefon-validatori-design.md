# Telefon raqami validatori (`phone_uz`) — dizayn spetsifikatsiyasi

**Sana:** 2026-08-17
**Kelib chiqishi:** `POST /api/auth-v2/send-otp` noto'g'ri formatdagi raqamda 500 qaytargani
**Shox:** `feature/error-responses` ustiga
**Bog'liq spetsifikatsiya:** `2026-08-17-xato-javoblari-design.md` (xato javoblarining izchilligi)

---

## 1. Kelib chiqishi

Ishlab chiqarish log'ida:

```
500 | POST | /api/auth-v2/send-otp | 110ms | 2026-08-17 14:26:18
err=[SmsEtcGateway] Send unexpected status code 503, response body: map[]
```

Telefon raqami noto'g'ri formatda bo'lganda (masalan raqamlar soni ko'p) SMS
provayderi 503 qaytaradi, bu xato yuqoriga ko'tarilib mijozga 500 bo'lib yetadi.

Bu yerda ikkita alohida nuqson bor:

1. **Mijoz xatosi server xatosi sifatida ko'rsatilmoqda.** Noto'g'ri raqam —
   4xx holat. 500 esa "bizda nosozlik" degani; monitoring va mijoz ikkalasini
   ham chalg'itadi.
2. **Yaroqsiz kirish butun zanjirni kezib chiqmoqda.** Tekshiruv yo'qligi
   sababli axlat raqam throttle limitini sarflaydi, DB'ga yozuv qoldiradi va
   tashqi xizmatga so'rov jo'natadi.

### 1.1. Validator kodi allaqachon yozilgan, lekin ulanmagan

`src/core/utils/phone_validator.go`:

```go
func IsValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^998([0-9]{2})[0-9]{7}$`)
	return re.MatchString(phone)
}
```

`grep` bo'yicha bu funksiya **butun loyihada hech qayerda chaqirilmagan** — o'lik
kod. Ya'ni asosiy ish qoidani yozish emas, uni ulash.

Ulash bilan birga bitta maqsadli tuzatish kiritiladi: regex hozir har chaqiruvda
qaytadan kompilyatsiya qilinadi. Funksiya umuman chaqirilmaganda bu ahamiyatsiz
edi, ammo endi u har bir `send-otp` so'rovining yo'liga tushadi, shuning uchun
paket darajasidagi o'zgaruvchiga ko'chiriladi:

```go
var phoneNumberPattern = regexp.MustCompile(`^998[0-9]{9}$`)

func IsValidPhoneNumber(phone string) bool {
	return phoneNumberPattern.MatchString(phone)
}
```

Naqsh ma'nosi o'zgarmaydi — `([0-9]{2})[0-9]{7}` guruhi hech qayerda
ishlatilmagan, `[0-9]{9}` unga teng. Bu doiradagi yagona refaktoring.

### 1.2. Hozirgi oqim

`SendOtpRequest.Phone` da faqat `validate:"required"` bor
(`handlers/authv2/schema/send_otp_schema.go:4`), shuning uchun bo'sh bo'lmagan
har qanday satr o'tadi va `SendOtpUseCase.Execute` ga yetib boradi:

```go
if this.throttle.CheckAndHitOTPSend(ctx, phoneNumber) { ... }  // ① limit sarflanadi
otp, err := this.service.Make(ctx, phoneNumber, purpose)        // ② DB'ga yozuv
if err := this.smsGateway.Send(phoneNumber, ...); err != nil {  // ③ 503 → 500
```

---

## 2. Yechim

### 2.1. `phone_uz` qoidasini ro'yxatdan o'tkazish

`go-playground/validator` ga `phone_uz` nomli custom qoida qo'shiladi. Qoida
tanasi mavjud `utils.IsValidPhoneNumber` ni chaqiradi — regex bitta joyda
qoladi va o'lik kod ishga tushadi.

Ro'yxatdan o'tkazish joyi — `NewRequestValidator`
(`entrypoint/presentation/interceptor/validator/request_validator.go`), chunki u
allaqachon `*validator.Validate` ni qabul qiladi va egallaydi:

```go
const phoneTag = "phone_uz"

func NewRequestValidator(v *validator.Validate) *RequestValidator {
	// RegisterValidation faqat tag nomi bo'sh bo'lganda xato qaytaradi,
	// bu yerda esa konstanta — shuning uchun xato tekshirilmaydi.
	_ = v.RegisterValidation(phoneTag, func(fl validator.FieldLevel) bool {
		return utils.IsValidPhoneNumber(fl.Field().String())
	})
	return &RequestValidator{validator: v}
}
```

**Nega `NewEcho` ichida emas.** Agar qoida `app/echo.go` da ro'yxatdan o'tsa,
`NewRequestValidator` ni to'g'ridan-to'g'ri qurgan test qoidasiz validator olardi
— `phone_uz` jim e'tiborsiz qolib, test noto'g'ri "o'tdi" deb ko'rsatardi.
Konstruktor ichida esa qoida obyekt bilan birga keladi.

### 2.2. Xato xabari

`RequestValidator.Validate` hozir go-playground xatosini xom uzatadi:

```go
return response.NewFailResponse(400, err.Error())
```

Bu `phone_uz` uchun quyidagi javobni berardi:

```
Key: 'SendOtpRequest.Phone' Error:Field validation for 'Phone' failed on the 'phone_uz' tag
```

Ichki struct nomlari oshkor bo'ladi va foydalanuvchi uchun tushunarsiz. Shuning
uchun `phone_uz` buzilgan holat alohida ushlab olinadi:

```go
// phoneMessage phone_uz qoidasi buzilgan bo'lsa tushunarli matn qaytaradi,
// aks holda validator xatosini o'zgarishsiz uzatadi.
func phoneMessage(err error) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			if e.Tag() == phoneTag {
				return "telefon raqami noto'g'ri formatda, kutilgan format: 998XXXXXXXXX"
			}
		}
	}
	return err.Error()
}
```

Boshqa taglar (`required`, `url`, `min`, `dive`) hozirgidek xom qoladi —
mavjud endpointlarning javobi o'zgarmaydi.

### 2.3. Schema

`handlers/authv2/schema/send_otp_schema.go`:

```go
Phone string `json:"phone" validate:"required,phone_uz"`
```

Boshqa schema'lar tegilmaydi.

---

## 3. Kutilgan xulq

| Kirish | Hozir | Keyin |
|---|---|---|
| `998901234567` | 200 | 200 (o'zgarmaydi) |
| `9989012345678` (13 raqam) | 500 | 400 + telefon xabari |
| `901234567` (998 siz) | 500 | 400 + telefon xabari |
| `+998901234567` | 500 | 400 + telefon xabari |
| `abc` | 500 | 400 + telefon xabari |
| `""` | 400 (`required`) | 400 (`required`, o'zgarmaydi) |

Tekshiruv `GetBody` ichida, handler tanasiga kirishdan **oldin** bajariladi
(`app/context/context_utils.go` → `c.Validate(data)`). Shu sababli noto'g'ri
raqamda 1.2-bo'limdagi ①, ② va ③ qadamlarning hech biri ishlamaydi: throttle
limiti sarflanmaydi, DB'ga OTP yozuvi tushmaydi, tashqi xizmatga so'rov ketmaydi.

---

## 4. Testlar

| Fayl | Nima tekshiriladi |
|---|---|
| `test/core/utils/phone_validator_test.go` | `IsValidPhoneNumber` jadval testi. Qabul qilinadi: `998901234567`. Rad etiladi: 13 raqam, 11 raqam, `+998901234567`, probelli variant, harfli satr, bo'sh satr, `997` prefiksi |
| `test/entrypoint/presentation/interceptor/validator/request_validator_test.go` | `NewRequestValidator(validator.New())` orqali `SendOtpRequest`: (a) to'g'ri raqam xatosiz o'tadi; (b) noto'g'ri raqam `*response.Response` 400 va aynan telefon xabari bilan qaytadi; (c) bo'sh `phone` da `required` xabari chiqadi, ya'ni telefon xabari uni bosib ketmaydi |

Uslub mavjud `test/core/application/usecase/journalusecases/validate_sort_order_test.go`
ga mos: `_test` paketi, o'zbekcha xato xabarlari, javob turini `errors.As`
orqali tekshirish.

`(c)` holati alohida muhim: bo'sh satr `phone_uz` ni ham buzadi, shuning uchun
xabar tanlash mantig'i `required` ni birinchi o'ringa qo'yishi kerak.
go-playground taglarni e'lon tartibida tekshiradi va birinchi buzilishda
to'xtaydi, shuning uchun `required,phone_uz` tartibi (aksincha emas) shart.

---

## 5. Qamrovdan tashqarida

- **Normalizatsiya.** `+998 90 123-45-67` kabi ko'rinishlarni tozalab
  `998901234567` ga keltirish ko'rib chiqildi va rad etildi: frontend qat'iy
  formatda yuboradi, normalizatsiya esa DB'dagi mavjud yozuvlar bilan mos
  kelmaslik xavfini keltiradi. Kerak bo'lsa keyinroq alohida ish sifatida.

- **Operator kodini tekshirish.** `33/50/55/77/88/90/...` ro'yxati bo'yicha
  qat'iyroq tekshiruv rad etildi: yangi operator kodi chiqqanda kod yangilanishi
  kerak bo'ladi, foyda esa kichik.

- **Boshqa endpointlar.** `POST /auth-v2/sandbox/login` (`SandboxLoginRequest.PhoneNumber`)
  va `GET /auth-v2/check-phone-number` (query parametri) ham telefon raqamini
  tekshirmasdan qabul qiladi. Ular xuddi shu nuqsonga ega, lekin bu ishda
  ataylab tegilmaydi — hozirgi vazifa `send-otp` bilan cheklangan. Ma'lum
  qarz sifatida qayd etiladi.

- **SMS gateway'ning xato ishlashi.** Provayder 503 qaytarganda (haqiqiy
  nosozlik holatida) mijoz hamon 500 oladi. Bu to'g'ri xulq emas — tashqi
  xizmat nosozligi 502/503 bo'lishi kerak — lekin u xato javoblari
  spetsifikatsiyasining mavzusi, bu ishniki emas.

---

## 6. Frontend uchun ta'sir

Noto'g'ri raqam endi 500 o'rniga 400 qaytaradi. Agar frontend hozir 500 ni
"server yiqildi" deb ko'rsatayotgan bo'lsa, 400 holatida javobdagi `message`
maydonini foydalanuvchiga chiqarish kerak bo'ladi.

To'g'ri formatdagi raqamlar uchun hech narsa o'zgarmaydi.
