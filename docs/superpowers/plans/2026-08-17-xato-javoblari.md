# Xato javoblarining izchilligi — amalga oshirish rejasi

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mijozga bo'sh `{"data":{}}` qaytaradigan 3 ta handler'ni tuzatish va bu shaklning qaytib kelishini AST qorovul bilan to'sish; `fmt.Println`/`log.Println` orqali yozilayotgan 6 ta joyni tozalash (biri — shaxsiy ma'lumot sizishi, CWE-200).

**Architecture:** Ish ikki mustaqil yo'nalishdan iborat. Birinchisi — **HTTP javob tanasi**: `JsonResponse(status, err)` chaqiruvida `err` `error` interfeysi bo'lgani va `*response.Response` dan boshqa xato turlarida eksport qilinadigan maydon bo'lmagani uchun `encoding/json` uni `{}` deb serializatsiya qiladi; yechim — `err.Error()`, qorovul esa shaklni butun `src/` da taqiqlaydi. Ikkinchisi — **loglash**: loyihada `zap` (`*logger.AsyncLogger`) mavjud va wire orqali ulanadi, lekin ROI oqimidagi kod hali ham stdout'ga yozadi. Yechim — mavjud logger'ni konstruktorga inyeksiya qilish.

**Tech Stack:** Go 1.25, Echo v4, `go.uber.org/zap` v1.27.1, Google wire (`wiregenx` + `wire`). Testlar — standart `testing`, `go/parser`, `go/ast`, `go.uber.org/zap/zaptest/observer`.

## Global Constraints

- **Yangi bog'liqlik yo'q.** `zap` va `zaptest/observer` allaqachon `go.mod` da.
- **Testlar ildizdagi `test/` katalogida**, `src/` tuzilishini takrorlaydi, paket nomi `<dir>_test`.
- **Generatsiya qilinadigan fayllarga qo'lda tegilmaydi:** `cmd/container/container.go` va `di/wire/provider.go` — faqat `make wire-build` orqali. `src/entrypoint/presentation/docs/*` (swaggo) bu ishda umuman o'zgarmaydi.
- **Xatoni yutish semantikasi o'zgarmaydi.** ROI navbatiga yuborish nosozligi hozir ham so'rovni yiqitmaydi; bu reja faqat `fmt.Println` ni tuzilgan log bilan almashtiradi. Xulq-atvorni o'zgartirish — alohida ish.
- **IDE diagnostikasi bu loyihada har doim eskirgan.** Faqat `go build ./...`, `go vet ./...`, `go test ./... -count=1` natijasiga ishoning.
- **Shox:** `feature/error-responses`. Bazaviy holat: `go build ./...` toza, **174 test o'tadi, 0 FAIL**.
- Izohlar, test xabarlari va commit xabarlari o'zbek tilida — loyihaning mavjud odati.
- Har bir task alohida commit.
- **`logs/async/` endi FON ISHLARI logi** (ROI navbat + push), "worker jarayoni logi" emas — `InitHttpServer` ham, worker jarayoni ham `logger.NewAsyncLogger()` quradi va ikkalasi ham shu katalogga yozadi. HTTP so'rov/javob logi bundan alohida, `logs/http/` da qoladi.
  - ROI yuborilmaganini qidirayotgan operator `logs/async/` ga qarashi kerak, `logs/http/` ga emas.
  - Ma'lum cheklov: ikki jarayon bir xil `rotatelogs` nishoniga yozishi mumkin — ikkalasi ham `logs/async/.log` symlink'ini qayta yaratadi va 30 kunlik `WithMaxAge` tozalashini yurgizadi. Past xavf deb qabul qilingan (`O_APPEND`, symlink/tozalash idempotent), lekin log qatorlari aralashib ketsa sababi shu.

---

## Fayl tuzilishi

**Yaratiladi:**

| Fayl | Mas'uliyati |
|---|---|
| `test/architecture/error_response_test.go` | Qorovul: `JsonResponse(status, err)` shaklini butun `src/` da taqiqlaydi + qorovulning o'zini tekshiruvchi self-test. |
| `test/core/application/tasks/roi_task_sender_log_test.go` | Uch ROI sender'i nashr xatosini yutishini va uni zap orqali loglashini isbotlaydi. |

