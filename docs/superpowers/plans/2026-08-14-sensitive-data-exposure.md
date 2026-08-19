# Muhim ma'lumotlarni ochiqlanishi (CWE-200) — amalga oshirish rejasi

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** PINFL va tug'ilgan sanani barcha API javoblaridan butunlay olib tashlash, ariza detalidagi telefon/email'ni so'rovchiga qarab redaksiya qilish, va `/journal/{id}/reviewers` ni login orqasiga qo'yib telefonini olib tashlash.

**Architecture:** PINFL/birth_date `json:"-"` teg bilan `UserEntity` dan chiqariladi — bu ularni chiqadigan barcha yo'llarda (ega, affiliation author) bir vaqtda yopadi. Telefon/email so'rovchiga qarab ikki kichik yordamchi bilan redaksiya qilinadi: `UserEntity` uchun (ega, affiliation) va `UserBasicEntity` uchun (reviewer). `go/ast` qorovuli PINFL/birth_date JSON tegining qaytishini to'sadi.

**Tech Stack:** Go 1.25, Echo v4, GORM, PostgreSQL. Testlar — standart `testing`, `encoding/json`, `go/ast`, qo'lda fake'lar.

## Global Constraints

- **Yangi bog'liqlik yo'q.** Faqat standart kutubxona va mavjud paketlar.
- **Testlar ildizdagi `test/` katalogida**, `src/` tuzilishini takrorlaydi, `package <dir>_test`.
- **Core qatlami (`src/core/**`) infratuzilmani import qilmaydi.** Redaksiya yordamchisi core'da.
- **Xatolik grantlamaydi, rad etadi.** Redaksiya deny-by-default: shubha bo'lsa bo'shatadi.
- **Generatsiya qilinadigan fayllar:** `cmd/container/container.go` (wire), `src/entrypoint/presentation/docs/*` (swaggo) — qo'lda tegilmaydi, generator bilan yangilanadi.
- **IDE diagnostikasi bu loyihada har doim eskirgan** — faqat `go build ./...`, `go vet ./...`, `go test ./... -count=1`.
- **Shox:** `feature/sensitive-data`, `feature/sql-injection` ustida. Bazaviy: **131 test o'tadi**, `go build` toza.
- Izohlar va commit xabarlari o'zbek tilida.

## Aniqlashtirilgan haqiqat (spec'dan aniqroq)

`ApplicationResponseEntity` da maxfiy ma'lumot uch yo'ldan chiqadi (mapper tekshirildi):

| Yo'l | Tur | PINFL | Telefon | Email | Hozirgi holat |
|---|---|---|---|---|---|
| `.User` | `*UserEntity` | bor | bor | bor | xom uzatiladi |
| `.Article.CoAuthors[].User` | `nil` | — | — | — | `slimApplicationCoAuthors` tozalaydi |
| `.Article.CoAuthorsWithAffiliation[].Author.User` | `*UserEntity` | bor | bor | bor | xom uzatiladi |
| `.ReviewStages[].Reviewer` | `*UserBasicEntity` | **yo'q** | bor | **yo'q** | `UserEntityToBasic` |

- PINFL/birth_date faqat `UserEntity` da → `json:"-"` ega va affiliation yo'llarini birdan yopadi.
- Telefon: `UserEntity` (ega, affiliation) va `UserBasicEntity` (reviewer) — ikki tur, ikki yordamchi.
- `UserBasicEntity` da email yo'q, `Pin` allaqachon izohga olingan (`user_basic_entity.go`).

---

## Fayl tuzilishi

**Yaratiladi:**

| Fayl | Mas'uliyati |
|---|---|
| `src/core/application/usecase/permissionusecases/redact_contact.go` | `RedactUserContact`, `RedactBasicContact`, `RedactApplicationContacts` |
| `test/core/application/usecase/permissionusecases/redact_contact_test.go` | Yordamchilar xatti-harakati (11 test) |
| `test/core/domain/entity/sensitive_json_test.go` | `UserEntity`/`ReviewerEntity` marshal — Pin/BirthDate/telefon chiqmasligi |
| `test/architecture/sensitive_json_test.go` | Qorovul: pin/birth_date JSON tegi taqiqi |

**O'zgartiriladi:**

