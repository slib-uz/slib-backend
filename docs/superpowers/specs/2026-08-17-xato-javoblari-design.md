# Xato javoblarining izchilligi va ko'rinuvchanligi — dizayn spetsifikatsiyasi

**Sana:** 2026-08-17
**Kelib chiqishi:** `GET /api/user/list` oddiy user tokeni bilan 500 qaytargani bo'yicha nosozlik tekshiruvi
**Shox:** `develop` ustiga
**Bog'liq audit topilmasi:** 2.2.6 (CWE-285), qisman 2.2.5 (CWE-200)

---

## 1. Kelib chiqishi va tekshiruv natijasi

Boshlang'ich shikoyat: oddiy foydalanuvchi tokeni bilan `/api/user/list` ga murojaat
qilinganda 500 Internal Server Error qaytgan.

**Ildiz sabab topildi va tuzatildi** (bu ish boshlanishidan oldin):
`UserListUseCase` huquq tekshiruvida `errors.New("forbidden")` qaytargan. Loyihada
xatoni HTTP statusiga aylantirish `ResponseMiddleware`
(`response_middleware.go:30`) orqali `errors.As(err, &resp)` bilan bo'ladi — ya'ni
xato `*response.Response` bo'lishi shart. Oddiy Go xatosi bu shartga tushmaydi va
Echo'ning standart ishlovchisiga borib 500 bo'ladi. Tuzatish:
`response.PermissionDeniedError` (403) va `response.UnauthorizedError` (401).

Bu, ayni paytda, `audit-javob-report.md` ning **2.2.6-bo'limidagi da'voni to'g'ri
qildi** — hisobot «past huquqli administrator 403 oladi» degan, lekin tuzatishgacha
haqiqatda 500 qaytgan.

### 1.1. Kengaytirilgan sidiruv: dastlabki xavotir asossiz chiqdi

Dastlab «9 usecase, 14 nuqta xuddi shu nuqsonga ega» deb baholangan edi. **Bu baho
noto'g'ri bo'ldi.** Faqat usecase'lar ko'rilgan, handler'lar ko'rilmagan.

