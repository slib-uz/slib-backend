# IDOR / resursga bog'langan ruxsat nazorati — amalga oshirish rejasi

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Foydalanuvchi so'rov parametrini o'zgartirib begona resursga kira olmasin (CWE-639).

**Architecture:** Ruxsat tekshiruvi usecase qatlamida — loyihada allaqachon o'rnatilgan uslub. Markaziy `JournalMemberPermissionUseCase` repozitoriy orqali jurnalning haqiqiy egasini aniqlaydi, shu sababli uning imzosi `(bool, error)` ga o'zgaradi va 20 ta chaqiruv joyi yangilanadi. Ikki zaif endpoint o'z usecase'ida `user` ni qabul qilib tekshiruv o'tkazadi.

**Tech Stack:** Go 1.25, Echo v4, GORM, google/wire (DI), standart `testing` paketi.

**Spec:** [2026-08-12-idor-access-control-design.md](../specs/2026-08-12-idor-access-control-design.md)

## Global Constraints

- Modul nomi `slib.uz`; barcha import yo'llari shundan boshlanadi.
- **Ruxsat tekshiruvi hech qachon xatoni `true` ga aylantirmaydi** (fail-closed). Xato yuzaga kelsa `(false, err)` qaytariladi va yuqoriga uzatiladi.
- **Barcha ko'rsatkichli maydonlar ishlatilishdan oldin `nil` ga tekshiriladi** — `role.JournalID`, `role.PublisherID`. `permissions/role_based_permission.go:53` da aynan shu yo'qligi sababli panic beradigan kod bor; bu yerda takrorlanmasin.
- Core qatlamida (`src/core/**`) loglash uchun faqat `github.com/labstack/gommon/log`. `zap` faqat `src/infrastructure/logger` ichida.
- `@inject` izohli konstruktor o'zgargach **`make wire-build`** ishga tushiriladi.
- Swagger izohi o'zgargan bo'lsa **`make generate-docs`**.
- Yangi bog'liqlik (test kutubxonasi va boshqa) **qo'shilmaydi**. Testlar standart `testing` va qo'lda yozilgan soxta obyektlar bilan.
- Katta interfeyslar uchun soxta obyektda interfeysning o'zi struct ichiga joylanadi (embedded) — faqat kerakli metodlar qayta yoziladi.
- Har bir task oxirida `go build ./...` va `go test ./... -count=1` toza o'tishi shart.

**Bog'liqlik eslatmasi:** `Makefile` da `test` target'i bo'lmasligi mumkin — uni
[2026-08-12-session-revocation.md](2026-08-12-session-revocation.md) rejasi kiritadi. Agar
yo'q bo'lsa, Task 1 uni qo'shadi; bor bo'lsa, o'sha qadam o'tkazib yuboriladi.

## File Structure

**Yaratiladigan fayllar:**

| Fayl | Mas'uliyati |
|---|---|
| `src/core/application/usecase/permissionusecases/is_admin.go` | Yagona admin tekshiruvi |
| `src/core/application/usecase/permissionusecases/is_admin_test.go` | Uning testlari |
| `src/core/application/usecase/permissionusecases/journal_member_permission_usecase_test.go` | Markaziy tekshiruv testlari |
| `src/core/application/usecase/userusecases/user_detail_usecase_test.go` | O'zi/admin testlari |
| `src/core/application/usecase/journalusecases/journal_members_list_usecase_test.go` | Jurnal a'zolari testlari |

**O'zgartiriladigan fayllar:**

| Fayl | O'zgarish |
|---|---|
| `permissionusecases/journal_member_permission_usecase.go` | Repozitoriy, imzo, publisher-admin mantig'i |
| 20 ta chaqiruv joyi (Task 2 jadvali) | `(bool, error)` ga moslashish |
| `userusecases/user_detail_usecase.go` | `user` parametri va tekshiruv |
| `journalusecases/journal_members_list_usecase.go` | `user` parametri va tekshiruv |
| `handlers/user/user_detail_handler.go` | `c.User` ni uzatish |
| `handlers/journal/journal_members_list_handler.go` | `c.User` ni uzatish |

---

## Task 1: `IsAdmin` yordamchisi

Loyihada uch xil "admin" tushunchasi bor va ular kelishmaydi: `permissions.AdminPermission`
faqat `IsAdmin` bayrog'ini, `UserListUseCase` ikkalasini, `JournalMemberPermissionUseCase`
faqat `RoleAdmin` rolini tekshiradi. Yangi kod yagona funksiyaga tayanadi.