| Fayl | O'zgarish |
|---|---|
| `src/core/domain/entity/user_entity.go` | `Pin`, `BirthDate` → `json:"-"` |
| `src/core/domain/entity/reviewer_entity.go` | `PhoneNumber` → `json:"-"` |
| `src/entrypoint/presentation/handlers/journal_config/schema/journal_config_response.go` | `CreatorResponse.Pin` olib tashlash |
| `src/core/application/usecase/article_applications_usecases/application_detail_usecase.go` | `Execute(requester, id)`, redaksiya |
| `src/entrypoint/presentation/handlers/article_application/application_detail_handler.go` | so'rovchini uzatish |
| `src/entrypoint/presentation/groups/journal_group.go` | reviewers → `AuthenticatedPermission` |

---

## Boshlashdan oldin: bazaviy holat

```bash
git branch --show-current     # feature/sensitive-data
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -5
```

131 test o'tishi kerak. O'tmasa — to'xtang.

---

### Task 1: PINFL va birth_date'ni JSON'dan chiqarish + qorovul

Auditning yuragi. `json:"-"` bilan PINFL/birth_date barcha `UserEntity` yo'llarida
bir vaqtda yopiladi.

**Files:**
- Modify: `src/core/domain/entity/user_entity.go:13,21`
- Modify: `src/entrypoint/presentation/handlers/journal_config/schema/journal_config_response.go`
- Test: `test/core/domain/entity/sensitive_json_test.go` (create)
- Test: `test/architecture/sensitive_json_test.go` (create)

**Interfaces:**
- Produces: `UserEntity` endi `pin`/`birth_date` JSON kalitlarini chiqarmaydi.

- [ ] **Step 1: `UserEntity` marshal testini yozish**

Yangi fayl `test/core/domain/entity/sensitive_json_test.go`:

```go
package entity_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"slib.uz/src/core/domain/entity"
)

// PINFL va tug'ilgan sana hech qanday API javobida chiqmasligi kerak (CWE-200).
func TestUserEntityDoesNotExposePinOrBirthDate(t *testing.T) {
	pin := "12345678901234"
	birth := time.Now()
	u := &entity.UserEntity{
		ID:          1,
		Pin:         &pin,
		BirthDate:   &birth,
		PhoneNumber: "+998901234567",
		Email:       "a@b.uz",
	}

	raw, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("marshal xatosi: %v", err)
	}
	out := string(raw)

	if strings.Contains(out, "\"pin\"") {
		t.Error("javobda \"pin\" kaliti bor, bo'lmasligi kerak")
	}
	if strings.Contains(out, "\"birth_date\"") {
		t.Error("javobda \"birth_date\" kaliti bor, bo'lmasligi kerak")
	}
	if strings.Contains(out, pin) {
		t.Error("javobda PINFL qiymati bor")
	}
	// Telefon bu darajada qoladi — u so'rovchiga qarab use case'da redaksiya qilinadi.
	if !strings.Contains(out, "phone_number") {
		t.Error("phone_number kaliti bo'lishi kerak edi (bu darajada redaksiya yo'q)")
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/core/domain/entity/... -run TestUserEntityDoesNotExposePinOrBirthDate -count=1
```

Kutilgan: FAIL — hozir `pin` va `birth_date` kalitlari chiqadi.

- [ ] **Step 3: `UserEntity` teglarini o'zgartirish**

`src/core/domain/entity/user_entity.go` — ikki qatorni o'zgartiring:

```go
	Pin            *string                  `json:"-"`
```

va

```go
	BirthDate      *time.Time               `json:"-"`
```

- [ ] **Step 4: Testni ishga tushirib, o'tishini tasdiqlash**

```bash
go test ./test/core/domain/entity/... -run TestUserEntityDoesNotExposePinOrBirthDate -count=1
```

Kutilgan: PASS.

- [ ] **Step 5: `journal_config` CreatorResponse.Pin'ni olib tashlash**

`src/entrypoint/presentation/handlers/journal_config/schema/journal_config_response.go` —
`CreatorResponse` struct'idan `Pin` maydonini o'chiring, va uni to'ldiruvchi blokni
(`:56-64` atrofidagi `if e.Creator != nil { pin := ... }`) soddalashtiring:

```go
	if e.Creator != nil {
		resp.Creator = &CreatorResponse{
			ID:       e.Creator.ID,
			FullName: e.Creator.FullName,
		}
	}
```