**O'zgartiriladi:**

| Fayl | O'zgarish |
|---|---|
| `src/entrypoint/presentation/handlers/social/usersocial_delete_handler.go` | `err` → `err.Error()` |
| `src/entrypoint/presentation/handlers/social/usersocial_update_handler.go` | `err` → `err.Error()` |
| `src/entrypoint/presentation/handlers/publisher/publisher_detail_handler.go` | `err` → `err.Error()` |
| `src/core/application/usecase/etaqrizusecases/find_reviewer_usecase.go` | Debug `fmt.Println` o'chiriladi + `fmt` importi |
| `src/core/application/usecase/authorshipclaimusecases/update_authorship_claim_status_usecase.go` | `fmt.Println` → zap; konstruktorga logger |
| `src/core/application/tasks/update_roi_task_sender.go` | `log.Println` → zap; konstruktorga logger |
| `src/core/application/tasks/set_roi_task_sender.go` | `log.Println` → zap; debug `fmt.Println` o'chiriladi; konstruktorga logger |
| `src/core/application/tasks/publish_roi_task_sender.go` | `log.Println` → zap; konstruktorga logger |
| `src/core/application/usecase/articleusecases/article_update_usecase.go` | O'lik `if err != nil` shoxi olib tashlanadi |
| `cmd/container/container.go`, `di/wire/provider.go` | `make wire-build` natijasi (qo'lda emas) |

**Import tozalash — har bir faylda aniq:**

- `find_reviewer_usecase.go`: `fmt` faqat shu bitta `Println` da ishlatiladi → **import olib tashlanadi**.
- `update_authorship_claim_status_usecase.go`: `fmt` **qoladi** — `sendToROI` da `fmt.Sprintf` bor.
- `article_update_usecase.go`: `fmt` **qoladi** — `checkId` da `fmt.Sprintf` bor.
- `set_roi_task_sender.go`: `fmt` va `log` **ikkalasi ham olib tashlanadi**.
- `update_roi_task_sender.go`, `publish_roi_task_sender.go`: `log` olib tashlanadi.

Bularning barchasini `go build ./...` baribir ushlaydi.

---

## Boshlashdan oldin: bazaviy holatni tasdiqlash

```bash
git branch --show-current            # feature/error-responses kutiladi
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | grep -c '^FAIL'    # 0 kutiladi
```

`FAIL` bo'lsa — **to'xtang va xabar bering**, bu rejaning aybi emas.

---

### Task 1: Qorovul + bo'sh `{"data":{}}` javoblari

Bu taskning testi qorovulning o'zi: avval u 3 ta mavjud buzilishni topib
yiqiladi, keyin handler'lar tuzatilgach yashil bo'ladi.

**Files:**
- Create: `test/architecture/error_response_test.go`
- Modify: `src/entrypoint/presentation/handlers/social/usersocial_delete_handler.go`
- Modify: `src/entrypoint/presentation/handlers/social/usersocial_update_handler.go`
- Modify: `src/entrypoint/presentation/handlers/publisher/publisher_detail_handler.go`

**Interfaces:**
- Consumes: `go/ast`, `go/parser`, `go/token` (standart kutubxona)
- Reuses: `sourceRoot(t)` va `relativeTo(root, path)` — `test/architecture/raw_sql_test.go` da, **bir xil `architecture_test` paketida**. Qayta e'lon qilinmaydi, aks holda kompilyatsiya xatosi.

**Step 1: Qorovul testini yozing (yiqilishi kerak)**

`test/architecture/error_response_test.go`:

```go
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

// responseMethods — javob tanasini JSON qilib serializatsiya qiladigan
// metodlar. JsonResponse — loyihaning o'z Context wrapper'i
// (src/entrypoint/presentation/app/context/context_wrap.go), JSON — echo'niki.
var responseMethods = map[string]bool{
	"JsonResponse": true,
	"JSON":         true,
}

// TestNoBareErrorInJSONResponse xato obyektini javob tanasi sifatida
// berishni taqiqlaydi.
//
// Nima uchun bu buzilish: JsonResponse(status, data) data ni
// encoding/json ga uzatadi. Go'ning `error` interfeysi — bu metod,
// maydon emas. errors.errorString kabi standart turlarda eksport
// qilinadigan maydon umuman yo'q, shuning uchun json.Marshal ularni
// {} deb yozadi. Natijada mijoz {"data":{},"status":400} oladi —
// status kod to'g'ri, lekin sabab yo'qoladi. To'g'ri yozuv —
// err.Error(), u {"data":"social does not belong to user"} beradi.
//
// CHEGARA (ataylab qoldirilgan): qorovul argument NOMI bo'yicha
// ishlaydi, TURI bo'yicha emas — `err` deb nomlangan yalang'och
// identifikatorni qidiradi. Ya'ni:
//
//  1. YOLG'ON SALBIY: `e` yoki `parseErr` deb nomlangan xato o'tib
//     ketadi. Loyihada xato o'zgaruvchisi deyarli har doim `err` deb
//     nomlanadi, shuning uchun amaliy qamrov yetarli. To'liq turga
//     asoslangan tekshiruv go/types va butun paketni yuklashni talab
//     qiladi — narxi foydasidan katta.
//  2. YOLG'ON IJOBIY: agar kimdir `err` deb nomlangan, lekin xato
//     bo'lmagan o'zgaruvchini javob sifatida bersa, qorovul yiqiladi.
//     Bunday holatda to'g'ri yechim — o'zgaruvchini qayta nomlash,
//     qorovulni yumshatish EMAS.
func TestNoBareErrorInJSONResponse(t *testing.T) {
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
			if !ok {
				return true
			}
			if !isBareErrResponse(call) {
				return true
			}
			selector := call.Fun.(*ast.SelectorExpr)
			pos := fset.Position(call.Lparen)
			violations = append(violations, fmt.Sprintf(
				"%s:%d  .%s(..., err) — xato obyekti javob tanasi sifatida berilgan",
				relativeTo(root, pos.Filename), pos.Line, selector.Sel.Name))
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("src/ daraxtini o'qib bo'lmadi: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf(
			"Xato obyekti javob tanasiga berilgan (%d ta joy). json.Marshal uni\n"+
				"{} deb yozadi va mijoz sababni ko'rmaydi. err o'rniga err.Error()\n"+
				"bering, yoki xatoni *response.Response bilan qaytaring — u holda\n"+
				"ResponseMiddleware statusni o'zi to'g'ri qo'yadi.\n\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

// isBareErrResponse chaqiruv .JsonResponse(status, err) yoki
// .JSON(status, err) shaklida ekanligini aniqlaydi.
func isBareErrResponse(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || !responseMethods[selector.Sel.Name] {
		return false
	}
	if len(call.Args) < 2 {
		return false
	}
	ident, ok := call.Args[1].(*ast.Ident)
	return ok && ident.Name == "err"
}
```

**Verify (yiqilishi kerak):**

```bash
go test ./test/architecture/ -run TestNoBareErrorInJSONResponse -count=1
```

Aynan **3 ta** buzilish ro'yxatlanishi kerak: `usersocial_delete_handler.go:40`,
`usersocial_update_handler.go:53`, `publisher_detail_handler.go:35`. Boshqa
son chiqsa — to'xtang va xabar bering.

**Step 2: Qorovulning tishlari borligini isbotlovchi self-test**

Bu test `src/` holatidan mustaqil: qorovul mantiqini ichki manba matnida
tekshiradi, shuning uchun kelajakda `src/` toza bo'lganda ham qorovul
ishlayotganini isbotlaydi. `test/architecture/error_response_test.go` oxiriga
qo'shing:

```go
// TestBareErrResponseDetectorHasTeeth qorovul mantiqining o'zini tekshiradi.
// src/ toza bo'lgach TestNoBareErrorInJSONResponse har doim yashil bo'ladi
// va u endi hech narsani isbotlamaydi — bu test esa detektorning haqiqatan
// ishlashini src/ holatidan qat'i nazar kafolatlaydi.
func TestBareErrResponseDetectorHasTeeth(t *testing.T) {
	const src = `package p

func handle(c ctx, err error, msg string) error {
	if bad {
		return c.JsonResponse(400, err)      // want: buzilish
	}
	if alsoBad {
		return c.JSON(500, err)              // want: buzilish
	}
	c.JsonResponse(400, err.Error())         // xavfsiz
	c.JsonResponse(400, msg)                 // xavfsiz
	c.JsonResponse(204, nil)                 // xavfsiz
	c.Logger().Error(err)                    // xavfsiz: JSON javob emas
	return nil
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "detector_teeth.go", src, 0)
	if err != nil {
		t.Fatalf("namuna manbani tahlil qilib bo'lmadi: %v", err)
	}

	var flagged []int
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if ok && isBareErrResponse(call) {
			flagged = append(flagged, fset.Position(call.Lparen).Line)
		}
		return true
	})

	want := []int{5, 8}
	if len(flagged) != len(want) {
		t.Fatalf("qorovul %d ta buzilish topdi, %d ta kutilgan edi: %v",
			len(flagged), len(want), flagged)
	}
	for i, line := range want {
		if flagged[i] != line {
			t.Errorf("%d-buzilish %d-qatorda topildi, %d kutilgan edi",
				i+1, flagged[i], line)
		}
	}
}
```

**Verify:**

```bash
go test ./test/architecture/ -run TestBareErrResponseDetectorHasTeeth -count=1
```

Bu **o'tishi** kerak (qorovul mantiqi to'g'ri), `TestNoBareErrorInJSONResponse`
esa hali ham yiqilib turadi (`src/` hali tuzatilmagan).

**Step 3: Uchta handler'ni tuzating**

`src/entrypoint/presentation/handlers/social/usersocial_delete_handler.go:40`:

```go
	err = this.usecase.Execute(uint(id), userID)
	if err != nil {
		return c.JsonResponse(http.StatusBadRequest, err.Error())
	}