**Files:**
- Create: `src/core/application/usecase/permissionusecases/is_admin.go`
- Test: `src/core/application/usecase/permissionusecases/is_admin_test.go`
- Modify: `Makefile` (agar `test` target'i yo'q bo'lsa)

**Interfaces:**
- Consumes: hech narsa (birinchi task)
- Produces: `permissionusecases.IsAdmin(user *entity.UserBasicEntity) bool` — Task 4 ishlatadi

- [ ] **Step 1: Testni yoz**

`src/core/application/usecase/permissionusecases/is_admin_test.go`:

```go
package permissionusecases

import (
	"testing"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

func TestIsAdminReturnsFalseForNilUser(t *testing.T) {
	if IsAdmin(nil) {
		t.Fatal("nil foydalanuvchi admin deb topildi")
	}
}

func TestIsAdminHonoursFlag(t *testing.T) {
	user := &entity.UserBasicEntity{ID: 1, IsAdmin: true}
	if !IsAdmin(user) {
		t.Fatal("IsAdmin bayrog'i e'tiborga olinmadi")
	}
}

func TestIsAdminHonoursRole(t *testing.T) {
	user := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleAdmin}},
	}
	if !IsAdmin(user) {
		t.Fatal("RoleAdmin roli e'tiborga olinmadi")
	}
}

func TestIsAdminRejectsOrdinaryUser(t *testing.T) {
	user := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleSecretary}},
	}
	if IsAdmin(user) {
		t.Fatal("oddiy foydalanuvchi admin deb topildi")
	}
}

func TestIsAdminSkipsNilRoles(t *testing.T) {
	user := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{nil, {Role: enum.RoleAdmin}},
	}
	if !IsAdmin(user) {
		t.Fatal("nil rol butun ro'yxatni buzdi")
	}
}
```

- [ ] **Step 2: Test yiqilishini tasdiqla**

Run: `go test ./src/core/application/usecase/permissionusecases/ -v`
Expected: FAIL — kompilyatsiya xatosi `undefined: IsAdmin`

- [ ] **Step 3: Funksiyani yoz**

`src/core/application/usecase/permissionusecases/is_admin.go`:

```go
package permissionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

// IsAdmin foydalanuvchining super-admin ekanini aniqlaydi.
//
// Loyihada admin ikki xil ifodalanadi: UserEntity dagi IsAdmin bayrog'i va
// RoleAdmin roli. Ikkalasi ham qabul qilinadi, shunda tekshiruv qaysi
// mexanizm ishlatilganidan qat'i nazar bir xil javob beradi.
func IsAdmin(user *entity.UserBasicEntity) bool {
	if user == nil {
		return false
	}

	if user.IsAdmin {
		return true
	}

	for _, role := range user.Roles {
		if role != nil && role.Role == enum.RoleAdmin {
			return true
		}
	}

	return false
}
```

- [ ] **Step 4: Testlar o'tishini tasdiqla**

Run: `go test ./src/core/application/usecase/permissionusecases/ -v`
Expected: PASS — 5 ta test

- [ ] **Step 5: `Makefile` da `test` target'i borligini tekshir**

Run: `grep -q '^test:' Makefile && echo "bor" || echo "yo'q"`

Agar `yo'q` bo'lsa, `Makefile` oxiriga qo'sh:

```makefile
test:
	@echo "Running tests..."
	@go test ./... -count=1
```

Agar `bor` bo'lsa, bu qadamni o'tkazib yubor.

- [ ] **Step 6: Commit**

```bash
git add src/core/application/usecase/permissionusecases/is_admin.go src/core/application/usecase/permissionusecases/is_admin_test.go Makefile
git commit -m "feat(authz): yagona IsAdmin yordamchisi

Loyihada admin uch xil tekshirilardi. Yangi kod yagona funksiyaga
tayanadi; mavjud uchta joy ataylab tegilmaydi."
```

---

## Task 2: `Execute` imzosini `(bool, error)` ga o'tkazish — xatti-harakat o'zgarmaydi

Bu task **sof mexanik**. Publisher-admin qisqa yo'li o'z joyida qoladi va hozirgi kabi
`true` qaytaradi — haqiqiy tuzatish Task 3 da. Shu sababli bu taskda hech qanday xatti-harakat
o'zgarmaydi va uni alohida ko'rib chiqish oson.

Imzo o'zgargani uchun loyiha 20 ta chaqiruv joyi yangilanmaguncha kompilyatsiya bo'lmaydi —
shuning uchun hammasi bitta taskda.

**Files:**
- Modify: `src/core/application/usecase/permissionusecases/journal_member_permission_usecase.go`
- Modify: quyidagi jadvaldagi 20 ta fayl

**Interfaces:**
- Consumes: hech narsa
- Produces: `(*JournalMemberPermissionUseCase).Execute(userRoles []*entity.UserRoleEntity, journalID uint) (bool, error)` — Task 3 va Task 5 shunga tayanadi

- [ ] **Step 1: Imzoni o'zgartir**

`src/core/application/usecase/permissionusecases/journal_member_permission_usecase.go` —
`Execute` metodini almashtir (mantiq **o'zgarmaydi**, faqat qaytish tipi):

```go
func (this *JournalMemberPermissionUseCase) Execute(userRoles []*entity.UserRoleEntity, journalID uint) (bool, error) {
	for _, role := range userRoles {

		if role.Role == enum.RoleAdmin {
			return true, nil
		}

		// TODO: Check if the user has a role in the journal and if that role is one of the allowed roles for the action.
		if role.Role == enum.RolePublisherAdmin {
			return true, nil
		}

		isMemberRole := role.Role == enum.RoleChiefEditor || role.Role == enum.RoleSecretary
		if role.JournalID != nil && *role.JournalID == journalID && isMemberRole {
			return true, nil
		}
	}
	return false, nil
}
```

`TODO` ataylab qoldiriladi — u Task 3 da olib tashlanadi.

- [ ] **Step 2: 16 ta odatiy chaqiruv joyini yangila**

Har birida quyidagi almashtirish bajariladi. Namuna
(`editionusecases/edition_create_usecase.go:21`):

```go
// hozir
if !this.memberPermission.Execute(user.Roles, edition.JournalID) {
	return response.PermissionDeniedError
}

// bo'ladi
allowed, err := this.memberPermission.Execute(user.Roles, edition.JournalID)
if err != nil {
	return err
}
if !allowed {
	return response.PermissionDeniedError
}
```

**Muhim:** mavjud xato qaytarish qatori (`response.PermissionDeniedError` yoki
`response.NewFailResponse(403, "...")`) **o'zgarmaydi** — faqat shart qayta quriladi.
Agar funksiyada allaqachon `err` nomli o'zgaruvchi bo'lsa, `allowed, err = ...` shaklida
qayta ishlating yoki `permErr` deb nomlang.

| # | Fayl | Qator | Qabul qiluvchi | Jurnal ID ifodasi |
|---|---|---|---|---|
| 1 | `editionusecases/edition_create_usecase.go` | 21 | `memberPermission` | `edition.JournalID` |
| 2 | `editionusecases/edition_update_usecase.go` | 28 | `memberPermission` | `edition.JournalID` |
| 3 | `editionusecases/edition_delete_usecase.go` | 37 | `memberPermission` | `edition.JournalID` |
| 4 | `editionusecases/edition_attach_articles_usecase.go` | 28 | `memberPermission` | `edition.JournalID` |
| 5 | `editionusecases/edition_detach_articles_usecase.go` | 34 | `memberPermission` | `edition.JournalID` |
| 6 | `journaleditorialusecases/journal_editorial_create_usecase.go` | 21 | `memberPermission` | `editorial.JournalID` |
| 7 | `journaleditorialusecases/journal_editorial_update_usecase.go` | 25 | `memberPermission` | `existing.JournalID` |
| 8 | `journaleditorialusecases/journal_editorial_delete_usecase.go` | 25 | `memberPermission` | `existing.JournalID` |
| 9 | `journalnewsusecases/journal_news_create_usecase.go` | 21 | `memberPermission` | `news.JournalID` |
| 10 | `journalnewsusecases/journal_news_update_usecase.go` | 25 | `memberPermission` | `existing.JournalID` |
| 11 | `journalnewsusecases/journal_news_delete_usecase.go` | 25 | `memberPermission` | `existing.JournalID` |
| 12 | `article_applications_usecases/peer_review_usecase.go` | 45 | `journalMemberPermission` | `reviewStage.Application.JournalID` |
| 13 | `article_applications_usecases/technical_reivew_usecase.go` | 43 | `journalMemberPermission` | `reviewStage.Application.JournalID` |
| 14 | `article_applications_usecases/application_publish_usecase.go` | 51 | `permission` | `application.Article.JournalID` |
| 15 | `articleusecases/direct_article_create_usecase.go` | 49 | `permissionUseCase` | `journalID` |
| 16 | `roiusecase/article_publish_roi_usecase.go` | 42 | `permission` | `application.Article.JournalID` |

Barchasi `src/core/application/usecase/` ostida.

- [ ] **Step 3: `presigned_url_usecase.go` ni yangila (nostandart — ijobiy tarmoq)**

`src/core/application/usecase/uploadusecases/presigned_url_usecase.go:53-57` — 55-qatorni
almashtir:

```go
func (this *PresignedURLUseCase) checkAccess(user *entity.UserBasicEntity, journalID, articleID uint) (bool, error) {

	isMember, err := this.journalMemberPermission.Execute(user.Roles, journalID)
	if err != nil {
		return false, err
	}
	if isMember {
		return true, nil
	}
```

Qolgan tarmoqlar (`authorPermission`, `articlePurchasedPermission`) tegilmaydi.

- [ ] **Step 4: `applications_list_usecase.go` ni yangila (nostandart — ijobiy tarmoq)**

`src/core/application/usecase/article_applications_usecases/applications_list_usecase.go:56-75` —
`checkAccess` metodini almashtir:

```go
func (this *ApplicationsListUseCase) checkAccess(user *entity2.UserBasicEntity, journalID uint) error {

	if user.IsAdmin {
		return nil
	}

	isMember, err := this.journalMemberPermission.Execute(user.Roles, journalID)
	if err != nil {
		return err
	}
	if isMember {
		return nil
	}

	journal, err := this.journalRepository.FindByID(journalID)
	if err != nil {
		return err
	}

	if this.publisherAdminPermission.Execute(user.Roles, journal.PublisherID) {
		return nil
	}
	return response.PermissionDeniedError
}
```

Eslatma: Task 3 dan keyin quyidagi `publisherAdminPermission` tarmog'i ortiqcha bo'lib qoladi,
chunki `journalMemberPermission` o'zi nashriyot egaligini tekshiradi. **Uni olib tashlamang** —
xatti-harakat bir xil qoladi va bu spec doirasidan tashqari.

- [ ] **Step 5: `extend_deadline_usecase.go` ni yangila (nostandart — qo'shma shart)**

`src/core/application/usecase/deadlineusecases/extend_deadline_usecase.go:61-72` —
`hasPermission` metodini almashtir:

```go
func (this *ExtendDeadlineUseCase) hasPermission(user *entity2.UserBasicEntity, journalID uint, deadlineType enum.DeadlineType) error {

	if deadlineType == enum.DeadlineTypeReviewDeadline && this.chiefEditorPermission.Execute(user.Roles, journalID) {
		return nil
	}

	// Qo'shma shart ikkiga bo'lindi: aks holda && qisqa tutashuvi bilan
	// xatoni qayta ishlash chalkash bo'lardi va deadlineType mos kelmaganda
	// ham DB so'rovi ketishi mumkin edi.
	if deadlineType == enum.DeadlineTypeResubmitDeadline {
		isMember, err := this.journalMemberPermission.Execute(user.Roles, journalID)
		if err != nil {
			return err
		}
		if isMember {
			return nil
		}
	}

	return response.PermissionDeniedError
}
```

- [ ] **Step 6: `application_reviewer_permission.go` ni yangila (nostandart — to'g'ridan-to'g'ri qaytish)**

`src/core/application/usecase/permissionusecases/application_reviewer_permission.go:23` —
`Execute` metodining oxirgi qatorini almashtir:

```go
func (this *ApplicationReviewerPermissionUseCase) Execute(userRoles []*entity.UserRoleEntity, applicationID uint) (bool, error) {
	application, err := this.applicationRepository.FindByID(applicationID)
	if err != nil {
		return false, err
	}
	return this.journalMemberPermission.Execute(userRoles, application.JournalID)
}
```

Oldingi `, nil` olib tashlanadi — endi ichki chaqiruv o'zi `(bool, error)` qaytaradi.

- [ ] **Step 7: Qurilishni tekshir**

Run: `go build ./...`
Expected: xatosiz (exit 0)

Agar `declared and not used: err` yoki `err redeclared` xatolari chiqsa, o'sha funksiyada
allaqachon `err` bor — `allowed, err = ...` (`:=` emas) ishlating yoki nomni `permErr` ga
o'zgartiring.

- [ ] **Step 8: Mavjud testlar hali o'tishini tasdiqla**

Run: `go test ./... -count=1`
Expected: PASS — bu task xatti-harakatni o'zgartirmaydi

- [ ] **Step 9: Commit**

```bash
git add src/core/application/usecase/
git commit -m "refactor(authz): JournalMemberPermission.Execute (bool, error) qaytaradi

Sof mexanik o'zgarish, xatti-harakat bir xil. Keyingi commitda
publisher-admin tekshiruviga repozitoriy qidiruvi qo'shiladi va
u xato qaytarishi mumkin bo'ladi."
```

---

## Task 3: Publisher-admin tekshiruvini tuzatish

Ishning markazi. Hozir `RolePublisherAdmin` roliga ega istalgan foydalanuvchi **istalgan**
jurnalga kiradi. Bu funksiya 20 ta joyda chaqilgani uchun bitta tuzatish hammasini yopadi.

**Files:**
- Modify: `src/core/application/usecase/permissionusecases/journal_member_permission_usecase.go`
- Test: `src/core/application/usecase/permissionusecases/journal_member_permission_usecase_test.go`

**Interfaces:**
- Consumes: `Execute(...) (bool, error)` imzosi (Task 2)
- Produces: `NewJournalMemberPermissionUseCase(journalRepository repository.JournalRepository) *JournalMemberPermissionUseCase` — konstruktor endi parametr oladi

- [ ] **Step 1: Testni yoz**

`src/core/application/usecase/permissionusecases/journal_member_permission_usecase_test.go`:

```go
package permissionusecases

import (
	"errors"
	"testing"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeJournalRepo — JournalRepository interfeysi embedded, faqat kerakli metod yoziladi.
type fakeJournalRepo struct {
	repository.JournalRepository
	publisherID uint
	err         error
	calls       int
}

func (f *fakeJournalRepo) GetPublisherIdByJournalId(journalID uint) (uint, error) {
	f.calls++
	return f.publisherID, f.err
}

func journalRole(role enum.UserRole, journalID uint) *entity.UserRoleEntity {
	id := journalID
	return &entity.UserRoleEntity{Role: role, JournalID: &id}
}

func publisherRole(publisherID uint) *entity.UserRoleEntity {
	id := publisherID
	return &entity.UserRoleEntity{Role: enum.RolePublisherAdmin, PublisherID: &id}
}

func TestChiefEditorAllowedOnOwnJournal(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{journalRole(enum.RoleChiefEditor, 7)}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if !allowed {
		t.Fatal("bosh muharrir o'z jurnaliga kira olmadi")
	}
	if repo.calls != 0 {
		t.Fatalf("keraksiz DB so'rovi ketdi: %d", repo.calls)
	}
}

func TestChiefEditorDeniedOnOtherJournal(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{journalRole(enum.RoleChiefEditor, 7)}, 99)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if allowed {
		t.Fatal("bosh muharrir begona jurnalga kirdi")
	}
}

// Zaiflikning o'zi: hozir istalgan nashriyot admini istalgan jurnalga kiradi.
func TestPublisherAdminDeniedOnOtherPublishersJournal(t *testing.T) {
	repo := &fakeJournalRepo{publisherID: 500} // jurnal 500-nashriyotga tegishli
	uc := NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{publisherRole(300)}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if allowed {
		t.Fatal("begona nashriyot admini jurnalga kirdi — IDOR ochiq")
	}
}

func TestPublisherAdminAllowedOnOwnPublishersJournal(t *testing.T) {
	repo := &fakeJournalRepo{publisherID: 300}
	uc := NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{publisherRole(300)}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if !allowed {
		t.Fatal("o'z nashriyoti jurnaliga publisher admin kira olmadi")
	}
}

func TestAdminAlwaysAllowedWithoutLookup(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{{Role: enum.RoleAdmin}}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if !allowed {
		t.Fatal("admin rad etildi")
	}
	if repo.calls != 0 {
		t.Fatalf("admin uchun keraksiz DB so'rovi ketdi: %d", repo.calls)
	}
}

func TestNoLookupWhenUserHasNoPublisherAdminRole(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := NewJournalMemberPermissionUseCase(repo)

	_, err := uc.Execute([]*entity.UserRoleEntity{journalRole(enum.RoleSecretary, 99)}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if repo.calls != 0 {
		t.Fatalf("publisher-admin roli yo'q, lekin DB so'rovi ketdi: %d", repo.calls)
	}
}

func TestRepositoryErrorIsPropagatedAndDenies(t *testing.T) {
	repo := &fakeJournalRepo{err: errors.New("db down")}
	uc := NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{publisherRole(300)}, 7)
	if err == nil {
		t.Fatal("repozitoriy xatosi yutib yuborildi")
	}
	if allowed {
		t.Fatal("xato holatida ruxsat berildi — fail-closed buzildi")
	}
}

func TestNilRolesAreSkipped(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{nil, {Role: enum.RoleAdmin}}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if !allowed {
		t.Fatal("nil rol butun ro'yxatni buzdi")
	}
}

func TestPublisherRoleWithNilPublisherIDDoesNotPanic(t *testing.T) {
	repo := &fakeJournalRepo{publisherID: 300}
	uc := NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{{Role: enum.RolePublisherAdmin}}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if allowed {
		t.Fatal("PublisherID nil bo'lgan rol ruxsat berdi")
	}
}
```

- [ ] **Step 2: Test yiqilishini tasdiqla**

Run: `go test ./src/core/application/usecase/permissionusecases/ -run TestPublisherAdminDenied -v`
Expected: FAIL — `NewJournalMemberPermissionUseCase` hozir parametr qabul qilmaydi
(`too many arguments`)

- [ ] **Step 3: Usecase'ni qayta yoz**

`src/core/application/usecase/permissionusecases/journal_member_permission_usecase.go` —
to'liq yangi mazmuni:

```go
package permissionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalMemberPermissionUseCase struct {
	journalRepository repository.JournalRepository
}

// @inject
func NewJournalMemberPermissionUseCase(journalRepository repository.JournalRepository) *JournalMemberPermissionUseCase {
	return &JournalMemberPermissionUseCase{journalRepository: journalRepository}
}

// Execute foydalanuvchining jurnal boshqaruviga kirish huquqini tekshiradi.
//
// Ruxsat beriladigan holatlar:
//   - RoleAdmin — har doim;
//   - RoleChiefEditor yoki RoleSecretary — faqat o'sha jurnalda;
//   - RolePublisherAdmin — faqat jurnal egasi bo'lgan nashriyotda.
//
// Nashriyot egaligini aniqlash uchun DB so'rovi faqat foydalanuvchida
// RolePublisherAdmin roli bo'lsa va arzon tekshiruvlar natija bermagan
// bo'lsagina bajariladi.
func (this *JournalMemberPermissionUseCase) Execute(userRoles []*entity.UserRoleEntity, journalID uint) (bool, error) {
	hasPublisherAdminRole := false

	for _, role := range userRoles {
		if role == nil {
			continue
		}

		if role.Role == enum.RoleAdmin {
			return true, nil
		}

		isMemberRole := role.Role == enum.RoleChiefEditor || role.Role == enum.RoleSecretary
		if isMemberRole && role.JournalID != nil && *role.JournalID == journalID {
			return true, nil
		}

		if role.Role == enum.RolePublisherAdmin && role.PublisherID != nil {
			hasPublisherAdminRole = true
		}
	}

	if !hasPublisherAdminRole {
		return false, nil
	}

	ownerPublisherID, err := this.journalRepository.GetPublisherIdByJournalId(journalID)
	if err != nil {
		return false, err
	}

	for _, role := range userRoles {
		if role == nil {
			continue
		}
		if role.Role == enum.RolePublisherAdmin && role.PublisherID != nil && *role.PublisherID == ownerPublisherID {
			return true, nil
		}
	}

	return false, nil
}
```

- [ ] **Step 4: Testlar o'tishini tasdiqla**

Run: `go test ./src/core/application/usecase/permissionusecases/ -v`
Expected: PASS — Task 1 dagi 5 ta va bu taskdagi 9 ta test

- [ ] **Step 5: Qurilish va DI**

Run: `make wire-build && go build ./... && go test ./... -count=1`
Expected: hammasi xatosiz

- [ ] **Step 6: Commit**

```bash
git add src/core/application/usecase/permissionusecases/ cmd/container/container.go
git commit -m "fix(authz): publisher admin faqat o'z nashriyoti jurnallariga kiradi

TODO qoldirilgan qisqa yo'l tufayli istalgan nashriyot admini istalgan
jurnalning boshqaruv ma'lumotlariga kirardi. Funksiya 20 ta joyda
chaqiriladi, shuning uchun bitta tuzatish hammasini yopadi.

Jurnal egasini aniqlash uchun DB so'rovi faqat publisher-admin roli
bo'lganda bajariladi."
```

---

## Task 4: `GET /api/user/{id}` — o'zi yoki admin

**Files:**
- Modify: `src/core/application/usecase/userusecases/user_detail_usecase.go`
- Modify: `src/entrypoint/presentation/handlers/user/user_detail_handler.go:43`
- Test: `src/core/application/usecase/userusecases/user_detail_usecase_test.go`

**Interfaces:**
- Consumes: `permissionusecases.IsAdmin(user *entity.UserBasicEntity) bool` (Task 1)
- Produces: `(*UserDetailUseCase).Execute(user *entity2.UserBasicEntity, id uint) (*entity2.UserDetailEntity, error)`

- [ ] **Step 1: Testni yoz**

`src/core/application/usecase/userusecases/user_detail_usecase_test.go`:

```go
package userusecases

import (
	"testing"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeUserRepo — UserRepository interfeysi embedded, faqat kerakli metod yoziladi.
type fakeUserRepo struct {
	repository.UserRepository
	calls int
}

func (f *fakeUserRepo) GetDetailByID(id uint) (*entity.UserDetailEntity, error) {
	f.calls++
	return &entity.UserDetailEntity{}, nil
}

func TestUserCanReadOwnProfile(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := NewUserDetailUseCase(repo)

	_, err := uc.Execute(&entity.UserBasicEntity{ID: 42}, 42)
	if err != nil {
		t.Fatalf("o'z profilini o'qiy olmadi: %v", err)
	}
}

func TestUserCannotReadOtherProfile(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := NewUserDetailUseCase(repo)

	_, err := uc.Execute(&entity.UserBasicEntity{ID: 42}, 99)
	if err == nil {
		t.Fatal("begona profil ochildi — IDOR")
	}
	if repo.calls != 0 {
		t.Fatalf("rad etilgan so'rov uchun DB'ga borildi: %d", repo.calls)
	}
}

func TestAdminCanReadAnyProfile(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := NewUserDetailUseCase(repo)

	admin := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleAdmin}},
	}

	if _, err := uc.Execute(admin, 99); err != nil {
		t.Fatalf("admin begona profilni o'qiy olmadi: %v", err)
	}
}

func TestNilUserIsRejected(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := NewUserDetailUseCase(repo)

	if _, err := uc.Execute(nil, 42); err == nil {
		t.Fatal("nil foydalanuvchi qabul qilindi")
	}
}
```

- [ ] **Step 2: Test yiqilishini tasdiqla**

Run: `go test ./src/core/application/usecase/userusecases/ -v`
Expected: FAIL — `too many arguments in call to uc.Execute`

- [ ] **Step 3: Usecase'ni yangila**

`src/core/application/usecase/userusecases/user_detail_usecase.go` — to'liq yangi mazmuni:

```go
package userusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserDetailUseCase struct {
	repository repository.UserRepository
}

// @inject
func NewUserDetailUseCase(repository repository.UserRepository) *UserDetailUseCase {
	return &UserDetailUseCase{repository: repository}
}

// Execute foydalanuvchi profilini qaytaradi.
//
// Profil shaxsga doir ma'lumotlarni (PINFL, tug'ilgan sana) o'z ichiga oladi,
// shuning uchun uni faqat egasining o'zi va admin ko'ra oladi. Ommaviy muallif
// ma'lumoti uchun /api/author/... va /api/users/find endpointlari mavjud.
func (this *UserDetailUseCase) Execute(user *entity2.UserBasicEntity, id uint) (*entity2.UserDetailEntity, error) {
	if user == nil {
		return nil, response.UnauthorizedError
	}

	if user.ID != id && !permissionusecases.IsAdmin(user) {
		return nil, response.PermissionDeniedError
	}

	return this.repository.GetDetailByID(id)
}
```

- [ ] **Step 4: Testlar o'tishini tasdiqla**

Run: `go test ./src/core/application/usecase/userusecases/ -v`
Expected: PASS — 4 ta test

- [ ] **Step 5: Handler'ni yangila**

`src/entrypoint/presentation/handlers/user/user_detail_handler.go:43` — chaqiruvni almashtir:

```go
	result, err := this.uc.Execute(user, uint(id))
```

Handler'dagi mavjud `if user == nil { return echo.ErrForbidden }` tekshiruvi qoladi.

- [ ] **Step 6: Qurilish va testlar**

Run: `go build ./... && go test ./... -count=1`
Expected: hammasi xatosiz

- [ ] **Step 7: Commit**

```bash
git add src/core/application/usecase/userusecases/ src/entrypoint/presentation/handlers/user/user_detail_handler.go
git commit -m "fix(authz): /user/{id} faqat o'ziga va adminga ochiq

Ilgari istalgan autentifikatsiyalangan foydalanuvchi begona profilni,
jumladan PINFL va pasport ma'lumotlarini o'qiy olardi."
```

---

## Task 5: `GET /api/journal-manage/{journalId}/members` — jurnal boshqaruvi

**Files:**
- Modify: `src/core/application/usecase/journalusecases/journal_members_list_usecase.go`
- Modify: `src/entrypoint/presentation/handlers/journal/journal_members_list_handler.go:35`
- Test: `src/core/application/usecase/journalusecases/journal_members_list_usecase_test.go`

**Interfaces:**
- Consumes: `NewJournalMemberPermissionUseCase(journalRepository repository.JournalRepository)` va `Execute(...) (bool, error)` (Task 2, Task 3)
- Produces: `(*JournalMembersListUsecase).Execute(user *entity.UserBasicEntity, journalID uint) ([]*entity.UserRoleWithBasicUserEntity, error)`

- [ ] **Step 1: Testni yoz**

`src/core/application/usecase/journalusecases/journal_members_list_usecase_test.go`:

```go
package journalusecases

import (
	"testing"

	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeUserRoleRepo — UserRoleRepository interfeysi embedded.
type fakeUserRoleRepo struct {
	repository.UserRoleRepository
	calls int
}

func (f *fakeUserRoleRepo) GetByJournalID(journalID uint) ([]*entity.UserRoleEntity, error) {
	f.calls++
	return []*entity.UserRoleEntity{}, nil
}

// fakeJournalRepo — JournalRepository interfeysi embedded.
type fakeJournalRepo struct {
	repository.JournalRepository
	publisherID uint
}

func (f *fakeJournalRepo) GetPublisherIdByJournalId(journalID uint) (uint, error) {
	return f.publisherID, nil
}

func chiefEditorOf(journalID uint) *entity.UserBasicEntity {
	id := journalID
	return &entity.UserBasicEntity{
		ID:    42,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleChiefEditor, JournalID: &id}},
	}
}

func newUseCase(roleRepo *fakeUserRoleRepo) *JournalMembersListUsecase {
	permission := permissionusecases.NewJournalMemberPermissionUseCase(&fakeJournalRepo{})
	return NewJournalMembersListUsecase(roleRepo, permission)
}

func TestJournalMemberCanListOwnJournalMembers(t *testing.T) {
	roleRepo := &fakeUserRoleRepo{}
	uc := newUseCase(roleRepo)

	if _, err := uc.Execute(chiefEditorOf(7), 7); err != nil {
		t.Fatalf("o'z jurnali a'zolarini ko'ra olmadi: %v", err)
	}
}

func TestForeignJournalMembersAreDenied(t *testing.T) {
	roleRepo := &fakeUserRoleRepo{}
	uc := newUseCase(roleRepo)

	if _, err := uc.Execute(chiefEditorOf(7), 99); err == nil {
		t.Fatal("begona jurnal a'zolari ochildi — IDOR")
	}
	if roleRepo.calls != 0 {
		t.Fatalf("rad etilgan so'rov uchun DB'ga borildi: %d", roleRepo.calls)
	}
}

func TestNilUserIsRejected(t *testing.T) {
	uc := newUseCase(&fakeUserRoleRepo{})

	if _, err := uc.Execute(nil, 7); err == nil {
		t.Fatal("nil foydalanuvchi qabul qilindi")
	}
}
```

- [ ] **Step 2: Test yiqilishini tasdiqla**

Run: `go test ./src/core/application/usecase/journalusecases/ -v`
Expected: FAIL — `too many arguments in call to NewJournalMembersListUsecase`

- [ ] **Step 3: Usecase'ni yangila**

`src/core/application/usecase/journalusecases/journal_members_list_usecase.go` — to'liq
yangi mazmuni:

```go
package journalusecases

import (
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalMembersListUsecase struct {
	repository       repository.UserRoleRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalMembersListUsecase(
	repository repository.UserRoleRepository,
	memberPermission *permissionusecases.JournalMemberPermissionUseCase,
) *JournalMembersListUsecase {
	return &JournalMembersListUsecase{repository: repository, memberPermission: memberPermission}
}

// Execute jurnal tahririyat a'zolari ro'yxatini qaytaradi.
//
// Javob a'zolarning shaxsiy maydonlarini o'z ichiga oladi, shuning uchun
// ro'yxat faqat jurnal boshqaruviga ochiq: bosh muharrir, kotib, jurnal
// egasi bo'lgan nashriyot admini va super-admin.
func (this *JournalMembersListUsecase) Execute(
	user *entity.UserBasicEntity, journalID uint,
) ([]*entity.UserRoleWithBasicUserEntity, error) {

	if user == nil {
		return nil, response.UnauthorizedError
	}

	allowed, err := this.memberPermission.Execute(user.Roles, journalID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, response.PermissionDeniedError
	}

	_members, err := this.repository.GetByJournalID(journalID)
	if err != nil {
		return nil, err
	}

	members := make([]*entity.UserRoleEntity, len(_members))
	copy(members, _members)

	return mapper.UserRoleEntityListToWithBasicUserEntityList(members), nil
}
```

- [ ] **Step 4: Testlar o'tishini tasdiqla**

Run: `go test ./src/core/application/usecase/journalusecases/ -v`
Expected: PASS — 3 ta test

- [ ] **Step 5: Handler'ni yangila**

`src/entrypoint/presentation/handlers/journal/journal_members_list_handler.go:35` — chaqiruvni
almashtir:

```go
	members, err := this.uc.Execute(c.User, uint(journalId))
```

Swagger izohiga `403` javobini qo'sh (24-25 qatorlar orasiga):

```go
// @Failure 403 {object} response.Response
```

- [ ] **Step 6: Qurilish, DI, testlar va hujjatlar**

Run: `make wire-build && go build ./... && go test ./... -count=1 && make generate-docs`
Expected: hammasi xatosiz

- [ ] **Step 7: Commit**

```bash
git add src/core/application/usecase/journalusecases/ src/entrypoint/presentation/handlers/journal/journal_members_list_handler.go cmd/container/container.go src/entrypoint/presentation/docs/
git commit -m "fix(authz): jurnal a'zolari ro'yxati faqat jurnal boshqaruviga ochiq

Ilgari istalgan autentifikatsiyalangan foydalanuvchi journalId ni
o'zgartirib istalgan jurnal tahririyatini ko'ra olardi."
```

---

## Yakuniy tekshiruv (kod emas, qo'lda)

Spec'ning 10-bo'limidagi muvaffaqiyat mezonlari. Ilovani ishga tushirib
(`make run`), turli rollardagi tokenlar bilan tekshiring:

- [ ] `GET /api/user/{begona_id}` → `403`
- [ ] `GET /api/user/{o'z_id}` → `200`
- [ ] Admin tokeni bilan `GET /api/user/{istalgan_id}` → `200`
- [ ] A jurnal bosh muharriri `GET /api/journal-manage/{B}/members` → `403`
- [ ] A nashriyoti admini B nashriyoti jurnalining tahririyatiga → `403` (ilgari `200` edi)
- [ ] O'z nashriyoti jurnalida publisher admin ilgarigidek ishlaydi → `200`