`CreatorResponse` struct e'lonidan `Pin` maydonini ham o'chiring. Ishlatilmay
qolgan o'zgaruvchi (`pin`) bo'lsa, `go build` ko'rsatadi.

- [ ] **Step 6: Qorovul testini yozish**

Yangi fayl `test/architecture/sensitive_json_test.go`:

```go
package architecture_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// forbiddenJSONKeys — bu maydonlar hech qachon JSON javobida chiqmasligi kerak.
// Telefon ataylab yo'q: u UserBasicEntity da (jurnal a'zolari uchun) qoladi va
// use case darajasida so'rovchiga qarab redaksiya qilinadi.
var forbiddenJSONKeys = map[string]bool{
	"pin":        true,
	"pinfl":      true,
	"birth_date": true,
}

// TestNoSensitiveJSONTags maxfiy maydonlarning JSON teg orqali oshkor bo'lishini
// taqiqlaydi (CWE-200). Kelajakda kimdir json:"pin" qo'shsa, test yiqiladi.
func TestNoSensitiveJSONTags(t *testing.T) {
	roots := []string{
		sourceDir(t, "core", "domain", "entity"),
		sourceDir(t, "entrypoint", "presentation", "handlers"),
	}

	var violations []string
	for _, root := range roots {
		walkStructTags(t, root, func(file string, line int, jsonKey string) {
			if forbiddenJSONKeys[jsonKey] {
				violations = append(violations, fmt.Sprintf("%s:%d  json:%q", file, line, jsonKey))
			}
		})
	}

	if len(violations) > 0 {
		t.Errorf(
			"maxfiy maydon JSON teg orqali oshkor bo'lgan (%d ta joy).\n"+
				"PINFL va tug'ilgan sana API javoblarida chiqmasligi kerak;\n"+
				"json:\"-\" ishlating.\n\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

// walkStructTags har bir struct maydonining json teg kalitini callback'ga beradi.
func walkStructTags(t *testing.T, root string, fn func(file string, line int, jsonKey string)) {
	t.Helper()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		f, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			return fmt.Errorf("%s: %w", path, perr)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			field, ok := n.(*ast.Field)
			if !ok || field.Tag == nil {
				return true
			}
			tagValue := strings.Trim(field.Tag.Value, "`")
			jsonTag := reflect.StructTag(tagValue).Get("json")
			if jsonTag == "" {
				return true
			}
			key := strings.Split(jsonTag, ",")[0]
			if key == "" || key == "-" {
				return true
			}
			fn(path, fset.Position(field.Pos()).Line, key)
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("daraxtni o'qib bo'lmadi: %v", err)
	}
}

func sourceDir(t *testing.T, parts ...string) string {
	t.Helper()
	all := append([]string{"..", "..", "src"}, parts...)
	dir, err := filepath.Abs(filepath.Join(all...))
	if err != nil {
		t.Fatalf("yo'lni aniqlab bo'lmadi: %v", err)
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Fatalf("katalog topilmadi: %s", dir)
	}
	return dir
}
```

- [ ] **Step 7: Qorovulni ishga tushirish**

```bash
go build ./... && go test ./test/architecture/... -run TestNoSensitiveJSONTags -count=1
```

Kutilgan: PASS (3-5-qadamlar `pin`/`birth_date` teglarini olib tashlagan).

Yiqilsa — chiqishdagi `fayl:qator` ni o'qing. Agar u 3-5-qadamdagi fayllardan
bo'lsa, teg to'liq o'zgartirilmagan. Boshqa fayl chiqsa — **to'xtang va xabar
bering**: reja hisobga olmagan maxfiy teg bor.

- [ ] **Step 8: Qorovulning o'z tishini tekshirish**

Vaqtincha `user_entity.go` da `Pin` ni `json:"pin"` ga qaytaring:

```bash
go test ./test/architecture/... -run TestNoSensitiveJSONTags -count=1
```

Kutilgan: YIQILADI, `user_entity.go` ni ko'rsatadi. **Yiqilmasa — qorovul
yaroqsiz, to'xtang.** Keyin `json:"-"` ga qaytaring va qayta tasdiqlang.

- [ ] **Step 9: To'liq testlar va commit**

```bash
go build ./... && go vet ./... && go test ./... -count=1 2>&1 | tail -5
```

Kutilgan: 133 test (131 + 2 yangi test funksiyasi).

```bash
git add src/core/domain/entity/user_entity.go \
        src/entrypoint/presentation/handlers/journal_config/schema/journal_config_response.go \
        test/core/domain/entity/sensitive_json_test.go \
        test/architecture/sensitive_json_test.go
