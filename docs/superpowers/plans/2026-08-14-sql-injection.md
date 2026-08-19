# SQL inyeksiya (CWE-89) — amalga oshirish rejasi

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `ORDER BY` ga foydalanuvchi kiritishini birlashtirish orqali yuzaga kelgan SQL inyeksiyani 15 ta nuqtada yopish va xavfli naqshning qaytib kelishini qorovul test bilan to'sish.

**Architecture:** Yangi `sorting` paketi API'da ko'rinadigan tartiblash nomini ruxsat etilgan SQL ustun ifodasiga bog'laydi. Foydalanuvchi kiritishi natijaga tushmaydi — u faqat `map` kaliti bo'lib qoladi, natija esa ikkita konstantadan yig'iladi. Har bir repozitoriy o'z ro'yxatini paket darajasidagi `var` sifatida e'lon qiladi. Oxirida `go/ast` asosidagi qorovul test butun `src/` daraxtida GORM metodlarining birinchi argumentida qator birlashtirish yoki `fmt.Sprintf` borligini taqiqlaydi.

**Tech Stack:** Go 1.25, GORM v1.31.1, PostgreSQL (pgx/v5), Echo v4. Testlar — standart `testing`, `go/parser`, `go/ast`.

## Global Constraints

- **Yangi bog'liqlik yo'q.** Faqat standart kutubxona va loyihada mavjud paketlar.
- **Testlar ildizdagi `test/` katalogida**, `src/` tuzilishini takrorlaydi, paket nomi `<dir>_test`.
- **Generatsiya qilinadigan fayllarga qo'lda tegilmaydi:** `cmd/container/container.go` (wire), `src/entrypoint/presentation/docs/*` (swaggo). Bu ishda ikkalasi ham o'zgarmasligi kerak.
- **Xatolik grantlamaydi, rad etadi.** Noma'lum tartiblash maydoni — har doim xatolik, hech qachon standart tartibga jimgina qaytish emas.
- **IDE diagnostikasi bu loyihada har doim eskirgan.** Faqat `go build ./...`, `go vet ./...`, `go test ./... -count=1` natijasiga ishoning.
- **Shox:** `feature/sql-injection`, `feature/upload-hardening` ustida. Bazaviy holat: **101 test o'tadi**, `go build ./...` toza.
- Izohlar va commit xabarlari o'zbek tilida — loyihaning mavjud odati.

---

## Fayl tuzilishi

**Yaratiladi:**

| Fayl | Mas'uliyati |
|---|---|
| `src/infrastructure/persistence/sorting/whitelist.go` | Yagona qaror nuqtasi: API nomi → SQL ifodasi. Butun ishning mantiqiy markazi. |
| `test/infrastructure/persistence/sorting/whitelist_test.go` | `sorting` xatti-harakatini to'liq qamrab oladi. |
| `test/architecture/raw_sql_test.go` | Qorovul: xavfli **shakl**ni butun `src/` da taqiqlaydi. |
| `test/infrastructure/persistence/repository/sort_whitelist_test.go` | Ro'yxatlar kontrakti: swagger `Enums(...)` bilan moslik. |

**O'zgartiriladi:**

| Fayl | O'zgarish |
|---|---|
| `src/core/application/response/response.go` | `InvalidSortFieldError` qo'shiladi |
| `src/infrastructure/persistence/repository/published_article_repository_impl.go` | ro'yxat + `Resolve`, `joinOr` yordamchisi |
| `src/infrastructure/persistence/repository/journal_repository_impl.go` | ro'yxat + `ResolvePair` |
| `src/infrastructure/persistence/repository/news_repository_impl.go` | ro'yxat + `Resolve` |
| `src/infrastructure/persistence/repository/report_repository_impl.go` | ro'yxat + `Resolve` |
| `src/infrastructure/persistence/repository/journal_rating_repository_impl.go` | ro'yxat + `Resolve` |
| `src/infrastructure/persistence/repository/support_dialog_repository_impl.go` | ro'yxat + `Resolve` (3 metod) |
| `src/infrastructure/persistence/repository/journal_application_repository_impl.go` | ISSN shartlar jadvali |

**Import tozalash — har bir faylda aniq:** tuzatishdan keyin `strings` importi
`journal_repository`, `news_repository`, `report_repository`,
`journal_rating_repository`, `support_dialog_repository` fayllarida
ishlatilmay qoladi va olib tashlanishi kerak. `journal_application_repository`
da `fmt` importi ishlatilmay qoladi. `published_article_repository` da `strings`
**qoladi** — uni `joinOr` ishlatadi. Buni `go build ./...` baribir ushlaydi.

---

## Boshlashdan oldin: bazaviy holatni tasdiqlash

```bash
git branch --show-current     # feature/sql-injection kutiladi
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -20
```

101 test o'tishi kerak. O'tmasa — **to'xtang va xabar bering**, bu rejaning aybi emas.

---

### Task 1: `sorting` paketi

Butun ishning mantiqiy markazi. Bundan keyingi barcha tasklar shu paketning
o'tkazgichlari bo'ladi.

**Files:**
- Create: `src/infrastructure/persistence/sorting/whitelist.go`
- Modify: `src/core/application/response/response.go`
- Test: `test/infrastructure/persistence/sorting/whitelist_test.go`

**Interfaces:**
- Consumes: `slib.uz/src/core/application/response` (mavjud paket)
- Produces:
  - `sorting.New(defaultOrder string, columns map[string]string) sorting.Whitelist`
  - `(sorting.Whitelist).Resolve(ordering string) (string, error)`
  - `(sorting.Whitelist).ResolvePair(field, direction string) (string, error)`
  - `(sorting.Whitelist).Fields() []string` — alifbo tartibida saralangan
  - `response.InvalidSortFieldError`

- [ ] **Step 1: Xatolikni qo'shish**

`src/core/application/response/response.go`, `InvalidFileError` qatoridan keyin:

```go
	InvalidFileError   = NewFailResponse(400, "invalid file")
	EntityToLargeError = NewFailResponse(413, "entity too large")

	// InvalidSortFieldError — so'rovda ruxsat etilmagan tartiblash maydoni
	// kelgan. 400, chunki bu mijoz xatosi; jimgina standart tartibga qaytish
	// frontend nosozligini yashirib qo'yardi.
	InvalidSortFieldError = NewFailResponse(400, "invalid sort field")
```

- [ ] **Step 2: Yiqiladigan testni yozish**

Yangi fayl `test/infrastructure/persistence/sorting/whitelist_test.go`:

```go
package sorting_test

import (
	"errors"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/infrastructure/persistence/sorting"
)

// newsFields — testlar uchun namunaviy ro'yxat. Haqiqiy ro'yxatlar
// repozitoriylarda e'lon qilinadi; bu yerda faqat xatti-harakat sinaladi.
func newsFields() sorting.Whitelist {
	return sorting.New("news.created_at DESC", map[string]string{
		"created_at": "news.created_at",
		"views":      "news.views_count",
	})
}

func TestResolveAllowedFieldAscending(t *testing.T) {
	got, err := newsFields().Resolve("created_at")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "news.created_at ASC" {
		t.Errorf("%q kutilgandi, %q keldi", "news.created_at ASC", got)
	}
}

func TestResolveMinusPrefixMeansDescending(t *testing.T) {
	got, err := newsFields().Resolve("-created_at")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "news.created_at DESC" {
		t.Errorf("%q kutilgandi, %q keldi", "news.created_at DESC", got)
	}
}

func TestResolveEmptyReturnsDefault(t *testing.T) {
	got, err := newsFields().Resolve("")
	if err != nil {
		t.Fatalf("bo'sh qiymat xatolik bermasligi kerak: %v", err)
	}
	if got != "news.created_at DESC" {
		t.Errorf("standart tartib kutilgandi, %q keldi", got)
	}
}

func TestResolveUnknownFieldIsRejected(t *testing.T) {
	got, err := newsFields().Resolve("password_hash")
	if !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("InvalidSortFieldError kutilgandi, %v keldi", err)
	}
	if got != "" {
		t.Errorf("xatolik bilan birga bo'sh natija kutilgandi, %q keldi", got)
	}
}

func TestResolveBareMinusIsRejected(t *testing.T) {
	if _, err := newsFields().Resolve("-"); !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("InvalidSortFieldError kutilgandi, %v keldi", err)
	}
}

// Auditning aynan payload'i. Bu test zaiflikni to'g'ridan-to'g'ri mahkamlaydi.
func TestResolveRejectsAuditPayload(t *testing.T) {
	payload := "(SELECT(*)FROM(*)pg_sleep(10))"
	got, err := newsFields().Resolve(payload)
	if !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("hujum satri rad etilishi kerak edi, xatolik: %v", err)
	}
	if got != "" {
		t.Fatalf("hujum satri natijaga tushdi: %q", got)
	}
}

func TestResolvePairDirections(t *testing.T) {
	cases := map[string]string{
		"":     "news.created_at ASC",
		"asc":  "news.created_at ASC",
		"ASC":  "news.created_at ASC",
		"desc": "news.created_at DESC",
		"DESC": "news.created_at DESC",
	}
	for direction, want := range cases {
		got, err := newsFields().ResolvePair("created_at", direction)
		if err != nil {
			t.Errorf("%q: xatolik kutilmagandi: %v", direction, err)
			continue
		}
		if got != want {
			t.Errorf("%q: %q kutilgandi, %q keldi", direction, want, got)
		}
	}
}

func TestResolvePairRejectsUnknownDirection(t *testing.T) {
	if _, err := newsFields().ResolvePair("created_at", "asc; DROP TABLE news"); !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("InvalidSortFieldError kutilgandi, %v keldi", err)
	}
}

func TestResolvePairRejectsUnknownField(t *testing.T) {
	if _, err := newsFields().ResolvePair("pg_sleep(10)", "desc"); !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("InvalidSortFieldError kutilgandi, %v keldi", err)
	}
}

func TestResolvePairEmptyFieldReturnsDefault(t *testing.T) {
	got, err := newsFields().ResolvePair("", "desc")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "news.created_at DESC" {
		t.Errorf("standart tartib kutilgandi, %q keldi", got)
	}
}

func TestFieldsIsSorted(t *testing.T) {
	got := newsFields().Fields()
	want := []string{"created_at", "views"}
	if len(got) != len(want) {
		t.Fatalf("%v kutilgandi, %v keldi", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%v kutilgandi, %v keldi", want, got)
		}
	}
}
```

**Diqqat — `errors.Is` bu yerda ishlashi kerak.** `response.InvalidSortFieldError`
— bu `*response.Response` ko'rsatkichi va u paket darajasidagi yagona nusxa.
`errors.Is` ko'rsatkichlarni solishtiradi, shuning uchun `Resolve` **aynan shu
o'zgaruvchini** qaytarishi shart, uni nusxalamasligi yoki o'ramasligi kerak.

- [ ] **Step 3: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/infrastructure/persistence/sorting/... -count=1
```

Kutilgan: kompilyatsiya xatosi — `sorting` paketi mavjud emas.

- [ ] **Step 4: `sorting` paketini yozish**

Yangi fayl `src/infrastructure/persistence/sorting/whitelist.go`:

```go
// Package sorting so'rov parametridan kelgan tartiblash nomini ruxsat etilgan
// SQL ustun ifodasiga aylantiradi.
//
// Nima uchun bu kerak: GORM'ning Order(string) metodi qiymatni
// clause.Column{Raw: true} qilib qo'yadi, ya'ni matn SQL'ga o'zgarishsiz
// tushadi. GORM'ning parametrlashtirish himoyasi faqat ? ga bog'lanadigan
// qiymatlarga tegishli; ustun nomlari undan tashqarida qoladi.
//
// Paket persistence qatlamida turadi, chunki ustun nomlari ma'lumotlar
// bazasining tushunchasi — domenning emas.
package sorting

import (
	"sort"
	"strings"

	"slib.uz/src/core/application/response"
)

const (
	ascending  = " ASC"
	descending = " DESC"
)

// Whitelist API'da ko'rinadigan tartiblash nomini SQL ifodasiga bog'laydi.
//
// columns qiymatlari va defaultOrder faqat kod ichidagi konstantalar bo'lishi
// SHART — ular SQL matniga o'zgarishsiz tushadi. Foydalanuvchi kiritishi
// natijaga hech qachon qo'shilmaydi: u faqat map kaliti sifatida ishlatiladi.
type Whitelist struct {
	columns      map[string]string
	defaultOrder string
}

// New ro'yxat yaratadi. defaultOrder — parametr berilmaganda ishlatiladigan
// to'liq ORDER BY ifodasi (masalan "articles.publication_date DESC").
func New(defaultOrder string, columns map[string]string) Whitelist {
	return Whitelist{columns: columns, defaultOrder: defaultOrder}
}