Tizimli sidiruv (oddiy xato qaytaradigan usecase'lar ∩ xatoni xom qaytaradigan
handler'lar) quyidagini ko'rsatdi:

| Dastlabki da'vo | Tekshirilgan haqiqat |
|---|---|
| 14 nuqta 500 qaytaradi | **0 nuqta** (`/user/list` dan tashqari, u tuzatilgan) |
| `fmt.Errorf` xavfli | `%w` xato zanjirini saqlaydi — `errors.As` baribir ishlaydi |
| Handler'lar xatoni xom qaytaradi | Ko'pchiligi `c.JsonResponse(400, err.Error())` bilan o'raydi |

Yolg'on ijobiylar aniqlangan sabablari:

- `article_author_affiliation_create_usecase.go:36` — kod **izohga olingan**
- `submit_for_peerreivew_usecase.go:65`, `support_dialog_create_answer_usecase.go:35` —
  `%w` bilan o'ralgan; ostidagi xato haqiqiy infratuzilma nosozligi bo'lsa 500 **to'g'ri**
- `update_authorship_claim_status_usecase.go:135` (`errors.New("Invalid ROI")`) —
  chaqiruv joyida (`:110-112`) xato **umuman qaytarilmaydi**, `fmt.Println` bilan
  yutiladi; mijozga hech qachon yetib bormaydi

**Xulosa (toraytirilgan):** bu faqat sweep aniqlagan sinf uchun to'g'ri — ya'ni
usecase `*response.Response` dan boshqa xato qaytaradi **∩** handler uni xom
qaytaradi kesishmasida qolgan nuqson yo'q. Bu status kodi bo'yicha barcha
muammolarni qamramaydi: masalan xuddi shu kunda hujjatlashtirilgan
`POST /api/auth-v2/send-otp` boshqa sinfga tegishli tirik 500'ni ko'rsatadi
(gateway 503 -> handler 500) — bu sweep'ning kesishma shartiga kirmaydi, chunki
muammo usecase/handler xato-qaytarish shaklida emas, pastki gateway javobining
yuqoriga qanday ko'tarilishida. Bu sinf ushbu to'lqin qamrovidan tashqarida.
Lekin sidiruv boshqa uchta haqiqiy muammoni ochdi — ular quyida.

---

## 2. Hal qilinadigan muammolar

### 2.1. B — mijozga bo'sh xato xabari

Uch handler `JsonResponse` ga `err.Error()` o'rniga `err` obyektini uzatadi.
`errors.errorString` ning yagona maydoni (`s`) eksport qilinmagan, shuning uchun
`encoding/json` uni `{}` deb serializatsiya qiladi:

```
{"data":{},"status":400}                                   ← hozirgi holat
{"data":"social does not belong to user","status":400}     ← kutilgan
```

Mijoz 400 oladi, lekin **nima uchun** ekanini bilmaydi.

| Fayl | Qator | Hozirgi kod |
|---|---|---|
| `handlers/social/usersocial_delete_handler.go` | 40 | `c.JsonResponse(http.StatusBadRequest, err)` |
| `handlers/social/usersocial_update_handler.go` | 53 | `ctx.JsonResponse(http.StatusBadRequest, err)` |
| `handlers/publisher/publisher_detail_handler.go` | 35 | `c.JsonResponse(400, err)` |

**O'zgarish:** `err` → `err.Error()`. Status kodi o'zgarmaydi (400 qoladi), shuning
uchun frontend uchun buzilish yo'q — faqat bo'sh `data` o'rniga matn keladi.
Loyihadagi qolgan 13 handler allaqachon aynan shu shaklda yozilgan.

### 2.2. E — yutib yuborilgan nosozliklar ko'rinmaydi

ROI (Research Output ID) sinxronizatsiyasi ikki joyda ishga tushadi va ikkalasida
ham nosozlik `fmt.Println` bilan stdout'ga chiqadi:

| Fayl | Qator |
|---|---|
| `usecase/authorshipclaimusecases/update_authorship_claim_status_usecase.go` | 111 |
| `usecase/articleusecases/article_update_usecase.go` | 105 |

Loyihada `zap` asosidagi strukturaviy logger bor (`src/infrastructure/logger/`),
kunlik fayllarga yozadi. `fmt.Println` bu infratuzilmani butunlay chetlab o'tadi —
ishlab chiqarishda ROI yuborilmagani **hech qayerda qayd etilmaydi**.

**O'zgarish:** ikki usecase konstruktoriga `*logger.AsyncLogger` inject qilinadi
va `fmt.Println` o'rniga strukturaviy yozuv:

```go
this.logger.Error("ROI yuborilmadi",
    zap.Uint("article_id", article.ID),
    zap.Error(err))
```

**Xato baribir yutiladi** (qaytarilmaydi) — bu ataylab qoldiriladi. Claim yoki
maqola allaqachon saqlangan; ROI push — yon ta'sir. Xatoni qaytarish muvaffaqiyatli
bajarilgan amaliyotni mijozga «xato» deb ko'rsatardi. O'zgarish faqat
**ko'rinuvchanlikda**, xatti-harakatda emas.

**Yangi naqsh:** hozircha hech bir usecase logger ishlatmaydi (0 ta). Bu birinchi
holat va keyingi usecase'lar uchun namuna bo'ladi. DI `@inject` annotatsiyasi va
`make wire-build` orqali qayta generatsiya qilinadi (`wiregenx` va `wire` mavjud).

### 2.3. etaqriz debug qoldig'i — CWE-200 ning log orqali ochiq qismi

```go
// usecase/etaqrizusecases/find_reviewer_usecase.go:34
reviewer, err = this.gateway.FindReviewerByScienceID(scienceID)
fmt.Println("reviewer from etaqriz/////////////:", reviewer)
```

Bu **butun `ReviewerEntity` obyektini** stdout'ga bosadi. 2.2.5 (CWE-200) ishida
`ReviewerEntity.PhoneNumber` ga `json:"-"` qo'yilgan edi — lekin `fmt.Println`
JSON teglarini **hurmat qilmaydi**. Ya'ni API javobidan olib tashlangan telefon
raqami loglarga tushishda davom etadi.

Qo'shimcha: `Println` `err` tekshiruvidan **oldin** turadi, ya'ni gateway xato
qaytarganda ham (reviewer `nil` bo'lganda) ishlaydi. Bu sof debug qoldig'i.

**O'zgarish:** qator butunlay olib tashlanadi. Mantiqqa ta'siri yo'q.

---

## 3. Qorovul

`test/architecture/error_response_test.go` — loyihadagi `raw_sql_test.go` va
`sensitive_json_test.go` naqshi bo'yicha, `go/ast` bilan.

**Qoida:** `JsonResponse(...)` yoki `JSON(...)` chaqiruvining ikkinchi argumenti
`err` deb nomlangan yalang'och identifikator bo'lishi taqiqlanadi.

**Qamrov:** qorovul `src/` daraxtini skanerlaydi (mavjud ikki qorovul bilan bir xil,
`sourceRoot()` yordamchisi qayta ishlatiladi). Metod nomi bo'yicha mos keladi —
qabul qiluvchi (`c`, `ctx` yoki boshqa) ahamiyatsiz. `JSON(...)` bugun bunday
shaklda ishlatilmaydi; u profilaktika uchun qoidaga kiritiladi.

**Nega aynan shu shakl:** qoida tor va sintaktik. `JsonResponse(200, result)` kabi
qonuniy chaqiruvlarga tegmaydi, chunki faqat `err` nomi tekshiriladi. Tuzatishdan
keyin `src/` bo'ylab mos keladigan joy qolmaydi — ya'ni **yolg'on ijobiy nol**.

**Chegarasi** (ataylab qoldiriladi va hujjatlashtiriladi): qorovul o'zgaruvchi
NOMI bo'yicha ishlaydi, TURI bo'yicha emas. Kimdir xato o'zgaruvchisini `e` yoki
`myErr` deb nomlab, uni yalang'och uzatsa — qorovul o'tkazib yuboradi. Turga
asoslangan aniq tekshiruv `go/packages` bilan to'liq tiplashtirishni talab qiladi;
mavjud ikki qorovul ham sof AST bo'lgani uchun izchillik saqlanadi. Bu chegara
`raw_sql_test.go` dagi «PROVENANCE» va «NOM TO'QNASHUVI» izohlari bilan bir xil
toifadagi kelishuv.

**Qorovulning o'z tishi:** `raw_sql_test.go:200` naqshi bo'yicha, tekshiruv
funksiyasiga ma'lum xavfli va xavfsiz shakllar berib, uning haqiqatan ishlashi
`src/` holatidan mustaqil tasdiqlanadi. Bu zarur, chunki tuzatishdan keyin `src/`
toza bo'ladi va asosiy test qorovul buzilganda ham yashil qolaverardi.

---

## 4. Testlar

- **B:** har uch handler uchun `httptest` bilan — javob tanasida xato matni bor,
  `{}` emas
- **Qorovul:** `src/` toza ekanini tasdiqlash + qorovulning o'z tishi (yuqoriga qarang)
- **E:** `zaptest/observer` bilan logger ulanadi; ROI nosozligida `Error` darajasida,
  `article_id` va `error` maydonlari bilan yozuv chiqishi tasdiqlanadi
- **E (xatti-harakat kafolati):** ROI nosozligida `Execute` baribir `nil` qaytarishi —
  ya'ni yutish saqlanib qolgani tasdiqlanadi

Ma'lumotlar bazasi talab qiladigan repozitoriylar testlanmaydi — loyihadagi mavjud
amaliyotga mos (fake repozitoriy naqshi: `user_detail_usecase_test.go`).

---

## 5. Qamrovdan tashqarida

- **C — kamida uchta xil xato kontrakti.** Bitta API'da bir nechta shakl mavjud:
  1. `JsonResponse(400, err.Error())` — xabarni `data` maydoniga qo'yadi
     (`{"status":400,"data":"..."}`), ~13 handler.
  2. `ResponseMiddleware` — sentinel xatoning (`*response.Response`) o'z tanasini
     qaytaradi, xabar `message` maydonida.
  3. `echo.NewHTTPError(400, err.Error())` — Echo'ning o'z konvertatsiyasi,
     natija `{"message":"..."}` (masalan
     `src/entrypoint/presentation/handlers/journal_rating/journal_rating_create_handler.go:36`).
  Birlashtirishning foydasi asosan estetik, narxi esa frontend uchun breaking
  change. Alohida, frontend jamoasi bilan kelishilgan ish.
- **D — avtorizatsiya xatolari 400 qaytarishi.** `socialusecases` dagi «social does
  not belong to user» egalik tekshiruvi bo'lib, 403 bo'lishi to'g'riroq. Frontend
  hozir 400 kutayotgan bo'lishi mumkin, shuning uchun alohida kelishiladi.
- **Qolgan 3 ta `fmt.Println`:** `applications_list_usecase.go:46` (debug qoldig'i),
  `application_resubmit_usecase.go:112,117` (ROI bilan bir xil yutish naqshi).
- **`publisher_detail_handler.go:31`** — `strconv.ParseUint` xatosi tekshirilmasdan
  keyingi qatorda qayta yoziladi. Alohida, mustaqil nuqson.
- **`audit-javob-report.md` ning qolgan da'volarini tasdiqlash.** 2.2.6 tuzatishgacha
  noto'g'ri edi; boshqa bo'limlarning «Natija» da'volari ham kod bilan solishtirilishi
  kerak. Alohida ish sifatida kelishilgan.

---

## 6. Frontend uchun ta'sir

Yo'q. Uchala o'zgarish ham status kodlarini o'zgartirmaydi:

- **B:** 400 → 400, faqat `data` maydoni `{}` o'rniga matn bo'ladi (yaxshilanish)
- **E:** javobga umuman tegmaydi, faqat server loglari
- **etaqriz:** javobga tegmaydi