git commit -m "fix(cwe-200): PINFL va tug'ilgan sana API javoblaridan olib tashlandi

UserEntity.Pin va .BirthDate endi json:\"-\" — bu ularni chiqadigan
barcha yo'llarda (ariza egasi, affiliation muallif) bir vaqtda yopadi.
journal_config CreatorResponse.Pin ham olib tashlandi.

go/ast qorovuli pin/pinfl/birth_date JSON tegining qaytishini to'sadi."
```

---

### Task 2: Aloqa ma'lumotlari redaksiya yordamchilari

`UserEntity` va `UserBasicEntity` uchun ikki kichik yordamchi. Core qatlamida,
`permissionusecases` da (`IsAdmin` yonida).

**Files:**
- Create: `src/core/application/usecase/permissionusecases/redact_contact.go`
- Test: `test/core/application/usecase/permissionusecases/redact_contact_test.go`

**Interfaces:**
- Consumes: `entity.UserEntity`, `entity.UserBasicEntity`, `entity.ApplicationResponseEntity`
- Produces:
  - `permissionusecases.RedactUserContact(user *entity.UserEntity, requesterID uint, isPrivileged bool)`
  - `permissionusecases.RedactBasicContact(user *entity.UserBasicEntity, requesterID uint, isPrivileged bool)`
  - `permissionusecases.RedactApplicationContacts(resp *entity.ApplicationResponseEntity, requesterID uint, isPrivileged bool)`

- [ ] **Step 1: Yiqiladigan testni yozish**

Yangi fayl `test/core/application/usecase/permissionusecases/redact_contact_test.go`:

```go
package permissionusecases_test

import (
	"testing"

	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
)

func TestRedactUserContactBlanksForStranger(t *testing.T) {
	u := &entity.UserEntity{ID: 5, PhoneNumber: "+998901112233", Email: "a@b.uz"}
	permissionusecases.RedactUserContact(u, 99, false)
	if u.PhoneNumber != "" || u.Email != "" {
		t.Errorf("begona uchun telefon/email bo'sh bo'lishi kerak, %q / %q keldi", u.PhoneNumber, u.Email)
	}
}

func TestRedactUserContactKeepsForOwner(t *testing.T) {
	u := &entity.UserEntity{ID: 5, PhoneNumber: "+998901112233", Email: "a@b.uz"}
	permissionusecases.RedactUserContact(u, 5, false)
	if u.PhoneNumber == "" || u.Email == "" {
		t.Error("egaga o'z aloqa ma'lumoti ko'rinishi kerak edi")
	}
}

func TestRedactUserContactKeepsForPrivileged(t *testing.T) {
	u := &entity.UserEntity{ID: 5, PhoneNumber: "+998901112233", Email: "a@b.uz"}
	permissionusecases.RedactUserContact(u, 99, true)
	if u.PhoneNumber == "" || u.Email == "" {
		t.Error("admin uchun aloqa ma'lumoti ko'rinishi kerak edi")
	}
}

func TestRedactUserContactNilSafe(t *testing.T) {
	permissionusecases.RedactUserContact(nil, 1, false) // panic bermasligi kerak
}

func TestRedactBasicContactBlanksForStranger(t *testing.T) {
	u := &entity.UserBasicEntity{ID: 5, PhoneNumber: "+998901112233"}
	permissionusecases.RedactBasicContact(u, 99, false)
	if u.PhoneNumber != "" {
		t.Errorf("begona uchun telefon bo'sh bo'lishi kerak, %q keldi", u.PhoneNumber)
	}
}

func TestRedactBasicContactKeepsForOwner(t *testing.T) {
	u := &entity.UserBasicEntity{ID: 5, PhoneNumber: "+998901112233"}
	permissionusecases.RedactBasicContact(u, 5, false)
	if u.PhoneNumber == "" {
		t.Error("egaga o'z telefoni ko'rinishi kerak edi")
	}
}

func TestRedactBasicContactNilSafe(t *testing.T) {
	permissionusecases.RedactBasicContact(nil, 1, false) // panic bermasligi kerak
}