// Fields ruxsat etilgan API nomlarini alifbo tartibida qaytaradi.
// Kontrakt testlari uchun: ro'yxat swagger hujjatiga mosligini mahkamlaydi.
func (w Whitelist) Fields() []string {
	fields := make([]string, 0, len(w.columns))
	for field := range w.columns {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}

// Resolve bitta parametrni ishlaydi: "-created_at" -> "<ustun> DESC".
// Bo'sh qiymat standart tartibni beradi. Ro'yxatda yo'q nom — xatolik.
func (w Whitelist) Resolve(ordering string) (string, error) {
	if ordering == "" {
		return w.defaultOrder, nil
	}

	field := ordering
	direction := ascending
	if strings.HasPrefix(ordering, "-") {
		field = strings.TrimPrefix(ordering, "-")
		direction = descending
	}

	column, ok := w.columns[field]
	if !ok {
		return "", response.InvalidSortFieldError
	}

	// Ikkala bo'lak ham konstanta: column — jadvaldan, direction — yuqoridagi
	// const bloklardan. Foydalanuvchi kiritishi bu yerga tushmaydi.
	return column + direction, nil
}

// ResolvePair maydon va yo'nalish alohida parametr bo'lganda ishlatiladi
// (jurnallar ro'yxatidagi sort_by + order).
func (w Whitelist) ResolvePair(field, direction string) (string, error) {
	if field == "" {
		return w.defaultOrder, nil
	}

	column, ok := w.columns[field]
	if !ok {
		return "", response.InvalidSortFieldError
	}

	switch strings.ToUpper(direction) {
	case "", "ASC":
		return column + ascending, nil
	case "DESC":
		return column + descending, nil
	default:
		return "", response.InvalidSortFieldError
	}
}
```

- [ ] **Step 5: Testlarni ishga tushirish**

```bash
go build ./... && go test ./test/infrastructure/persistence/sorting/... -count=1 -v
go test ./... -count=1 2>&1 | tail -20
```

Kutilgan: 11 ta test funksiyasi ham PASS, umumiy hisob **112** (101 + 11).

`errors.Is` yiqilsa — `Resolve` xatolikni o'rab yubormaganini tekshiring.

- [ ] **Step 6: Mutatsiya tekshiruvi — testlar haqiqatan tishlaydimi**

Testlar o'tishi ular biror narsani ushlashini isbotlamaydi. Vaqtincha
`whitelist.go` dagi `Resolve` ichidagi ro'yxat tekshiruvini buzing —
`return "", response.InvalidSortFieldError` o'rniga:

```go
	column, ok := w.columns[field]
	if !ok {
		column = field          // MUTATSIYA: eski xavfli xatti-harakat
	}
```

```bash
go test ./test/infrastructure/persistence/sorting/... -count=1
```

Kutilgan: `TestResolveUnknownFieldIsRejected`, `TestResolveBareMinusIsRejected`
va `TestResolveRejectsAuditPayload` yiqiladi. **Yiqilmasa — testlar
yaroqsiz, to'xtang va xabar bering.**

Keyin mutatsiyani qo'lda bekor qiling (fayl hali commit qilinmagan, shuning
uchun `git checkout` yordam bermaydi): `column = field` qatorini o'chirib,
`return "", response.InvalidSortFieldError` ni qaytaring.

```bash
go test ./test/infrastructure/persistence/sorting/... -count=1
```

Kutilgan: hammasi yana PASS.

- [ ] **Step 7: Commit**

```bash
git add src/infrastructure/persistence/sorting/whitelist.go \
        src/core/application/response/response.go \
        test/infrastructure/persistence/sorting/whitelist_test.go
git commit -m "feat(sorting): tartiblash maydonlari uchun allow-list

GORM'ning Order(string) metodi qiymatni Raw: true bilan qo'yadi, ya'ni
matn SQL'ga o'zgarishsiz tushadi. Endi foydalanuvchi kiritishi natijaga
umuman qo'shilmaydi — u faqat map kaliti, natija esa ikkita konstantadan
yig'iladi.

Auditning aynan payload'i test bilan mahkamlandi."
```

---

### Task 2: Maqolalar va jurnallar ro'yxati

Ekspertiza aynan maqolalar endpointini sinagan — bu task auditning topilmasini
yopadi.

**Files:**
- Modify: `src/infrastructure/persistence/repository/published_article_repository_impl.go:401-410`
- Modify: `src/infrastructure/persistence/repository/journal_repository_impl.go:130-135`
- Test: `test/infrastructure/persistence/repository/sort_whitelist_test.go` (yaratiladi)

**Interfaces:**
- Consumes: `sorting.New`, `(Whitelist).Resolve`, `(Whitelist).ResolvePair`, `(Whitelist).Fields` (Task 1)
- Produces:
  - `repository.ArticleSortFields` (`sorting.Whitelist`)
  - `repository.JournalSortFields` (`sorting.Whitelist`)

Ro'yxatlar **eksport qilinadi**, chunki kontrakt testlari `repository_test`
paketidan turib ularni ko'rishi kerak.

- [ ] **Step 1: Yiqiladigan kontrakt testini yozish**

Yangi fayl `test/infrastructure/persistence/repository/sort_whitelist_test.go`:

```go
package repository_test

import (
	"testing"

	"slib.uz/src/infrastructure/persistence/repository"
)

func assertFields(t *testing.T, name string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: %v kutilgandi, %v keldi", name, want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s: %v kutilgandi, %v keldi", name, want, got)
		}
	}
}

// Ro'yxatlar handler'lardagi swagger Enums(...) hujjatiga mos bo'lishi kerak.
// Bu test API kontraktini kodga mahkamlaydi: kimdir ro'yxatga maydon qo'shsa
// yoki olib tashlasa, hujjat bilan farq shu yerda ko'rinadi.
func TestArticleSortFieldsMatchSwaggerEnums(t *testing.T) {
	// articles_list_handler.go: Enums(views_count,rating_sum,publication_date)
	assertFields(t, "articles",
		repository.ArticleSortFields.Fields(),
		[]string{"publication_date", "rating_sum", "views_count"})
}

func TestJournalSortFieldsMatchSwaggerEnums(t *testing.T) {
	// journal_list_handler.go: Enums(views_count,rating_sum,established_date)
	assertFields(t, "journals",
		repository.JournalSortFields.Fields(),
		[]string{"established_date", "rating_sum", "views_count"})
}

