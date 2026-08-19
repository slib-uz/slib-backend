# SQL inyeksiya (CWE-89) — dizayn spetsifikatsiyasi

**Sana:** 2026-08-14
**Audit topilmasi:** 2.2.4 SQL Injection, xavflilik darajasi **Yuqori**
**CWE:** CWE-89 — Improper Neutralization of Special Elements used in an SQL Command
**OWASP Top 2025:** 5-o'rin (Injection)
**Shox:** `feature/sql-injection`, `feature/upload-hardening` ustiga (zanjirning 3-halqasi)

---

## 1. Zaiflik

Ekspertiza vaqtga asoslangan ko'r inyeksiyani (time-based blind SQL injection)
maqolalar ro'yxati endpointida namoyish etgan:

```
GET /api/articles/list?page=1&page_size=10&sort=-views_count          ← normal so'rov
GET /api/articles/list?page=1&page_size=10&sort=(SELECT(*)FROM(*)pg_sleep(10))
                                                → javob roppa-rosa 10 sekundda keldi
```

Sabab — `src/infrastructure/persistence/repository/published_article_repository_impl.go:401`:

```go
if sort != "" {
    if strings.HasPrefix(sort, "-") {
        sort = strings.TrimPrefix(sort, "-")
        query = query.Order(sort + " DESC")
    } else {
        query = query.Order(sort + " ASC")
    }
}
```

Foydalanuvchi yuborgan qator to'g'ridan-to'g'ri SQL matniga qo'shiladi.

### 1.1. GORM buni himoya qilmaydi

Bu taxmin emas — kutubxona darajasida tasdiqlangan. `gorm.io/gorm@v1.31.1`,
`chainable_api.go`:

```go
case string:
    if v != "" {
        tx.Statement.AddClause(clause.OrderBy{
            Columns: []clause.OrderByColumn{{
                Column: clause.Column{Name: v, Raw: true},
            }},
        })
    }
```

`Raw: true` — "bu matnni o'zgartirmasdan, qo'shtirnoqqa olmasdan SQL'ga qo'y".
GORM'ning parametrlashtirish himoyasi faqat `?` bog'lanadigan **qiymatlarga**
tegishli; ustun nomlari va `ORDER BY` ifodalari undan tashqarida.

### 1.2. Ekspluatatsiya chegarasi

Loyiha `github.com/jackc/pgx/v5` drayverini ishlatadi (`gorm.io/driver/postgres`
orqali). pgx kengaytirilgan so'rov protokolida ishlaydi va bitta so'rovda ikkinchi
buyruqni bajarishga yo'l qo'ymaydi — ya'ni `sort=id; DROP TABLE articles;--`
ko'rinishidagi ketma-ket so'rovlar (stacked queries) **bu yerda ishlamaydi**.

Bu zaiflikni kamaytirmaydi. `ORDER BY` ichiga ichki so'rov qo'yish mumkin, va
ekspertiza aynan shuni qilgan. Vaqtga asoslangan ko'r inyeksiya orqali
ma'lumotlar bazasidagi istalgan jadvalni belgima-belgi o'qib olsa bo'ladi:
parol xeshlari, PINFL, telefon raqamlari. Sekinroq, lekin to'liq. Hisobotning
17-rasmi aynan ma'lumot chiqarish jarayonini ko'rsatadi.

"Yuqori" darajasi o'rinli.

---

## 2. Bu bitta xato emas — bitta naqsh

To'liq sidiruv shuni ko'rsatdi: bir xil kod 7 ta faylda takrorlangan.

| Fayl | Qator | Metod | Parametr |
|---|---|---|---|
| `published_article_repository_impl.go` | 404, 406 | `GetAll` | `sort` — **ekspertiza sinagan** |
| `journal_repository_impl.go` | 132 | `GetListByPage` | `sort_by` + `order` |
| `report_repository_impl.go` | 141, 143 | `GetByPaging` | `ordering` |
| `journal_rating_repository_impl.go` | 71, 73 | `GetByJournalID` | `ordering` |
| `news_repository_impl.go` | 55, 57 | `GetByPaging` | `ordering` |
| `support_dialog_repository_impl.go` | 60, 62 | `GetByPaging` | `ordering` |
| `support_dialog_repository_impl.go` | 91, 93 | `GetByChatID` | `ordering` |
| `support_dialog_repository_impl.go` | 133, 135 | `GetChatsByPaging` | `ordering` |