// RedactApplicationContacts javobning uchala maxfiy yo'lini birdan tozalaydi:
// ega (.User), affiliation muallif (.Article.CoAuthorsWithAffiliation[].Author.User),
// va reviewer (.ReviewStages[].Reviewer). Fake repo kerak emas — response entity
// to'g'ridan-to'g'ri quriladi.
func TestRedactApplicationContactsBlanksStranger(t *testing.T) {
	resp := &entity.ApplicationResponseEntity{
		User: &entity.UserEntity{ID: 10, PhoneNumber: "+998900000010", Email: "owner@x.uz"},
		Article: &entity.ArticleInputEntity{
			CoAuthorsWithAffiliation: []*entity.ArticleAuthorAffiliationEntity{
				{Author: &entity.AuthorEntity{User: &entity.UserEntity{ID: 20, PhoneNumber: "+998900000020", Email: "aff@x.uz"}}},
			},
		},
		ReviewStages: []*entity.ReviewStageResponseEntity{
			{Reviewer: &entity.UserBasicEntity{ID: 30, PhoneNumber: "+998900000030"}},
		},
	}

	permissionusecases.RedactApplicationContacts(resp, 999, false)

	if resp.User.PhoneNumber != "" || resp.User.Email != "" {
		t.Error("ega telefon/email bo'sh bo'lishi kerak")
	}
	if aff := resp.Article.CoAuthorsWithAffiliation[0].Author.User; aff.PhoneNumber != "" || aff.Email != "" {
		t.Error("affiliation muallif telefon/email bo'sh bo'lishi kerak")
	}
	if resp.ReviewStages[0].Reviewer.PhoneNumber != "" {
		t.Error("reviewer telefon bo'sh bo'lishi kerak")
	}
}

func TestRedactApplicationContactsKeepsOwnerRecord(t *testing.T) {
	// So'rovchi = ega (ID 10). Eganing o'z yozuvi ko'rinadi, begonalarniki emas.
	resp := &entity.ApplicationResponseEntity{
		User: &entity.UserEntity{ID: 10, PhoneNumber: "+998900000010", Email: "owner@x.uz"},
		ReviewStages: []*entity.ReviewStageResponseEntity{
			{Reviewer: &entity.UserBasicEntity{ID: 30, PhoneNumber: "+998900000030"}},
		},
	}

	permissionusecases.RedactApplicationContacts(resp, 10, false)

	if resp.User.PhoneNumber == "" {
		t.Error("egaga o'z telefoni ko'rinishi kerak edi")
	}
	if resp.ReviewStages[0].Reviewer.PhoneNumber != "" {
		t.Error("begona reviewer telefoni bo'sh bo'lishi kerak")
	}
}

func TestRedactApplicationContactsNilSafe(t *testing.T) {
	permissionusecases.RedactApplicationContacts(nil, 1, false) // panic bermasligi kerak
	// Article nil va ReviewStages nil holatlar ham xavfsiz:
	permissionusecases.RedactApplicationContacts(&entity.ApplicationResponseEntity{}, 1, false)
}
```

**Reja bosqichida tasdiqlang:** `ReviewStageResponseEntity.Reviewer` maydon turi
`*UserBasicEntity` (mapper `UserEntityToBasic` qaytaradi). Yuqoridagi test shunga
tayanadi. Agar `go build` uni `*UserEntity` deb ko'rsatsa, testda ham, funksiyada
ham `RedactUserContact` ishlating.

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/core/application/usecase/permissionusecases/... -run Redact -count=1
```

Kutilgan: kompilyatsiya xatosi — funksiyalar aniqlanmagan.

- [ ] **Step 3: Yordamchilarni yozish**

Yangi fayl `src/core/application/usecase/permissionusecases/redact_contact.go`:

```go
package permissionusecases

import "slib.uz/src/core/domain/entity"

// RedactUserContact so'rovchi egasi yoki privilegiyalangan (admin) bo'lmasa
// UserEntity dagi aloqa ma'lumotlarini bo'shatadi.
//
// PINFL va tug'ilgan sana bu yerda kerak emas — ular json:"-" bilan baribir
// chiqmaydi. Bu yordamchi faqat telefon va email uchun, chunki ular ba'zi
// kontekstlarda (o'z yozuvi, admin) ko'rinishi kerak.
func RedactUserContact(user *entity.UserEntity, requesterID uint, isPrivileged bool) {
	if user == nil || user.ID == requesterID || isPrivileged {
		return
	}
	user.PhoneNumber = ""
	user.Email = ""
}

// RedactBasicContact UserBasicEntity uchun aynan shu qoidani qo'llaydi.
// UserBasicEntity da email yo'q, shuning uchun faqat telefon bo'shatiladi.
func RedactBasicContact(user *entity.UserBasicEntity, requesterID uint, isPrivileged bool) {
	if user == nil || user.ID == requesterID || isPrivileged {
		return
	}
	user.PhoneNumber = ""
}

// RedactApplicationContacts ariza detali javobining uchala maxfiy yo'lini
// tozalaydi. UserEntity yo'llarida (ega, affiliation muallif) telefon+email,
// UserBasicEntity yo'lida (reviewer) telefon. PINFL/birth_date json:"-" bilan
// baribir chiqmaydi.
func RedactApplicationContacts(resp *entity.ApplicationResponseEntity, requesterID uint, isPrivileged bool) {
	if resp == nil {
		return
	}

	RedactUserContact(resp.User, requesterID, isPrivileged)

	if resp.Article != nil {
		for _, aff := range resp.Article.CoAuthorsWithAffiliation {
			if aff != nil && aff.Author != nil {
				RedactUserContact(aff.Author.User, requesterID, isPrivileged)
			}
		}
	}

	for _, stage := range resp.ReviewStages {
		if stage != nil {
			RedactBasicContact(stage.Reviewer, requesterID, isPrivileged)
		}
	}
}
```

**Reja bosqichida tasdiqlang:** import blokiga faqat `entity` kerak. Agar
`ReviewStageResponseEntity.Reviewer` `*UserEntity` bo'lsa (kutilmaydi),
`RedactBasicContact` o'rniga `RedactUserContact` ishlating.

- [ ] **Step 4: Testlarni ishga tushirish**

```bash
go build ./... && go test ./test/core/application/usecase/permissionusecases/... -run Redact -count=1 -v
```

Kutilgan: 11 ta test PASS.

- [ ] **Step 5: Mutatsiya tekshiruvi**

Vaqtincha `RedactUserContact` da `user.ID == requesterID` ni `user.ID != requesterID`
ga o'zgartiring:

```bash
go test ./test/core/application/usecase/permissionusecases/... -run Redact -count=1
```

Kutilgan: `TestRedactUserContactBlanksForStranger` va `KeepsForOwner` yiqiladi.
**Yiqilmasa — to'xtang.** Keyin qaytaring.

- [ ] **Step 6: Commit**

```bash
git add src/core/application/usecase/permissionusecases/redact_contact.go \
        test/core/application/usecase/permissionusecases/redact_contact_test.go
git commit -m "feat(cwe-200): aloqa ma'lumotlari redaksiya yordamchilari

RedactUserContact (UserEntity: telefon+email) va RedactBasicContact
(UserBasicEntity: telefon). So'rovchi egasi yoki admin bo'lmasa bo'shatadi —
deny-by-default."
```

---

### Task 3: Ariza detali use case va handler

`application_detail_usecase` so'rovchini oladi va Task 2 ning
`RedactApplicationContacts` funksiyasini chaqiradi. Redaksiya mantiqi Task 2 da
allaqachon testlangani uchun bu task integratsiyaga — imzo o'zgarishi va so'rovchini
uzatishga — qaratilgan.

**Files:**
- Modify: `src/core/application/usecase/article_applications_usecases/application_detail_usecase.go`
- Modify: `src/entrypoint/presentation/handlers/article_application/application_detail_handler.go`

**Interfaces:**
- Consumes: `permissionusecases.RedactApplicationContacts`, `permissionusecases.IsAdmin` (Task 2 + mavjud)
- Produces: `ApplicationDetailUsecase.Execute(requester *entity.UserBasicEntity, applicationId uint)`

- [ ] **Step 1: Use case imzosini o'zgartirish va redaksiya chaqirish**

`application_detail_usecase.go` — `Execute` ni o'zgartiring:

```go
func (this *ApplicationDetailUsecase) Execute(requester *entity.UserBasicEntity, applicationId uint) (*entity.ApplicationResponseEntity, error) {
	application, err := this.repository.GetByIDWithRelations(applicationId)
	if err != nil {
		return nil, err
	}

	references, err := this.referenceRepo.GetListByArticleID(application.ArticleID)
	if err != nil {
		return nil, err
	}
	application.Article.References = references

	response := mapper.ApplicationEntityToResponse(application)

	var requesterID uint
	if requester != nil {
		requesterID = requester.ID
	}
	permissionusecases.RedactApplicationContacts(response, requesterID, permissionusecases.IsAdmin(requester))

	return response, nil
}
```