```

`src/entrypoint/presentation/handlers/social/usersocial_update_handler.go:53`:

```go
	err = this.usecase.Execute(uint(id), dto, userID)
	if err != nil {
		return ctx.JsonResponse(http.StatusBadRequest, err.Error())
	}
```

`src/entrypoint/presentation/handlers/publisher/publisher_detail_handler.go:35`:

```go
	publisher, err := this.uc.Execute(uint(uintId))
	if err != nil {
		return c.JsonResponse(400, err.Error())
	}
```

**Verify:**

```bash
go build ./... && go vet ./...
go test ./test/architecture/ -count=1
go test ./... -count=1 2>&1 | grep -c '^FAIL'    # 0
```

**Commit:** `fix(handlers): xato obyekti o'rniga xabar matni qaytariladi + AST qorovul`

---

### Task 2: eTaqriz debug logi — shaxsiy ma'lumot sizishi (CWE-200)

`find_reviewer_usecase.go:34` butun `ReviewerEntity` ni stdout'ga bosadi.
Entity ichida `PhoneNumber` bor va unga oldingi CWE-200 ishida `json:"-"`
tegi qo'yilgan — ya'ni maydon HTTP javobidan ataylab yashirilgan, lekin
`fmt.Println` `%v` orqali uni baribir logga chiqaradi. Bu — o'sha ishning
qamrovidan chetda qolgan sizish nuqtasi.