Jami: **7 fayl, 8 metod, 15 birlashtirish nuqtasi.** Oltitasida naqsh belgima-belgi
bir xil:

```go
if strings.HasPrefix(ordering, "-") {
    sortField := strings.TrimPrefix(ordering, "-")
    query = query.Order(sortField + " DESC")
} else {
    query = query.Order(ordering + " ASC")
}
```

Bu tuzatishning shaklini belgilaydi. 15 ta qatorni alohida yamash yetarli emas —
naqshning o'zi almashtirilishi va qaytib kelishi to'silishi kerak.

### 2.1. Ro'yxat allaqachon mavjud edi

Har bir handler'da swagger izohi ruxsat etilgan qiymatlarni sanab o'tadi:

```go
// @Param sort query string false "Sort by field" Enums(views_count,rating_sum,publication_date)
// @Param ordering query string false "Ordering" Enums(created_at,-created_at)
```

Ya'ni ruxsat etilgan maydonlar ro'yxati **hujjatda bor, lekin kodda yo'q**.
Zaiflik aynan shu bo'shliqda yashaydi. Tuzatish yangi kontrakt o'ylab topmaydi —
mavjud kontraktni majburlaydi.

---

## 3. Yechim

### 3.1. `sorting` paketi

Yangi paket: `src/infrastructure/persistence/sorting`.

Nima uchun persistence qatlamida? Ustun nomlari ma'lumotlar bazasining tushunchasi,
domenning emas. Yordamchi so'rov qurilayotgan joyning yonida turishi kerak — bu
uni handler'dan repozitoriygacha bo'lgan yo'lda "unutib qo'yish"ni qiyinlashtiradi.

```go
package sorting

// Whitelist API'da ko'rinadigan tartiblash nomini SQL ifodasiga bog'laydi.
// columns qiymatlari faqat kod ichidagi konstantalar bo'lishi shart —
// ular SQL'ga o'zgarishsiz tushadi.
type Whitelist struct {
    columns      map[string]string
    defaultOrder string
}

func New(defaultOrder string, columns map[string]string) Whitelist

// Resolve bitta parametrni ishlaydi: "-created_at" → "created_at DESC"
func (w Whitelist) Resolve(ordering string) (string, error)

// ResolvePair alohida maydon va yo'nalish parametrlarini ishlaydi
// (journal ro'yxatidagi sort_by + order).
func (w Whitelist) ResolvePair(field, direction string) (string, error)
```

**Xatti-harakat jadvali:**

| Kirish | Natija |
|---|---|
| `""` | `defaultOrder`, xatosiz |
| `created_at` | `<ustun> ASC` |
| `-created_at` | `<ustun> DESC` |
| `-` (faqat prefiks) | `400` |
| ro'yxatda yo'q nom | `400` |
| `direction` bo'sh | `ASC` |
| `direction` = `asc`/`ASC`/`desc`/`DESC` | mos yo'nalish |
| boshqa `direction` | `400` |

Xatolik: `response.InvalidSortFieldError` — `NewFailResponse(400, "invalid sort field")`,
`src/core/application/response/response.go` ga qo'shiladi.

Infratuzilma qatlamining `core/application/response` dan xatolik qaytarishi
loyihada mavjud naqsh: `src/infrastructure/storage/minio_storage.go` xuddi
shunday `response.InvalidFileError` qaytaradi.

**Markaziy xususiyat:** `Resolve` qaytaradigan qator ikkita konstantadan yig'iladi —
jadvaldagi ustun ifodasi va `" ASC"` / `" DESC"`. Foydalanuvchi kiritishi natijaga
umuman tushmaydi; u faqat `map` kaliti sifatida ishlatiladi. Shu sababli
`Order`ning `Raw: true` xatti-harakati endi ahamiyatsiz bo'lib qoladi.

### 3.2. Endpoint bo'yicha ro'yxatlar