// Maqolalar so'rovi qidiruv parametrlari berilganda journals jadvalini JOIN
// qiladi, va views_count/rating_sum ustunlari IKKALA jadvalda ham mavjud.
// Ustunlar to'liq nom bilan yozilmasa, PostgreSQL "ambiguous column" xatosini
// beradi — ya'ni hujjatlashtirilgan qiymat ishlamaydi.
func TestArticleSortColumnsAreTableQualified(t *testing.T) {
	for _, field := range repository.ArticleSortFields.Fields() {
		got, err := repository.ArticleSortFields.Resolve(field)
		if err != nil {
			t.Fatalf("%s: xatolik kutilmagandi: %v", field, err)
		}
		if len(got) < 9 || got[:9] != "articles." {
			t.Errorf("%s: ustun \"articles.\" bilan boshlanishi kerak, %q keldi", field, got)
		}
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/infrastructure/persistence/repository/... -count=1
```

Kutilgan: kompilyatsiya xatosi — `ArticleSortFields` va `JournalSortFields`
aniqlanmagan.

- [ ] **Step 3: Maqolalar ro'yxatini qo'shish**

`published_article_repository_impl.go` — import blokiga qo'shing:

```go
	"slib.uz/src/infrastructure/persistence/sorting"
```

`GetAll` funksiyasidan **oldin**, fayl darajasida:

```go
// ArticleSortFields — /api/articles/list uchun ruxsat etilgan tartiblash
// maydonlari. Nomlar articles_list_handler.go dagi Enums(...) ro'yxatiga mos.
//
// Ustunlar jadval nomi bilan to'liq yozilgan: qidiruv parametrlari berilganda
// so'rov journals jadvalini JOIN qiladi, va views_count/rating_sum ikkala
// jadvalda ham bor.
var ArticleSortFields = sorting.New("articles.publication_date DESC", map[string]string{
	"views_count":      "articles.views_count",
	"rating_sum":       "articles.rating_sum",
	"publication_date": "articles.publication_date",
})
```

- [ ] **Step 4: `GetAll` ichidagi tartiblashni almashtirish**

`GetAll` boshida, `var articles []*models.ArticleModel` qatoridan keyin qo'shing:

```go
	// Tartiblash birinchi tekshiriladi: yaroqsiz qiymat uchun bazaga
	// murojaat qilishning hojati yo'q.
	orderExpr, err := ArticleSortFields.Resolve(sort)
	if err != nil {
		return nil, err
	}
```

Keyin 401-410 qatorlardagi blokni:

```go
	if sort != "" {
		if strings.HasPrefix(sort, "-") {
			sort = strings.TrimPrefix(sort, "-")
			query = query.Order(sort + " DESC")
		} else {
			query = query.Order(sort + " ASC")
		}
	} else {
		query = query.Order("publication_date DESC")
	}
```

bittagina qatorga almashtiring:

```go
	query = query.Order(orderExpr)
```

- [ ] **Step 5: Jurnallar ro'yxatini qo'shish**

`journal_repository_impl.go` — import blokiga qo'shing:

```go
	"slib.uz/src/infrastructure/persistence/sorting"
```

`GetListByPage` funksiyasidan **oldin**, fayl darajasida:

```go
// JournalSortFields — jurnallar ro'yxati uchun ruxsat etilgan tartiblash
// maydonlari. Nomlar journal_list_handler.go dagi Enums(...) ro'yxatiga mos.
var JournalSortFields = sorting.New("journals.established_date DESC", map[string]string{
	"views_count":      "journals.views_count",
	"rating_sum":       "journals.rating_sum",
	"established_date": "journals.established_date",
})
```

`GetListByPage` boshida, `var journals []*models.JournalModel` qatoridan keyin:

```go
	orderExpr, err := JournalSortFields.ResolvePair(sortBy, order)
	if err != nil {
		return nil, err
	}
```

130-135 qatorlardagi blokni:

```go
	// Apply sorting
	if sortBy != "" {
		query = query.Order(sortBy + " " + strings.ToUpper(order))
	} else {
		query = query.Order("established_date DESC")
	}
```

almashtiring:

```go
	query = query.Order(orderExpr)
```

- [ ] **Step 6: Ishlatilmay qolgan importni olib tashlash**

`journal_repository_impl.go` da `strings` endi ishlatilmaydi — import blokidan
`"strings"` qatorini o'chiring.

`published_article_repository_impl.go` da `strings` **qoladi** (390-qatordagi
`strings.Join` uni ishlatadi) — tegmang.

- [ ] **Step 7: Qurish va testlar**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -20
```

Kutilgan: qurish toza, **115** test o'tadi (112 + 3 ta yangi).

`strings` importi haqida xato chiqsa — 6-qadamni bajarmagansiz.

- [ ] **Step 8: Commit**

```bash
git add src/infrastructure/persistence/repository/published_article_repository_impl.go \
        src/infrastructure/persistence/repository/journal_repository_impl.go \
        test/infrastructure/persistence/repository/sort_whitelist_test.go
git commit -m "fix(sql): maqolalar va jurnallar ro'yxatida ORDER BY inyeksiyasi yopildi

Ekspertiza aynan /api/articles/list?sort= parametrini sinagan va
pg_sleep(10) bilan vaqtga asoslangan ko'r inyeksiyani ko'rsatgan.

Yo'l-yo'lakay mavjud nosozlik ham tuzaldi: views_count va rating_sum
ustunlari articles va journals jadvallarida ham bor, shuning uchun
qidiruv bilan birga ?sort=views_count PostgreSQL'ning ambiguous column
xatosini berardi. Endi ustunlar to'liq nom bilan yoziladi."
```

---

### Task 3: Yangiliklar, shikoyatlar va reytinglar

Uchta bir xil tuzatish. Naqsh belgima-belgi bir xil, faqat jadval nomi farq qiladi.

**Files:**
- Modify: `src/infrastructure/persistence/repository/news_repository_impl.go:52-59`
- Modify: `src/infrastructure/persistence/repository/report_repository_impl.go:138-145`
- Modify: `src/infrastructure/persistence/repository/journal_rating_repository_impl.go:68-75`
- Modify: `test/infrastructure/persistence/repository/sort_whitelist_test.go`

**Interfaces:**
- Consumes: `sorting.New`, `(Whitelist).Resolve`, `(Whitelist).Fields` (Task 1)
- Produces: `repository.NewsSortFields`, `repository.ReportSortFields`, `repository.JournalRatingSortFields`

**Jadval nomlari taxmin qilinmagan** — har biri modelning `TableName()`
metodidan olingan. Ikkitasi GORM'ning standart ko'plik qoidasiga bo'ysunmaydi:

| Model | `TableName()` |
|---|---|
| `NewsModel` | `news` — `newss` emas |
| `ReportModel` | `report` — **birlikda**, `reports` emas |
| `JournalRatingModel` | `journal_ratings` |

- [ ] **Step 1: Yiqiladigan testlarni qo'shish**

`test/infrastructure/persistence/repository/sort_whitelist_test.go` fayliga
qo'shing:

```go
func TestNewsSortFieldsMatchSwaggerEnums(t *testing.T) {
	// news_list_handler.go: Enums(created_at,-created_at)
	assertFields(t, "news",
		repository.NewsSortFields.Fields(),
		[]string{"created_at"})
}

func TestReportSortFieldsMatchSwaggerEnums(t *testing.T) {
	// report_list_handler.go: Enums(created_at,-created_at)
	assertFields(t, "report",
		repository.ReportSortFields.Fields(),
		[]string{"created_at"})
}

func TestJournalRatingSortFieldsMatchSwaggerEnums(t *testing.T) {
	// journal_rating_list_handler.go: Enums(created_at,-created_at)
	assertFields(t, "journal_ratings",
		repository.JournalRatingSortFields.Fields(),
		[]string{"created_at"})
}

// ReportModel.TableName() "report" qaytaradi — BIRLIKDA. GORM'ning standart
// ko'plik qoidasiga ishonib "reports" deb yozilsa, so'rov ishga tushganda
// relation "reports" does not exist xatosi chiqadi.
func TestReportSortColumnUsesSingularTableName(t *testing.T) {
	got, err := repository.ReportSortFields.Resolve("created_at")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "report.created_at ASC" {
		t.Errorf("%q kutilgandi, %q keldi", "report.created_at ASC", got)
	}
}
```

- [ ] **Step 2: Testlarni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/infrastructure/persistence/repository/... -count=1
```

Kutilgan: kompilyatsiya xatosi — uchala ro'yxat ham aniqlanmagan.

- [ ] **Step 3: Yangiliklar repozitoriysini tuzatish**

`news_repository_impl.go` — import blokiga `"slib.uz/src/infrastructure/persistence/sorting"`
qo'shing, `"strings"` ni o'chiring.

`GetByPaging` dan oldin, fayl darajasida:

```go
// NewsSortFields — yangiliklar ro'yxati uchun ruxsat etilgan tartiblash
// maydonlari. NewsModel.TableName() "news" qaytaradi.
var NewsSortFields = sorting.New("news.created_at DESC", map[string]string{
	"created_at": "news.created_at",
})
```

`GetByPaging` boshida, `var total int64` qatoridan oldin:

```go
	orderExpr, err := NewsSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}
```

Keyin bu blokni:

```go
	// Apply ordering
	if ordering != "" {
		if strings.HasPrefix(ordering, "-") {
			sortField := strings.TrimPrefix(ordering, "-")
			query = query.Order(sortField + " DESC")
		} else {
			query = query.Order(ordering + " ASC")
		}
	} else {
		// Default ordering by created_at DESC
		query = query.Order("created_at DESC")
	}
```

almashtiring:

```go
	query = query.Order(orderExpr)
```

- [ ] **Step 4: Shikoyatlar repozitoriysini tuzatish**

`report_repository_impl.go` — import blokiga `sorting` qo'shing, `"strings"` ni
o'chiring.

`GetByPaging` dan oldin:

```go
// ReportSortFields — shikoyatlar ro'yxati uchun ruxsat etilgan tartiblash
// maydonlari. ReportModel.TableName() "report" qaytaradi — birlikda.
var ReportSortFields = sorting.New("report.created_at DESC", map[string]string{
	"created_at": "report.created_at",
})
```

`GetByPaging` boshida, `var total int64` qatoridan oldin:

```go
	orderExpr, err := ReportSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}