Import blokiga `"slib.uz/src/core/application/usecase/permissionusecases"` qo'shing.

- [ ] **Step 2: Handler'ni so'rovchini uzatishga o'zgartirish**

`application_detail_handler.go` — `Handle` ni o'zgartiring:

```go
func (this *ApplicationDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	applicationId, err := context.GetIntPathParam(ctx, "applicationId")
	if err != nil {
		return err
	}

	app, err := this.uc.Execute(c.User, uint(applicationId))
	if err != nil {
		return err
	}
	return c.JsonResponse(http.StatusOK, app)
}
```

`c.User` — `*entity.UserBasicEntity` (IDOR ishida tasdiqlangan). Endpoint
`AuthenticatedPermission` bilan, shuning uchun `c.User` odatda bor; `RedactUserContact`
`nil` so'rovchida ham xavfsiz (requesterID=0, hech kim ega bo'lmaydi, hammasi bo'shaladi).

- [ ] **Step 4: Testni to'ldirib, o'tkazish**

1-qadamdagi `buildAppWithContacts` va `TestApplicationDetailRedactsStrangerContacts`
ni haqiqiy struct va fake repo bilan to'ldiring (`t.Skip` olib tashlanadi).

```bash
go build ./... && go test ./test/core/application/usecase/article_applications_usecases/... -count=1 -v
```

Kutilgan: PASS — begona so'rovchi uchun barcha telefon/email bo'sh.

- [ ] **Step 5: Butun to'plam**

```bash
go build ./... && go vet ./... && go test ./... -count=1 2>&1 | tail -5
```

Kutilgan: qurish toza. Boshqa chaqiruvchilar (`myapplication` — u boshqa use case,
bu imzo o'zgarishi unga tegmaydi; wire generatsiyasi) buzilmasligini tekshiring.
Agar `Execute` ning boshqa chaqiruvchisi bo'lsa (`go build` ko'rsatadi), uni ham
moslang.

- [ ] **Step 6: Commit**

```bash
git add src/core/application/usecase/article_applications_usecases/application_detail_usecase.go \
        src/entrypoint/presentation/handlers/article_application/application_detail_handler.go \
        test/core/application/usecase/article_applications_usecases/application_detail_redaction_test.go
git commit -m "fix(cwe-200): ariza detalida telefon/email so'rovchiga qarab redaksiya

Auditning aynan endpointi (/article-application/detail). Javobdagi har bir
maxfiy foydalanuvchi yozuvi (ariza egasi, affiliation muallif, reviewer)
begona so'rovchi uchun telefon/email'siz qaytadi. So'rovchi egasi yoki
admin bo'lsa ko'radi."
```

---

### Task 4: Reviewer endpointini yopish va telefonini olib tashlash