**Files:**
- Modify: `src/core/application/usecase/etaqrizusecases/find_reviewer_usecase.go`

**Step 1: Sizish borligini tasdiqlang**

```bash
grep -n "json:\"-\"" src/core/domain/entity/reviewer*.go
grep -n "fmt.Println" src/core/application/usecase/etaqrizusecases/find_reviewer_usecase.go
```

`PhoneNumber` da `json:"-"` borligini o'z ko'zingiz bilan ko'ring. Bo'lmasa —
to'xtang va xabar bering: reja noto'g'ri asosga qurilgan bo'ladi.

**Step 2: Qatorni va endi ishlatilmaydigan importni o'chiring**

`find_reviewer_usecase.go`, 34-qator butunlay o'chiriladi:

```go
	reviewer, err = this.gateway.FindReviewerByScienceID(scienceID)
	if err != nil {
		return nil, err
	}
```

Import bloki `fmt` siz qoladi (`fmt` shu faylda boshqa hech qayerda
ishlatilmaydi):

```go
import (
	"errors"

	core "slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)
```

Debug logi o'rniga hech narsa qo'yilmaydi — bu diagnostik chiqish edi,
funksional talab emas. Xato holati allaqachon `return nil, err` orqali
chaqiruvchiga yetkaziladi.