Har bir repozitoriyda paket darajasidagi `var` sifatida, tegishli metodning yonida.

| Repozitoriy / metod | API nomi | SQL ifodasi | Standart |
|---|---|---|---|
| `published_article.GetAll` | `views_count` | `articles.views_count` | `articles.publication_date DESC` |
| | `rating_sum` | `articles.rating_sum` | |
| | `publication_date` | `articles.publication_date` | |
| `journal.GetListByPage` | `views_count` | `journals.views_count` | `journals.established_date DESC` |
| | `rating_sum` | `journals.rating_sum` | |
| | `established_date` | `journals.established_date` | |
| `news.GetByPaging` | `created_at` | `news.created_at` | `news.created_at DESC` |
| `report.GetByPaging` | `created_at` | `report.created_at` | `report.created_at DESC` |
| `journal_rating.GetByJournalID` | `created_at` | `journal_ratings.created_at` | `journal_ratings.created_at DESC` |
| `support_dialog.GetByPaging` | `created_at` | `support_dialogs.created_at` | `support_dialogs.created_at DESC` |
| `support_dialog.GetByChatID` | `created_at` | `support_dialogs.created_at` | `support_dialogs.created_at DESC` |
| `support_dialog.GetChatsByPaging` | `created_at` | `support_dialogs.created_at` | `support_dialogs.created_at DESC` |

Jadvaldagi API nomlari handler'lardagi swagger `Enums(...)` ro'yxatidan olingan —
ya'ni API kontrakti o'zgarmaydi.

Jadval nomlari taxmin qilinmagan, har bir modelning `TableName()` metodidan
tasdiqlangan. Ikkitasi GORM'ning standart ko'plik qoidasiga **bo'ysunmaydi** va
ularni taxmin qilish xato bo'lardi:

- `NewsModel.TableName()` → `news` (`newss` emas)
- `ReportModel.TableName()` → `report` — **birlikda** (`reports` emas)