**Files:**
- Modify: `src/core/domain/entity/reviewer_entity.go:8`
- Modify: `src/entrypoint/presentation/groups/journal_group.go:48`
- Test: `test/core/domain/entity/sensitive_json_test.go` (Task 1 fayliga qo'shiladi)

**Interfaces:**
- Consumes: `permissions.AuthenticatedPermission` (mavjud)

- [ ] **Step 1: Reviewer telefon marshal testini qo'shish**

`test/core/domain/entity/sensitive_json_test.go` fayliga qo'shing:

```go
func TestReviewerEntityDoesNotExposePhone(t *testing.T) {
	r := &entity.ReviewerEntity{ID: 1, FullName: "A B", PhoneNumber: "+998901234567"}
	raw, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal xatosi: %v", err)
	}
	if strings.Contains(string(raw), "phone_number") {
		t.Error("reviewer javobida phone_number bor, bo'lmasligi kerak")
	}
	if strings.Contains(string(raw), "+998901234567") {
		t.Error("reviewer javobida telefon qiymati bor")
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

```bash
go test ./test/core/domain/entity/... -run TestReviewerEntityDoesNotExposePhone -count=1
```

Kutilgan: FAIL — hozir `phone_number` chiqadi.

- [ ] **Step 3: `ReviewerEntity.PhoneNumber` ni yopish**

`src/core/domain/entity/reviewer_entity.go:8`:

```go
	PhoneNumber string           `json:"-"`
```

- [ ] **Step 4: Reviewer endpointini login orqasiga qo'yish**

`src/entrypoint/presentation/groups/journal_group.go:48`:

```go
	group.GET("/:journalId/reviewers", this.reviewersList.Handle, permissions.AuthenticatedPermission)
```

Import blokida `permissions` bor-yo'qligini tekshiring; yo'q bo'lsa qo'shing
(`"slib.uz/src/entrypoint/presentation/interceptor/permissions"`).

- [ ] **Step 5: Testlar**

```bash
go build ./... && go vet ./... && go test ./... -count=1 2>&1 | tail -5
```

Kutilgan: qurish toza, `TestReviewerEntityDoesNotExposePhone` PASS, umumiy 145.

- [ ] **Step 6: Commit**

```bash
git add src/core/domain/entity/reviewer_entity.go \
        src/entrypoint/presentation/groups/journal_group.go \
        test/core/domain/entity/sensitive_json_test.go
git commit -m "fix(cwe-200): reviewer ro'yxati login talab qiladi, telefon olib tashlandi

/journal/{id}/reviewers hozir anonim ochiq edi va telefon ro'yxatini
qaytarardi. Endi AuthenticatedPermission bilan va ReviewerEntity.PhoneNumber
json:\"-\"."
```

---

### Task 5: Yakuniy tekshiruv

**Files:** hech qaysi (tekshirish va hujjat)

- [ ] **Step 1: Generatorlar idempotent**

```bash
git add -A && git status --short
make wire-build
make generate-docs
git diff --stat
```

Kutilgan: `git diff --stat` bo'sh. `application_detail_usecase.Execute` imzosi
o'zgargani uchun wire buni qayta generatsiya qilishi mumkin — agar diff chiqsa,
u faqat generatsiya natijasi bo'lsa commit qiling; qo'lda o'zgartirmang.

- [ ] **Step 2: PINFL/birth_date hech qayerda JSON'ga chiqmasligini tasdiqlash**

```bash
grep -rn 'json:"pin"\|json:"pinfl"\|json:"birth_date"' src/
```

Kutilgan: hech qanday natija.

- [ ] **Step 3: `.Pin`/`.BirthDate` JSON round-trip qilmasligini tasdiqlash**

```bash
grep -rn '\.Pin\b\|\.BirthDate\b' src/ | grep -viE '_test|json:'
```

Har bir natijани ko'rib chiqing — hech biri `json.Marshal` orqali PINFL'ni
tashqariga chiqarmasligi kerak (kutilgan: faqat ichki qidiruv `GetByPin` va
mapper'lar). Yangi oshkor yo'l topilsa — **to'xtang va xabar bering**.

- [ ] **Step 4: To'liq to'plam**

```bash
go build ./... && go vet ./... && go test ./... -count=1
```

Kutilgan: 145 test, 0 yiqilish, ishchi daraxt toza.

- [ ] **Step 5: Spec bilan solishtirish**

`docs/superpowers/specs/2026-08-14-sensitive-data-exposure-design.md` ni ochib
har bir bo'limni bajarilgan ish bilan solishtiring. Spec §3.2 "har bir UserEntity"
degan — amalda reviewer `UserBasicEntity` (bu reja aniqlashtirilgan haqiqat
jadvalida qayd etilgan). Farq topilsa — kodni emas, avval farqni xabar qiling.

- [ ] **Step 6: Shoxni yakunlash**

**REQUIRED SUB-SKILL:** `superpowers:finishing-a-development-branch`

Bazaviy shox: `feature/sql-injection` (zanjir #18 → #21 → #22 → shu ish).

---

## Kutilayotgan yakuniy holat

| Ko'rsatkich | Boshlanish | Yakun |
|---|---|---|
| Testlar | 131 | 145 |
| PINFL JSON teglari (`src/`) | bor | 0 |
| Anonim ochiq PII endpoint | `/journal/{id}/reviewers` | yopilgan |
| Yangi bog'liqlik | — | 0 |

## Deploy oldidan (kod ishi emas)

Frontend jamoasi bilan uchta o'zgarishni solishtirish: (1) PINFL/birth_date barcha
javoblardan yo'qoldi, (2) `/journal/{id}/reviewers` anonim → login va telefonsiz,
(3) ariza detalida begona uchun telefon/email bo'sh.