**Verify:**

```bash
go build ./... && go vet ./...
grep -rn "fmt" src/core/application/usecase/etaqrizusecases/find_reviewer_usecase.go   # bo'sh chiqishi kerak
go test ./... -count=1 2>&1 | grep -c '^FAIL'    # 0
```

**Commit:** `fix(cwe-200): eTaqriz debug logi retsenzent ma'lumotlarini stdout'ga chiqarmaydi`

---

### Task 3: Mualliflik da'vosi — `fmt.Println` o'rniga zap

`update_authorship_claim_status_usecase.go:111` ROI xatosini `fmt.Println`
bilan stdout'ga yozadi: kontekst yo'q (qaysi maqola?), daraja yo'q, tuzilgan
maydonlar yo'q, va production log fayliga tushmaydi.

**Files:**
- Modify: `src/core/application/usecase/authorshipclaimusecases/update_authorship_claim_status_usecase.go`
- Regenerate: `cmd/container/container.go`, `di/wire/provider.go` (`make wire-build`)

**Interfaces:**
- Consumes: `slib.uz/src/infrastructure/logger` (`*logger.AsyncLogger` — embedded `*zap.Logger`), `go.uber.org/zap`
- `core/` dan `infrastructure/logger` ga import — bu shoxda YANGI naqsh (tekshirilgan: bu o'zgarishgacha `src/core/` da uni import qiladigan 0 ta fayl bor edi). Eng yaqin mavjud pretsedent — `core/` dan `infrastructure/config` ga import (7 faylda). Port abstraksiyasidan ataylab voz kechilgan: port `zap.Field` ni oshkor qilishi kerak bo'lardi, bu esa `infrastructure` ni `core` ga sizdirardi — mavjud bo'lmagan pretsedent asosida emas.
- `NewAsyncLogger()` wire'da allaqachon provayder sifatida ro'yxatdan o'tgan (`cmd/container/container.go:832`), shuning uchun `make wire-build` uni o'zi ulaydi.

**Step 1: Struktura va konstruktorga logger qo'shing**

```go
type UpdateAuthorshipClaimStatusUseCase struct {
	claimRepo                  repository.AuthorshipClaimRepository
	articleRepo                repository.ArticleRepository
	publishedArticleRepository repository.PublishedArticleRepository
	authorRepo                 repository.AuthorRepository
	roiGateway                 gateway.ROIGateway
	frontendURL                string
	log                        *logger.AsyncLogger
}

// @inject
func NewUpdateAuthorshipClaimStatusUseCase(
	claimRepo repository.AuthorshipClaimRepository,
	articleRepo repository.ArticleRepository,
	publishedArticleRepository repository.PublishedArticleRepository,
	authorRepo repository.AuthorRepository,
	roiGateway gateway.ROIGateway,
	configAdapter conf.ConfigAdapter,
	log *logger.AsyncLogger,
) *UpdateAuthorshipClaimStatusUseCase {
	return &UpdateAuthorshipClaimStatusUseCase{
		claimRepo:                  claimRepo,
		articleRepo:                articleRepo,
		publishedArticleRepository: publishedArticleRepository,
		authorRepo:                 authorRepo,
		roiGateway:                 roiGateway,
		frontendURL:                configAdapter.GetFrontendURL(),
		log:                        log,
	}
}
```

Import blokiga qo'shing (`fmt` **qoladi** — `sendToROI` da `fmt.Sprintf` bor):

```go
	"go.uber.org/zap"
	"slib.uz/src/infrastructure/logger"
```

**Step 2: `fmt.Println` ni almashtiring**

`Execute` ichida, 110–112-qatorlar:

```go
		if err := this.sendToROI(article); err != nil {
			// Xato ataylab yutiladi: da'vo allaqachon tasdiqlangan va
			// bazaga yozilgan, ROI'ga yuborish esa yordamchi qadam.
			this.log.Error("ROI'ga maqola yuborilmadi",
				zap.Uint("article_id", article.ID),
				zap.Uint("claim_id", input.ClaimID),
				zap.Error(err))
		}
```

**Step 3: Wire'ni qayta generatsiya qiling**

```bash
make wire-build
git diff --stat cmd/container/container.go di/wire/provider.go
```

`container.go` da `NewUpdateAuthorshipClaimStatusUseCase(...)` chaqiruviga
`asyncLogger` argumenti qo'shilgan bo'lishi kerak. Fayl qo'lda tahrir
QILINMAYDI — agar `make wire-build` yiqilsa, to'xtang va xabar bering.

**Verify:**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | grep -c '^FAIL'    # 0
```

**Commit:** `refactor(authorship-claim): ROI xatosi zap orqali tuzilgan holda loglanadi`

---

### Task 4: ROI task sender'lari — zap va o'lik shoxni olib tashlash

Uchala ROI sender'i ham `Publish` xatosini `log.Println` bilan yutadi va
**doim `nil` qaytaradi**. Shu sababli `article_update_usecase.go:103–106`
dagi `if err != nil` shoxi hech qachon bajarilmaydi — o'lik kod. Task ikkala
narsani bir vaqtda tuzatadi: loglashni to'g'rilaydi va endi yolg'on
tasalli beradigan tekshiruvni olib tashlaydi.

**Files:**
- Create: `test/core/application/tasks/roi_task_sender_log_test.go`
- Modify: `src/core/application/tasks/update_roi_task_sender.go`
- Modify: `src/core/application/tasks/set_roi_task_sender.go`
- Modify: `src/core/application/tasks/publish_roi_task_sender.go`
- Modify: `src/core/application/usecase/articleusecases/article_update_usecase.go`
- Regenerate: `cmd/container/container.go`, `di/wire/provider.go`

**Interfaces:**
- Consumes: `publisher.TaskPublisher` (mavjud port), `*logger.AsyncLogger`, `go.uber.org/zap`
- Test uchun: `go.uber.org/zap/zaptest/observer` — yozilgan yozuvlarni xotirada ushlab qoladi.
- `*logger.AsyncLogger` — `*zap.Logger` ni embed qiladigan struktura, shuning uchun testda `&logger.AsyncLogger{Logger: zap.New(core)}` deb qurish mumkin.

**Step 1: Testni yozing (yiqilishi kerak — konstruktorlar hali logger olmaydi)**

`test/core/application/tasks/roi_task_sender_log_test.go`:

```go
package tasks_test

import (
	"errors"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/logger"
)

// failingPublisher — Publish har doim xato qaytaradigan soxta nashr etuvchi.
type failingPublisher struct {
	publisherPort
	err error
}

func (this *failingPublisher) Publish(task *entity.TaskEntity[any], maxRetryCount int) error {
	return this.err
}

func observedLogger() (*logger.AsyncLogger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.ErrorLevel)
	return &logger.AsyncLogger{Logger: zap.New(core)}, logs
}