```

`if ordering != "" { ... } else { query = query.Order("created_at DESC") }`
blokini almashtiring:

```go
	query = query.Order(orderExpr)
```

- [ ] **Step 5: Reytinglar repozitoriysini tuzatish**

`journal_rating_repository_impl.go` — import blokiga `sorting` qo'shing,
`"strings"` ni o'chiring.

`GetByJournalID` dan oldin:

```go
// JournalRatingSortFields — jurnal reytinglari ro'yxati uchun ruxsat etilgan
// tartiblash maydonlari.
var JournalRatingSortFields = sorting.New("journal_ratings.created_at DESC", map[string]string{
	"created_at": "journal_ratings.created_at",
})
```

`GetByJournalID` boshida, `var total int64` qatoridan oldin:

```go
	orderExpr, err := JournalRatingSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}
```

`if ordering != "" { ... } else { query = query.Order("created_at DESC") }`
blokini almashtiring:

```go
	query = query.Order(orderExpr)
```

- [ ] **Step 6: Qurish va testlar**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -20
```

Kutilgan: qurish toza, **119** test o'tadi (115 + 4 ta yangi).

- [ ] **Step 7: Commit**

```bash
git add src/infrastructure/persistence/repository/news_repository_impl.go \
        src/infrastructure/persistence/repository/report_repository_impl.go \
        src/infrastructure/persistence/repository/journal_rating_repository_impl.go \
        test/infrastructure/persistence/repository/sort_whitelist_test.go
git commit -m "fix(sql): yangilik, shikoyat va reyting ro'yxatlarida ORDER BY inyeksiyasi yopildi

Uchala joyda bir xil nusxa-ko'chirma naqsh edi. Jadval nomlari taxmin
qilinmadi, TableName() metodlaridan olindi — ReportModel \"report\"
qaytaradi, birlikda."
```

---

### Task 4: Qo'llab-quvvatlash dialoglari

Bitta faylda uchta metod, hammasi bir xil modelni so'raydi.

**Files:**
- Modify: `src/infrastructure/persistence/repository/support_dialog_repository_impl.go:57-64, 88-97, 130-139`
- Modify: `test/infrastructure/persistence/repository/sort_whitelist_test.go`

**Interfaces:**
- Consumes: `sorting.New`, `(Whitelist).Resolve`, `(Whitelist).Fields` (Task 1)
- Produces: `repository.SupportDialogSortFields`

Uchala metod ham `SupportDialogModel` ni so'raydi, `TableName()` →
`support_dialogs`. Shuning uchun bitta umumiy ro'yxat yetarli.

- [ ] **Step 1: Yiqiladigan testni qo'shish**

`test/infrastructure/persistence/repository/sort_whitelist_test.go` fayliga:

```go
func TestSupportDialogSortFieldsMatchSwaggerEnums(t *testing.T) {
	// chat_list_handler.go, support_dialog_list_question_handler.go,
	// support_dialog_list_answer_handler.go: Enums(created_at,-created_at)
	assertFields(t, "support_dialogs",
		repository.SupportDialogSortFields.Fields(),
		[]string{"created_at"})
}

// GetChatsByPaging so'rovi last_msg nomli quyi so'rovni JOIN qiladi.
// Ustun to'liq nom bilan yozilishi shu sababli muhim.
func TestSupportDialogSortColumnIsTableQualified(t *testing.T) {
	got, err := repository.SupportDialogSortFields.Resolve("-created_at")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "support_dialogs.created_at DESC" {
		t.Errorf("%q kutilgandi, %q keldi", "support_dialogs.created_at DESC", got)
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/infrastructure/persistence/repository/... -count=1
```

Kutilgan: kompilyatsiya xatosi — `SupportDialogSortFields` aniqlanmagan.

- [ ] **Step 3: Ro'yxatni qo'shish**

`support_dialog_repository_impl.go` — import blokiga
`"slib.uz/src/infrastructure/persistence/sorting"` qo'shing, `"strings"` ni
o'chiring. **`"fmt"` qoladi** — uni 69-qatordagi `fmt.Errorf` ishlatadi.

Birinchi metoddan oldin, fayl darajasida:

```go
// SupportDialogSortFields — qo'llab-quvvatlash dialoglari uchun ruxsat etilgan
// tartiblash maydonlari. Uchala metod ham SupportDialogModel ni so'raydi,
// TableName() "support_dialogs" qaytaradi.
var SupportDialogSortFields = sorting.New("support_dialogs.created_at DESC", map[string]string{
	"created_at": "support_dialogs.created_at",
})
```

- [ ] **Step 4: Uchala metodni tuzatish**

Uchala metodda ham bir xil ish: metod boshida, `var total int64` qatoridan
oldin qo'shing —

```go
	orderExpr, err := SupportDialogSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}
```

— va tartiblash blokini `query = query.Order(orderExpr)` ga almashtiring.

Metodlar va ularning eski standart tartiblari:

| Metod | Eski standart |
|---|---|
| `GetByPaging` | `created_at DESC` |
| `GetByChatID` | `created_at DESC` |
| `GetChatsByPaging` | `support_dialogs.created_at DESC` |

Uchalasi ham endi `support_dialogs.created_at DESC` bo'ladi — natija bir xil,
chunki uchala so'rov ham `support_dialogs` jadvalidan boshlanadi.