Qolganlari: `ArticleModel` → `articles`, `JournalModel` → `journals`,
`JournalRatingModel` → `journal_ratings`, `SupportDialogModel` → `support_dialogs`
(uchala metod ham shu modelni so'raydi).

### 3.3. Yo'l-yo'lakay tuzatiladigan mavjud nosozlik

`published_article.GetAll` so'rovi qidiruv parametrlari berilganda `journals j` ni
JOIN qiladi. `views_count` va `rating_sum` ustunlari **ikkala jadvalda ham mavjud**
(`article_model.go:34,37` va `journal_model.go:58,60`).

Ya'ni hozir hujjatlashtirilgan `?sort=views_count` so'rovi PostgreSQL'ning
`column reference "views_count" is ambiguous` xatosini beradi — qidiruv bilan birga
ishlatilganda. Ro'yxatda ustunlar `articles.views_count` deb to'liq yozilishi bu
nosozlikni o'z-o'zidan yopadi.

### 3.4. Qorovul test

`test/architecture/raw_sql_test.go` — standart kutubxonaning `go/parser` va
`go/ast` paketlari bilan butun `src/` daraxtini tahlil qiladi. Yangi bog'liqlik
kerak emas.

**Qoida:** quyidagi metodlarning **birinchi argumenti** `+` birlashtirishni
(`*ast.BinaryExpr`, `token.ADD`) yoki `fmt.Sprintf` chaqiruvini o'z ichiga olmasligi
kerak:

```
Order, Where, Or, Not, Having, Group, Select, Table, Joins, Raw, Exec
```

Nima uchun aynan birinchi argument? GORM'da birinchi argument har doim SQL matni,
qolganlari `?` ga bog'lanadigan qiymatlar. Shuning uchun:

```go
query.Where("name ILIKE ?", "%"+name+"%")   // o'tadi — `+` qiymat argumentida
query.Where("name = '" + name + "'")        // yiqiladi — `+` SQL matnida
query.Order(sort + " DESC")                 // yiqiladi
query.Order(orderExpr)                      // o'tadi — o'zgaruvchi, birlashtirish yo'q
```

Chegara aynan xavflilik chegarasida o'tadi.

Test yiqilganda xabar fayl:qator va `sorting` paketiga ishora beradi.

### 3.5. Qorovul majburlaydigan ikkita qo'shimcha tozalash

**`published_article_repository_impl.go:390`**

```go
query = query.Where("("+strings.Join(conditions, " OR ")+")", args...)
```

Bu kod xavfsiz — `conditions` faqat shu funksiya ichidagi konstanta bo'laklardan
yig'iladi. Lekin u xavfli koddan tashqi ko'rinishi bilan farq qilmaydi.
Nomlangan yordamchiga chiqariladi:

```go
// joinOr faqat shu fayldagi konstanta SQL bo'laklarini birlashtiradi.
// Foydalanuvchi kiritishi bu yerga tushmaydi — u args orqali ? ga bog'lanadi.
func joinOr(conditions []string) string {
    return "(" + strings.Join(conditions, " OR ") + ")"
}
```

**`journal_application_repository_impl.go:105`**

```go
Where(fmt.Sprintf("journals.%s = ?", issnField), issn)
```

Hozir `issnField` handler'da tekshiriladi (`CheckISSNHandler:33` — `issn_paper`
yoki `issn_online` bo'lmasa 400 qaytaradi), lekin himoya SQL qurilayotgan joydan
uzoqda turibdi — IDOR ishida ko'rgan naqshimiz. Handler tekshiruvi joyida qoladi;
repozitoriy o'zining mustaqil qorovulini oladi. To'liq SQL bo'lagi konstanta
jadvalidan olinadi:

```go
var issnConditions = map[string]string{
    "issn_paper":  "journals.issn_paper = ?",
    "issn_online": "journals.issn_online = ?",
}

condition, ok := issnConditions[issnField]
if !ok {
    return nil, response.NewFailResponse(400, "invalid ISSN type")
}
query = query.Where(condition, issn)
```

E'tibor bering: jadvalda ustun nomi emas, **butun shart** saqlanadi. Aks holda
`condition + " = ?"` yana birlashtirish bo'lardi va qorovul uni to'xtatardi —
bu qorovul kodni to'g'ri shaklga itarayotganining namunasi.

### 3.6. Qorovulning chegarasi

Buni yashirmaslik kerak: qorovul **sintaktik**. U chaqiruv joyidagi ifodani
ko'radi, o'zgaruvchining qayerdan kelganini emas.

`journal_repository_impl.go:442-479` (`GetJournalStatisticsV2`) `filterClause` ni
`+=` bilan quradi va `Raw(countQuery, filterArgs...)` ga o'zgaruvchi sifatida
beradi. Bu kod xavfsiz — barcha bo'laklar konstanta, barcha qiymatlar `?` orqali
o'tadi — lekin qorovul buni **isbotlamaydi**, shunchaki o'tkazib yuboradi.

Xuddi shu narsa `sorting.Resolve` natijasiga ham tegishli: qorovul
`query.Order(orderExpr)` ni o'tkazadi, chunki `orderExpr` — oddiy identifikator.
Uning xavfsizligini qorovul emas, `sorting` paketining testlari kafolatlaydi.

Ya'ni himoya ikki qismdan iborat va ikkalasi ham kerak:
qorovul — xavfli **shaklni** to'sadi; `sorting` testlari — xavfsiz yo'lning
haqiqatan xavfsizligini isbotlaydi.

---

## 4. Testlar

Loyihaning mavjud tartibi: testlar ildizdagi `test/` katalogida, paket tuzilishini
takrorlaydi, `package <dir>_test`. Yangi bog'liqlik yo'q — standart `testing` va
qo'lda yozilgan fake'lar.

**`test/infrastructure/persistence/sorting/whitelist_test.go`**

1. Ruxsat etilgan maydon → `<ustun> ASC`
2. `-` prefiksi → `<ustun> DESC`
3. Bo'sh qator → standart tartib, xatosiz
4. Ro'yxatda yo'q maydon → xatolik, bo'sh natija
5. Faqat `-` → xatolik
6. `pg_sleep` hujum satri (auditning aynan payloadi) → xatolik
7. `ResolvePair` yo'nalishlari: `asc`, `ASC`, `desc`, `DESC` → mos natija
8. `ResolvePair` bo'sh yo'nalish → `ASC`
9. `ResolvePair` yaroqsiz yo'nalish (`asc; DROP`) → xatolik
10. `ResolvePair` yaroqsiz maydon → xatolik

**`test/architecture/raw_sql_test.go`**

11. Butun `src/` daraxti qoidaga mos — bu test tuzatishdan **keyin** o'tishi kerak,
    tuzatishdan oldin esa 15 ta joyni sanab yiqilishi kerak. RED bosqichi zaiflikni
    ro'yxat sifatida ko'rsatadi.

**Ro'yxat kontrakti testlari** — `test/infrastructure/persistence/repository/`

12. Maqolalar ro'yxati aynan `views_count`, `rating_sum`, `publication_date` ni
    qabul qiladi va boshqa hech narsani
13. Jurnallar ro'yxati aynan `views_count`, `rating_sum`, `established_date`
14. `created_at` guruhi (news, report, rating, support ×3) aynan `created_at`

Bu testlar swagger'dagi hujjatlashtirilgan kontraktni kodga mahkamlaydi: agar
kimdir ro'yxatga maydon qo'shsa yoki olib tashlasa, test buni ko'rsatadi.

**Repozitoriylarning o'zi testlanmaydi.** Ular ma'lumotlar bazasi talab qiladi,
loyihada esa baza bilan test yo'q. Tuzatishdan keyin ular mantiqsiz o'tkazgichga
aylanadi — butun qaror `sorting` paketida, va u to'liq testlangan.

Kutilayotgan natija: 101 → ~115 test.

---

## 5. Frontend uchun ta'sir

Hujjatlashtirilgan qiymatlar (`views_count`, `rating_sum`, `publication_date`,
`established_date`, `created_at`, va ularning `-` prefiksli shakllari) avvalgidek
ishlaydi. API kontrakti o'zgarmaydi.

O'zgarish faqat bitta holatda seziladi: ro'yxatdan **tashqari** qiymat yuborilsa,
avval u jimgina SQL'ga tushardi (yoki xato berardi), endi `400 invalid sort field`
qaytadi.

**Deploy oldidan bajarilishi kerak:** frontend jamoasi bilan har bir ro'yxat
sahifasi qanday `sort` / `sort_by` / `ordering` qiymatlarini yuborayotganini
solishtirib chiqish. Agar biror sahifa jadvalda yo'q qiymat yuborayotgan bo'lsa,
u sahifa 400 oladi. Bunday qiymat topilsa — u xavfsiz ustun bo'lsa jadvalga
qo'shiladi, aks holda frontend tuzatiladi.

Bu tekshiruv zaruriy, chunki 400 qaytarish qarori ataylab tanlangan: jimgina
standart tartibga qaytish frontend xatosini yillar davomida yashirib turishi
mumkin edi.

---

## 6. Qamrovdan tashqarida

- **`GetJournalStatisticsV2` ning `filterClause` qurilishi** — xavfsiz, lekin
  qorovul buni isbotlamaydi (3.6-bo'lim). O'zgartirilmaydi.
- **Baza bilan integratsiya testlari** — alohida infratuzilma ishi. Bu ish
  qorovul va sof funksiya testlari bilan cheklanadi.
- **Boshqa inyeksiya turlari** (OS command, LDAP, template) — sidiruvda
  topilmadi, alohida ish emas.
- **Qorovulni provenance tahliliga kengaytirish** (o'zgaruvchining kelib chiqishini
  kuzatish) — hozirgi 15 ta nuqtani yopish uchun kerak emas, murakkabligi
  foydasidan katta.

---

## 7. Shox va zanjir

Bu ish `feature/upload-hardening` ustiga quriladi, chunki `test/` katalogi
`develop`da hali yo'q — u PR #18 bilan keladi. Zanjir:

```
develop ← PR #18  (feature/security-hardening)   CWE-613 + CWE-639
        ← PR #21  (feature/upload-hardening)     CWE-434 + CWE-79
        ← yangi   (feature/sql-injection)        CWE-89
```

Merge tartibi majburiy: #18 → #21 → yangi.