// TestRoiSendersSwallowPublishErrorAndLogIt — uchala ROI sender'i ham
// navbatga qo'yish xatosini chaqiruvchiga qaytarmaydi (ROI yuborish
// yordamchi qadam), lekin uni jimgina yo'qotmaydi ham: xato zap orqali
// ERROR darajasida yoziladi.
func TestRoiSendersSwallowPublishErrorAndLogIt(t *testing.T) {
	publishErr := errors.New("navbat mavjud emas")

	tests := []struct {
		name string
		run  func(pub *failingPublisher, log *logger.AsyncLogger) error
	}{
		{
			name: "UpdateRoiSenderTask",
			run: func(pub *failingPublisher, log *logger.AsyncLogger) error {
				return tasks.NewUpdateRoiSenderTask(pub, log).
					Run(tasks.UpdateRoiSenderPayload{ArticleID: 42})
			},
		},
		{
			name: "PublishRoiSenderTask",
			run: func(pub *failingPublisher, log *logger.AsyncLogger) error {
				return tasks.NewPublishRoiSenderTask(pub, log).
					Run(tasks.PublishRoiSenderPayload{ArticleID: 42})
			},
		},
		{
			name: "SetRoiSenderTask",
			run: func(pub *failingPublisher, log *logger.AsyncLogger) error {
				return tasks.NewSetRoiSenderTask(pub, log).
					Run(tasks.SetRoiSenderPayload{ApplicationID: 42})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, logs := observedLogger()
			pub := &failingPublisher{err: publishErr}

			if err := tt.run(pub, log); err != nil {
				t.Fatalf("%s.Run xato qaytardi: %v (nil kutilgan edi — "+
					"ROI navbati asosiy oqimni yiqitmasligi kerak)", tt.name, err)
			}

			entries := logs.All()
			if len(entries) != 1 {
				t.Fatalf("%s: %d ta ERROR yozuvi yozildi, 1 ta kutilgan edi",
					tt.name, len(entries))
			}
			if !containsErrorField(entries[0], publishErr) {
				t.Errorf("%s: log yozuvida asl xato yo'q: %+v",
					tt.name, entries[0].ContextMap())
			}
		})
	}
}