Almashtiriladigan blok (uchala joyda ham bir xil ko'rinishda, faqat oxirgi
metodda standart tartib to'liq nom bilan yozilgan):

```go
	if ordering != "" {
		if strings.HasPrefix(ordering, "-") {
			sortField := strings.TrimPrefix(ordering, "-")
			query = query.Order(sortField + " DESC")
		} else {
			query = query.Order(ordering + " ASC")
		}
	} else {
		query = query.Order("created_at DESC")
	}
```

- [ ] **Step 5: Qurish va testlar**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -20
```

Kutilgan: qurish toza, **121** test o'tadi (119 + 2 ta yangi).

- [ ] **Step 6: Uchala nuqta ham tuzatilganini tasdiqlash**

```bash
grep -n 'Order(.*+' src/infrastructure/persistence/repository/support_dialog_repository_impl.go
```

Kutilgan: **hech qanday natija yo'q.** Natija chiqsa — biror metod
o'tkazib yuborilgan.

- [ ] **Step 7: Commit**

```bash
git add src/infrastructure/persistence/repository/support_dialog_repository_impl.go \
        test/infrastructure/persistence/repository/sort_whitelist_test.go
git commit -m "fix(sql): qo'llab-quvvatlash dialoglarida ORDER BY inyeksiyasi yopildi

Bitta faylda uchta metod, hammasida bir xil naqsh. Uchalasi ham
support_dialogs jadvalini so'raydi, shuning uchun bitta umumiy ro'yxat."
```

---

### Task 5: Qolgan ikkita dinamik SQL nuqtasi

Bu ikkisi tartiblash bilan bog'liq emas, lekin keyingi taskdagi qorovul
ularni to'sadi — shuning uchun undan oldin tozalanishi kerak.

**Files:**
- Modify: `src/infrastructure/persistence/repository/published_article_repository_impl.go:390`
- Modify: `src/infrastructure/persistence/repository/journal_application_repository_impl.go:99-112`
- Test: `test/infrastructure/persistence/repository/sort_whitelist_test.go`

**Interfaces:**
- Produces: `repository.ISSNConditions` (`map[string]string`)

- [ ] **Step 1: Yiqiladigan testni qo'shish**

`test/infrastructure/persistence/repository/sort_whitelist_test.go` fayliga:

```go
// ISSNConditions jadvalida ustun nomi emas, BUTUN shart saqlanadi.
// Aks holda repozitoriyda condition + " = ?" birlashtirish bo'lardi va
// qorovul test uni to'xtatardi.
func TestISSNConditionsAreCompleteSQLFragments(t *testing.T) {
	want := map[string]string{
		"issn_paper":  "journals.issn_paper = ?",
		"issn_online": "journals.issn_online = ?",
	}
	if len(repository.ISSNConditions) != len(want) {
		t.Fatalf("%d ta shart kutilgandi, %d keldi", len(want), len(repository.ISSNConditions))
	}
	for key, expected := range want {
		got, ok := repository.ISSNConditions[key]
		if !ok {
			t.Errorf("%q kaliti yo'q", key)
			continue
		}
		if got != expected {
			t.Errorf("%q: %q kutilgandi, %q keldi", key, expected, got)
		}
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/infrastructure/persistence/repository/... -count=1
```

Kutilgan: kompilyatsiya xatosi — `ISSNConditions` aniqlanmagan.

- [ ] **Step 3: `joinOr` yordamchisini ajratish**

`published_article_repository_impl.go` — `GetAll` funksiyasidan oldin qo'shing:

```go
// joinOr faqat shu fayldagi konstanta SQL bo'laklarini birlashtiradi.
// Foydalanuvchi kiritishi bu yerga tushmaydi — u args orqali ? ga bog'lanadi.
//
// Alohida funksiyaga chiqarilgani bejiz emas: chaqiruv joyida bu ifoda
// xavfli koddan tashqi ko'rinishi bilan farq qilmasdi.
func joinOr(conditions []string) string {
	return "(" + strings.Join(conditions, " OR ") + ")"
}
```

390-qatorni:

```go
			query = query.Where("("+strings.Join(conditions, " OR ")+")", args...)
```

almashtiring:

```go
			query = query.Where(joinOr(conditions), args...)
```

- [ ] **Step 4: ISSN shartlar jadvalini qo'shish**

`journal_application_repository_impl.go` — `GetByISSN` funksiyasidan oldin:

```go
// ISSNConditions ISSN turini to'liq SQL shartiga bog'laydi.
//
// Jadvalda ustun nomi emas, BUTUN shart saqlanadi. Aks holda quyida
// condition + " = ?" birlashtirish kerak bo'lardi — ya'ni aynan biz
// yo'q qilayotgan naqsh.
//
// Handler ham issn_type ni tekshiradi (CheckISSNHandler); bu yerdagi
// tekshiruv undan mustaqil va so'rov qurilayotgan joyning o'zida turadi.
var ISSNConditions = map[string]string{
	"issn_paper":  "journals.issn_paper = ?",
	"issn_online": "journals.issn_online = ?",
}
```

`GetByISSN` tanasini almashtiring:

```go
func (this *JournalApplicationRepositoryImpl) GetByISSN(issn string, issnField string) (*entity.JournalApplicationEntity, error) {
	var application models.JournalApplicationModel

	condition, ok := ISSNConditions[issnField]
	if !ok {
		return nil, response.NewFailResponse(400, "invalid ISSN type")
	}

	if err := (this.db().
		Preload("Journal").
		Joins("JOIN journals ON journals.id = journal_applications.journal_id").
		Where(condition, issn).
		Order("created_at desc").
		First(&application)).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.JournalApplicationModelToEntity(&application), nil
}
```

Import blokiga `"slib.uz/src/core/application/response"` qo'shing va `"fmt"` ni
o'chiring — u endi ishlatilmaydi.

- [ ] **Step 5: Qurish va testlar**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -20
```

Kutilgan: qurish toza, **122** test o'tadi (121 + 1 ta yangi).

`response` importi to'qnashuv bersa (fayl allaqachon shu nomdagi biror narsani
import qilgan bo'lsa) — taxallus bering: `resp "slib.uz/src/core/application/response"`
va chaqiruvni `resp.NewFailResponse(...)` ga o'zgartiring.

- [ ] **Step 6: Commit**

```bash
git add src/infrastructure/persistence/repository/published_article_repository_impl.go \
        src/infrastructure/persistence/repository/journal_application_repository_impl.go \
        test/infrastructure/persistence/repository/sort_whitelist_test.go
git commit -m "refactor(sql): qolgan ikkita dinamik SQL nuqtasi tozalandi

Ikkalasi ham xavfsiz edi, lekin xavfli koddan tashqi ko'rinishi bilan
farq qilmasdi. ISSN shartlari jadvalida ustun nomi emas, butun shart
saqlanadi — aks holda yana birlashtirish kerak bo'lardi."
```

---

### Task 6: Qorovul test

Endi `src/` toza bo'lishi kerak. Bu task xavfli **shakl**ni butun daraxtda
taqiqlaydi.

**Files:**
- Create: `test/architecture/raw_sql_test.go`

**Interfaces:**
- Consumes: hech narsa (faqat standart kutubxona)
- Produces: hech narsa

- [ ] **Step 1: Qorovul testini yozish**

Yangi fayl `test/architecture/raw_sql_test.go`:

```go
// Package architecture_test butun kod bazasiga tegishli qoidalarni tekshiradi.
package architecture_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// guardedMethods — GORM'ning birinchi argumenti SQL matni bo'lgan metodlari.
// Qolgan argumentlar ? ga bog'lanadigan qiymatlar, ular xavfsiz.
var guardedMethods = map[string]bool{
	"Order":  true,
	"Where":  true,
	"Or":     true,
	"Not":    true,
	"Having": true,
	"Group":  true,
	"Select": true,
	"Table":  true,
	"Joins":  true,
	"Raw":    true,
	"Exec":   true,
}

// TestNoConcatenatedSQL SQL matniga foydalanuvchi kiritishini birlashtirishni
// taqiqlaydi (CWE-89).
//
// Nima uchun aynan birinchi argument: GORM'da birinchi argument har doim SQL
// matni. Shuning uchun Where("name ILIKE ?", "%"+name+"%") o'tadi — u yerdagi
// + qiymat argumentida va ? orqali parametrlashtiriladi. Where("name = '" +
// name + "'") esa yiqiladi.
//
// Qorovulning chegarasi: u SINTAKTIK. Chaqiruv joyidagi ifodani ko'radi,
// o'zgaruvchining qayerdan kelganini emas. journal_repository_impl.go dagi
// GetJournalStatisticsV2 filterClause ni += bilan quradi va Raw ga
// o'zgaruvchi sifatida beradi — qorovul buni o'tkazib yuboradi. O'sha kod
// xavfsiz, lekin buni qorovul ISBOTLAMAYDI.
//
// Xavfsiz yo'lning haqiqatan xavfsizligini sorting paketining testlari
// kafolatlaydi. Himoya ikki qismdan iborat va ikkalasi ham kerak.
func TestNoConcatenatedSQL(t *testing.T) {
	root := sourceRoot(t)

	var violations []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		file, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return fmt.Errorf("%s: %w", path, parseErr)
		}

		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || !guardedMethods[selector.Sel.Name] {
				return true
			}
			if reason := unsafeSQLExpr(call.Args[0]); reason != "" {
				pos := fset.Position(call.Lparen)
				violations = append(violations, fmt.Sprintf(
					"%s:%d  .%s(...) birinchi argumentida %s",
					relativeTo(root, pos.Filename), pos.Line, selector.Sel.Name, reason))
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("src/ daraxtini o'qib bo'lmadi: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf(
			"SQL matni foydalanuvchi kiritishi bilan birlashtirilgan (%d ta joy).\n"+
				"Tartiblash uchun src/infrastructure/persistence/sorting paketidan\n"+
				"foydalaning; boshqa hollarda SQL bo'lagini konstanta jadvalidan oling.\n\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

// unsafeSQLExpr ifoda ichida SQL matnini yig'ish belgilarini qidiradi.
// Bo'sh qator qaytsa — ifoda xavfsiz.
func unsafeSQLExpr(arg ast.Expr) string {
	reason := ""
	ast.Inspect(arg, func(node ast.Node) bool {
		if reason != "" {
			return false
		}
		switch typed := node.(type) {
		case *ast.BinaryExpr:
			if typed.Op == token.ADD {
				reason = "qatorlarni + bilan birlashtirish bor"
				return false
			}
		case *ast.CallExpr:
			if selector, ok := typed.Fun.(*ast.SelectorExpr); ok {
				pkg, isIdent := selector.X.(*ast.Ident)
				if isIdent && pkg.Name == "fmt" && selector.Sel.Name == "Sprintf" {
					reason = "fmt.Sprintf chaqiruvi bor"
					return false
				}
			}
		}
		return true
	})
	return reason
}

// sourceRoot src/ katalogining yo'lini qaytaradi. Test test/architecture/
// katalogida ishlaydi, shuning uchun ikki pog'ona yuqoriga chiqiladi.
func sourceRoot(t *testing.T) string {
	t.Helper()

	root, err := filepath.Abs(filepath.Join("..", "..", "src"))
	if err != nil {
		t.Fatalf("src/ yo'lini aniqlab bo'lmadi: %v", err)
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		t.Fatalf("src/ katalogi topilmadi: %s", root)
	}
	return root
}

func relativeTo(root, path string) string {
	rel, err := filepath.Rel(filepath.Dir(root), path)
	if err != nil {
		return path
	}
	return rel
}
```

- [ ] **Step 2: Testni ishga tushirish**

```bash
go test ./test/architecture/... -count=1 -v
```

Kutilgan: **PASS.** Oldingi tasklar barcha 17 ta nuqtani tozalagan.

Yiqilsa — chiqishdagi ro'yxatni o'qing. Agar u 1-5-tasklarda ko'rsatilgan
fayllardan bo'lsa, o'sha taskda biror nuqta o'tkazib yuborilgan. Agar
**boshqa** fayl chiqsa — **to'xtang va xabar bering**: bu reja hisobga
olmagan joy, uni o'z bilganicha tuzatmang.

- [ ] **Step 3: Mutatsiya tekshiruvi — qorovul haqiqatan tishlaydimi**

O'tgan test biror narsani ushlashini isbotlamaydi. Zaiflikni vaqtincha
qaytaring — `published_article_repository_impl.go` da
`query = query.Order(orderExpr)` qatoridan **keyin** yangi qator qo'shing:

```go
	query = query.Order(orderExpr)
	query = query.Order(sort + " DESC")   // MUTATSIYA
```

Almashtirmang, **qo'shing**: `orderExpr` ishlatilmay qolsa Go
"declared and not used" xatosi bilan qurishdan bosh tortadi va qorovul
umuman ishga tushmaydi.

```bash
go build ./... 2>&1 | head -5
go test ./test/architecture/... -count=1
```

Kutilgan: qorovul YIQILADI va aynan shu fayl:qatorni ko'rsatadi.

**Yiqilmasa — qorovul yaroqsiz, to'xtang va xabar bering.**

Keyin bekor qiling:

```bash
git checkout src/infrastructure/persistence/repository/published_article_repository_impl.go
go build ./... && go test ./test/architecture/... -count=1
```

- [ ] **Step 4: To'liq testlar**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -20
```

Kutilgan: **123** test o'tadi (122 + 1 ta yangi).

- [ ] **Step 5: Commit**

```bash
git add test/architecture/raw_sql_test.go
git commit -m "test(arch): SQL matnini birlashtirish qorovuli

Naqsh 7 ta faylda takrorlangan edi — ya'ni 8-marta ham ko'chirilishi
mumkin. Bu test uni to'sadi: GORM metodlarining birinchi argumentida
qator birlashtirish yoki fmt.Sprintf bo'lsa, test yiqiladi.

Qorovul sintaktik — nimani isbotlamasligi test izohida yozilgan."
```

---

### Task 7: Yakuniy tekshiruv

**Files:** hech qaysi (faqat tekshirish va hujjat)

- [ ] **Step 1: Generatorlar idempotent ekanini tasdiqlash**

Bu ish wire yoki swagger tomonidan generatsiya qilinadigan hech narsaga
tegmasligi kerak. Avval hamma narsani stage'ga qo'ying, keyin generatorlarni
ishga tushiring:

```bash
git add -A && git status --short
make wire-build
make generate-docs
git diff --stat
```

Kutilgan: `git diff --stat` **bo'sh**. Natija chiqsa — generatsiya qilinadigan
fayl o'zgargan, sababini aniqlang.

- [ ] **Step 2: Butun daraxtda eski naqsh qolmaganini tasdiqlash**

```bash
grep -rn --include=*.go -E "\.Order\([^\"')]" src/ | grep -v "Order(order" | grep -v "Order(this.db"
```

Kutilgan: faqat `journal_editorial_repository_impl.go:52` dagi
`Order(\`"order" ASC\`)` (bu konstanta, xavfsiz) va `legacy_author_repository_impl.go:29`
dagi `Order(this.db().Raw("similarity(full_name, ?) DESC", fullName))`
(parametrlashtirilgan, xavfsiz). Boshqa natija chiqmasligi kerak.

- [ ] **Step 3: To'liq test to'plami**

```bash
go build ./... && go vet ./...
go test ./... -count=1
```

Kutilgan: 123 test, 0 yiqilish, ishchi daraxt toza.

- [ ] **Step 4: Spec bilan solishtirish**

`docs/superpowers/specs/2026-08-14-sql-injection-design.md` ni ochib, har bir
bo'limni bajarilgan ish bilan solishtiring. Farq topilsa — kodni emas, avval
farqni xabar qiling.

- [ ] **Step 5: Shoxni yakunlash**

**REQUIRED SUB-SKILL:** `superpowers:finishing-a-development-branch`

Bazaviy shox: `feature/upload-hardening` (`develop` emas — zanjir
#18 → #21 → shu ish).

---

## Kutilayotgan yakuniy holat

| Ko'rsatkich | Boshlanish | Yakun |
|---|---|---|
| Testlar | 101 | 123 |
| Xavfli SQL nuqtalari (`src/`) | 17 | 0 |
| Yangi bog'liqliklar | — | 0 |
| O'zgargan API kontrakti | — | yo'q (hujjatlashtirilgan qiymatlar avvalgidek) |

## Deploy oldidan (kod ishi emas)

Frontend jamoasi bilan har bir ro'yxat sahifasi qanday `sort` / `sort_by` /
`ordering` qiymatlarini yuborayotganini solishtirib chiqish. Ro'yxatda yo'q
qiymat yuborayotgan sahifa endi `400 invalid sort field` oladi. Bunday qiymat
topilsa — xavfsiz ustun bo'lsa ro'yxatga qo'shiladi, aks holda frontend
tuzatiladi.