func containsErrorField(entry observer.LoggedEntry, want error) bool {
	for _, field := range entry.Context {
		if field.Key == "error" && field.Interface == want {
			return true
		}
	}
	return false
}
```

`publisherPort` embed'i uchun fayl boshiga qo'shing (interfeysni embed qilish
— loyihaning soxta repozitoriy odati: faqat kerakli metod amalga oshiriladi):

```go
import portpublisher "slib.uz/src/core/domain/ports/tasks/publisher"

type publisherPort = portpublisher.TaskPublisher
```

**Diqqat:** `TaskPublisher` interfeysining aniq metodlarini avval tekshiring:

```bash
cat src/core/domain/ports/tasks/publisher/*.go
```

`Publish(task *entity.TaskEntity[any], maxRetryCount int) error` imzosi
mos kelmasa — testni haqiqiy imzoga moslang, teskarisiga emas.

**Verify (yiqilishi kerak):** kompilyatsiya xatosi — konstruktorlar 1 ta
argument kutadi, test 2 ta beradi. Bu kutilgan holat.

```bash
go test ./test/core/application/tasks/ -count=1
```

**Step 2: Uchala sender'ni logger bilan qayta yozing**

`src/core/application/tasks/update_roi_task_sender.go` to'liq:

```go
package tasks

import (
	"go.uber.org/zap"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
	"slib.uz/src/infrastructure/logger"
)

type UpdateRoiSenderTask struct {
	publisher publisher.TaskPublisher
	log       *logger.AsyncLogger
}

// @inject
func NewUpdateRoiSenderTask(publisher publisher.TaskPublisher, log *logger.AsyncLogger) *UpdateRoiSenderTask {
	return &UpdateRoiSenderTask{publisher: publisher, log: log}
}

type UpdateRoiSenderPayload struct {
	ArticleID uint
}

// Run maqolani ROI yangilash navbatiga qo'yadi. Navbat nosozligi
// chaqiruvchiga qaytarilmaydi — ROI yangilash yordamchi qadam va u
// maqolani saqlashni yiqitmasligi kerak. Shuning uchun Run har doim
// nil qaytaradi; imzodagi error kelajakdagi kengaytirish uchun qolgan.
func (this *UpdateRoiSenderTask) Run(payload UpdateRoiSenderPayload) error {
	task := entity.NewTaskEntity[any](enum.TaskUpdateRoi, payload)

	if err := this.publisher.Publish(task, 5); err != nil {
		this.log.Error("ROI yangilash vazifasi navbatga qo'yilmadi",
			zap.Uint("article_id", payload.ArticleID),
			zap.Error(err))
	}

	return nil
}
```

`publish_roi_task_sender.go` — xuddi shu shakl, `enum.TaskPublishRoi`,
xabar `"ROI nashr vazifasi navbatga qo'yilmadi"`, maydon
`zap.Uint("article_id", payload.ArticleID)`.

`set_roi_task_sender.go` — xuddi shu shakl, `enum.TaskSetRoi`, xabar
`"ROI o'rnatish vazifasi navbatga qo'yilmadi"`, maydon
`zap.Uint("application_id", payload.ApplicationID)`. Bundan tashqari
**28-qatordagi `fmt.Println("SetRoiSenderTask.Run")` butunlay o'chiriladi**
va `fmt` importi ham olib tashlanadi. Bu — production'da har chaqiruvda
stdout'ga tushadigan diagnostik chiqish, uning o'rniga hech narsa
qo'yilmaydi.

**Step 3: Wire'ni qayta generatsiya qiling**

```bash
make wire-build
go build ./...
```

Uchala sender `cmd/container/container.go` da bir necha joyda quriladi
(`setRoiSenderTask` — 218 va 897-qatorlar atrofida, `updateRoiSenderTask` —
283, `publishRoiSenderTask` — 289). Hammasiga `asyncLogger` qo'shilishi
kerak; `make wire-build` buni o'zi bajaradi.

**Verify:**

```bash
go test ./test/core/application/tasks/ -count=1    # endi o'tishi kerak
```

**Step 4: `article_update_usecase.go` dagi o'lik shoxni olib tashlang**

`Run` hech qachon xato qaytarmagani uchun 103–106-qatorlardagi tekshiruv
o'lik kod. Uni ochiq-oydin e'tiborsizlik bilan almashtiring:

```go
	// UpdateRoiSenderTask.Run navbat xatosini o'zi loglaydi va har doim
	// nil qaytaradi — maqolani yangilash ROI navbatiga bog'liq emas.
	_ = this.sendToROI(publishedArticle)

	return nil
}
```

`fmt` importi **qoladi** — `checkId` da `fmt.Sprintf` ishlatiladi.

**Verify:**

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | grep -c '^FAIL'    # 0
grep -rn "fmt.Println\|log.Println" src/core/application/tasks/*roi*    # bo'sh
```

**Commit:** `refactor(tasks): ROI sender'lari zap orqali loglaydi, o'lik xato tekshiruvi olib tashlandi`

---

## Yakuniy tekshirish

```bash
go build ./... && go vet ./...
go test ./... -count=1 2>&1 | tail -5
go test ./... -count=1 -v 2>&1 | grep -c '^--- PASS'
```

Kutilgan natija:

- `go build`, `go vet` — toza.
- `FAIL` — 0 ta.
- O'tuvchi testlar soni **174 dan yuqori** (Task 1 dan 2 ta, Task 4 dan 1 ta yuqori darajali test qo'shiladi).
- `git status` — `cmd/container/container.go` va `di/wire/provider.go` o'zgargan (wire generatsiyasi), boshqa generatsiya qilinadigan fayllar tegilmagan.

Yakuniy nazorat grep'lari — hammasi bo'sh chiqishi kerak:

```bash
grep -rnE "(JsonResponse|JSON)\(.*,\s*err\s*\)" --include=*.go src/
grep -rn "fmt.Println" src/core/application/usecase/etaqrizusecases/
grep -rn "log.Println" src/core/application/tasks/
```

---

## Qamrovdan tashqarida (ataylab)

- **`send_article_kafka_task_sender.go:21`** dagi `fmt.Println("SendArticleKafkaTaskSender.Run")` — bir xil naqsh, lekin ROI oqimiga tegishli emas. Alohida ish sifatida qoldiriladi.
- **Xatoni yutish semantikasini o'zgartirish.** ROI navbatiga yuborish nosozligi so'rovni yiqitishi kerakmi degan savol — mahsulot qarori, texnik tozalash emas.
- **`audit-javob-report.md` dagi qolgan "Natija" da'volarini kod bilan solishtirish** — alohida ish sifatida kelishilgan.
- **`response_middleware.go` ni `echo.HTTPErrorHandler` ga ko'chirish** — dizayn hujjatida ko'rib chiqilgan va rad etilgan (`docs/superpowers/specs/2026-08-17-xato-javoblari-design.md`, 5-bo'lim).
