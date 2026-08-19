# Fayl yuklash va saqlangan XSS himoyasi — implementatsiya rejasi

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fayl yuklash zanjirini yopish — ruxsat etilmagan kengaytma uchun presigned URL berilmaydi, MinIO Content-Type va hajmni majburlaydi, brauzer yuklangan faylni ijro etmaydi.

**Architecture:** Uch qatlam. L1 — `src/core/domain/entity/enum` dagi allow-list, URL berish paytida ishlaydi. L2 — `PresignedPostPolicy`, MinIO yuklash paytida Content-Type, hajm va obyekt kalitini majburlaydi va `Content-Disposition` ni obyekt metadatasiga yozadi (ommaviy bucket anonim o'qilgani uchun bu yerda **yagona** himoya). L3 — presigned GET javob sarlavhalari, maxfiy bucket uchun.

**Tech Stack:** Go 1.25, Echo v4.13.3, minio-go v7.0.100, google/wire, swaggo.

**Spec:** `docs/superpowers/specs/2026-08-13-upload-hardening-design.md`

## Global Constraints

- **Yangi bog'liqlik yo'q.** `go.mod` va `go.sum` o'zgarmasligi shart. Testlar standart `testing` va qo'lda yozilgan fake'lar bilan.
- **Testlar joyi:** loyiha rootidagi `test/` katalogi, paket yo'lini takrorlab. Masalan `src/core/domain/entity/enum/upload_type.go` → `test/core/domain/entity/enum/upload_type_test.go`, paket nomi `enum_test`. Yopiq identifikatorlarga kirish yo'q — faqat eksport qilinganlar orqali.
- **Core qatlami** (`src/core/**`) faqat `github.com/labstack/gommon/log` bilan loglaydi, hech qachon `zap`. Core infratuzilmani to'g'ridan-to'g'ri import qilmaydi — faqat `src/core/domain/ports/**` orqali.
- **Xatolik holatida rad etish.** Tur aniqlanmasa, magic-byte mos kelmasa yoki shubha bo'lsa — yuklash rad etiladi, hech qachon ruxsat berilmaydi.
- **`cmd/container/container.go` (wire) va `src/entrypoint/presentation/docs/*` (swaggo) generatsiya qilinadi.** Qo'lda tahrirlash mumkin emas. Konstruktor imzosi o'zgarsa `make wire-build`, swagger annotatsiyasi o'zgarsa `make generate-docs` ishga tushiriladi; keyin `git diff` bo'sh bo'lishi tekshiriladi.
- **Yangi env o'zgaruvchisi kerak emas.** `UploadedFileMaxSize` (`env:"UPLOADED_FILE_MAX_SIZE,required"`) `src/infrastructure/config/env.go:51` da allaqachon mavjud.
- **IDE diagnostikasi bu repozitoriyda muntazam eskirgan.** Faqat `go build ./...`, `go vet ./...`, `go test ./... -count=1` natijalariga ishoning.
- **Bazaviy holat:** shox `feature/upload-hardening`, `go build` toza, `go test ./... -count=1` 74 test o'tadi.

---

## Fayl tuzilishi

| Fayl | Mas'uliyat | Task |
|---|---|---|
| `src/core/domain/entity/enum/upload_type.go` | **Yangi.** Allow-list: kengaytma → MIME → magic-byte → disposition. Yagona haqiqat manbai | 1 |
| `test/core/domain/entity/enum/upload_type_test.go` | **Yangi.** Allow-list testlari | 1 |
| `src/infrastructure/storage/post_policy.go` | **Yangi.** `BuildPostPolicy` — sof funksiya, MinIO serversiz testlanadi | 2 |
| `test/infrastructure/storage/post_policy_test.go` | **Yangi.** Policy shartlari testlari | 2 |
| `src/core/domain/ports/storage/file_storage.go` | `PutObjectPresignedUrl` → `PostPolicyPresignedUrl` | 2 |
| `src/infrastructure/storage/minio_storage.go` | Port implementatsiyasi, hajm chegarasi, presigned GET parametrlari | 2, 4, 5 |
| `src/core/application/usecase/uploadusecases/upload_file_usecase.go` | Allow-list, papka→bucket, kalit generatsiyasi | 3 |
| `test/core/application/usecase/uploadusecases/upload_file_usecase_test.go` | **Yangi.** Presigned yuklash testlari | 3 |
| `src/entrypoint/presentation/handlers/upload/schema/presigned_put_url_schema.go` | So'rov/javob shakli | 3 |
| `src/entrypoint/presentation/handlers/upload/presigned_put_url_handler.go` | Handler | 3 |
| `src/core/application/usecase/uploadusecases/upload_tempfile_usecase.go` | Magic-byte tekshiruvi | 4 |
| `test/core/application/usecase/uploadusecases/upload_tempfile_usecase_test.go` | **Yangi.** Magic-byte testlari | 4 |
| `src/core/application/usecase/uploadusecases/upload_base64_file_usecase.go` | Umumiy magic-byte funksiyasiga o'tish | 5 |
| `src/entrypoint/presentation/app/app.go` | Xavfsizlik sarlavhalari | 6 |
| `test/entrypoint/presentation/app/security_headers_test.go` | **Yangi.** Sarlavha testlari | 6 |

---

## Task 1: Allow-list — yagona haqiqat manbai

**Files:**
- Create: `src/core/domain/entity/enum/upload_type.go`
- Test: `test/core/domain/entity/enum/upload_type_test.go`

**Interfaces:**
- Consumes: hech narsa (birinchi task)
- Produces:
  - `type UploadType struct { Extension, ContentType, Disposition string; Magic [][]byte }`
  - `func LookupUploadType(ext string) (*UploadType, bool)`
  - `func (u *UploadType) MatchesMagic(head []byte) bool`
  - Konstantalar: `DispositionAttachment = "attachment"`, `DispositionInline = "inline"`
  - `const MagicPrefixLen = 8`

Bu task butun ishning poydevori — keyingi barcha tasklar shu funksiyalarni iste'mol qiladi. Boshqa hech qaysi faylda kengaytma yoki MIME turi ro'yxati yozilmasligi kerak.

- [ ] **Step 1: Failing testni yozish**

`test/core/domain/entity/enum/upload_type_test.go`:

```go
package enum_test

import (
	"testing"

	"slib.uz/src/core/domain/entity/enum"
)

func TestAllowedExtensionsAreFound(t *testing.T) {
	cases := []struct {
		ext         string
		contentType string
		disposition string
	}{
		{".pdf", "application/pdf", enum.DispositionAttachment},
		{".jpg", "image/jpeg", enum.DispositionInline},
		{".jpeg", "image/jpeg", enum.DispositionInline},
		{".png", "image/png", enum.DispositionInline},
		{".doc", "application/msword", enum.DispositionAttachment},
		{".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", enum.DispositionAttachment},
	}

	for _, c := range cases {
		ut, ok := enum.LookupUploadType(c.ext)
		if !ok {
			t.Fatalf("%s ruxsat etilgan bo'lishi kerak edi, topilmadi", c.ext)
		}
		if ut.ContentType != c.contentType {
			t.Errorf("%s: content-type %q kutilgandi, %q keldi", c.ext, c.contentType, ut.ContentType)
		}
		if ut.Disposition != c.disposition {
			t.Errorf("%s: disposition %q kutilgandi, %q keldi", c.ext, c.disposition, ut.Disposition)
		}
	}
}

func TestLookupIsCaseInsensitive(t *testing.T) {
	for _, ext := range []string{".PDF", ".pdf", ".pDf"} {
		ut, ok := enum.LookupUploadType(ext)
		if !ok {
			t.Fatalf("%s topilishi kerak edi", ext)
		}
		if ut.ContentType != "application/pdf" {
			t.Errorf("%s: application/pdf kutilgandi, %q keldi", ext, ut.ContentType)
		}
	}
}

func TestDangerousExtensionsAreRejected(t *testing.T) {
	// .svg ataylab ro'yxatda yo'q: SVG ichida <script> ishlaydi.
	for _, ext := range []string{".html", ".htm", ".svg", ".exe", ".php", ".js", ".sh", ""} {
		if _, ok := enum.LookupUploadType(ext); ok {
			t.Errorf("%q rad etilishi kerak edi, lekin qabul qilindi", ext)
		}
	}
}

func TestMagicBytesMatch(t *testing.T) {
	pdf, _ := enum.LookupUploadType(".pdf")
	if !pdf.MatchesMagic([]byte("%PDF-1.7\n")) {
		t.Error("haqiqiy PDF boshi qabul qilinishi kerak edi")
	}
	if pdf.MatchesMagic([]byte("<html><scr")) {
		t.Error("HTML baytlari PDF sifatida QABUL QILINMASLIGI kerak edi")
	}

	png, _ := enum.LookupUploadType(".png")
	if !png.MatchesMagic([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		t.Error("haqiqiy PNG boshi qabul qilinishi kerak edi")
	}

	jpg, _ := enum.LookupUploadType(".jpg")
	if !jpg.MatchesMagic([]byte{0xFF, 0xD8, 0xFF, 0xE0}) {
		t.Error("haqiqiy JPEG boshi qabul qilinishi kerak edi")
	}
}

func TestShortInputIsRejectedNotPanicked(t *testing.T) {
	png, _ := enum.LookupUploadType(".png")
	if png.MatchesMagic([]byte{0x89, 0x50}) {
		t.Error("imzodan qisqa kirish rad etilishi kerak edi")
	}
	if png.MatchesMagic(nil) {
		t.Error("bo'sh kirish rad etilishi kerak edi")
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

Run: `go test ./test/core/domain/entity/enum/... -count=1`
Expected: FAIL — `undefined: enum.LookupUploadType` (kompilyatsiya xatosi)

- [ ] **Step 3: Minimal implementatsiyani yozish**

`src/core/domain/entity/enum/upload_type.go`:

```go
package enum

import (
	"bytes"
	"strings"
)

// Content-Disposition qiymatlari.
const (
	DispositionAttachment = "attachment"
	DispositionInline     = "inline"
)

// MagicPrefixLen — imzoni tekshirish uchun yetarli bosh baytlar soni.
// Ro'yxatdagi eng uzun imzo 8 bayt (PNG).
const MagicPrefixLen = 8

// UploadType ruxsat etilgan bitta fayl turini tavsiflaydi.
// Kengaytma, MIME turi, imzo va disposition bir joyda turadi — ular
// uch xil faylda yashasa, vaqt o'tishi bilan uzoqlashib ketadi.
type UploadType struct {
	Extension   string
	ContentType string
	Disposition string
	// Magic — ruxsat etilgan bosh baytlar variantlari.
	// Bo'sh bo'lsa imzo tekshirilmaydi (hozircha bunday tur yo'q).
	Magic [][]byte
}

var uploadTypes = map[string]*UploadType{
	".pdf": {
		Extension:   ".pdf",
		ContentType: "application/pdf",
		Disposition: DispositionAttachment,
		Magic:       [][]byte{[]byte("%PDF-")},
	},
	".jpg": {
		Extension:   ".jpg",
		ContentType: "image/jpeg",
		Disposition: DispositionInline,
		Magic:       [][]byte{{0xFF, 0xD8, 0xFF}},
	},
	".jpeg": {
		Extension:   ".jpeg",
		ContentType: "image/jpeg",
		Disposition: DispositionInline,
		Magic:       [][]byte{{0xFF, 0xD8, 0xFF}},
	},
	".png": {
		Extension:   ".png",
		ContentType: "image/png",
		Disposition: DispositionInline,
		Magic:       [][]byte{{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
	},
	".doc": {
		Extension:   ".doc",
		ContentType: "application/msword",
		Disposition: DispositionAttachment,
		Magic:       [][]byte{{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	},
	".docx": {
		Extension:   ".docx",
		ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		Disposition: DispositionAttachment,
		// PK\x03\x04 — ZIP konteyneri. Bu imzo .docx ni boshqa ZIP asosidagi
		// formatlardan ajratmaydi; konteyner ichini ochish spec doirasidan tashqarida.
		Magic: [][]byte{{0x50, 0x4B, 0x03, 0x04}},
	},
}

// LookupUploadType kengaytma bo'yicha ruxsat etilgan turni qaytaradi.
// Katta-kichik harf ahamiyatsiz. Topilmasa (nil, false) — chaqiruvchi
// yuklashni RAD ETISHI shart.
func LookupUploadType(ext string) (*UploadType, bool) {
	ut, ok := uploadTypes[strings.ToLower(ext)]
	return ut, ok
}

// MatchesMagic head baytlari shu turning imzolaridan biriga mos kelishini tekshiradi.
// head imzodan qisqa bo'lsa false qaytaradi — panic emas.
func (u *UploadType) MatchesMagic(head []byte) bool {
	if len(u.Magic) == 0 {
		return true
	}
	for _, sig := range u.Magic {
		if bytes.HasPrefix(head, sig) {
			return true
		}
	}
	return false
}
```

- [ ] **Step 4: Testni ishga tushirib, o'tishini tasdiqlash**

Run: `go test ./test/core/domain/entity/enum/... -count=1 -v`
Expected: PASS — 5 ta test

- [ ] **Step 5: To'liq to'plamni tekshirish**

Run: `go build ./... && go vet ./... && go test ./... -count=1`
Expected: build va vet toza, 74 + 5 = 79 test o'tadi

- [ ] **Step 6: Commit**

```bash
git add src/core/domain/entity/enum/upload_type.go test/core/domain/entity/enum/upload_type_test.go
git commit -m "feat(upload): ruxsat etilgan fayl turlari uchun yagona allow-list"
```

---

## Task 2: POST policy quruvchi va storage porti

**Files:**
- Create: `src/infrastructure/storage/post_policy.go`
- Test: `test/infrastructure/storage/post_policy_test.go`
- Modify: `src/core/domain/ports/storage/file_storage.go:21`
- Modify: `src/infrastructure/storage/minio_storage.go:54-56`

**Interfaces:**
- Consumes: Task 1 dan `enum.UploadType`
- Produces:
  - `func BuildPostPolicy(bucketName, objectName string, ut *enum.UploadType, maxSize int64, expires time.Time) (*minio.PostPolicy, error)`
  - Port metodi: `PostPolicyPresignedUrl(ctx context.Context, bucket enum.Bucket, objectName string, ut *enum.UploadType, maxSize int64) (*url.URL, map[string]string, error)`

`BuildPostPolicy` ataylab sof funksiya sifatida ajratiladi: `MinioStorage` ning qolgan qismi haqiqiy MinIO serverini talab qiladi, policy mantig'i esa `PostPolicy.String()` orqali serversiz testlanadi.

- [ ] **Step 1: Failing testni yozish**

`test/infrastructure/storage/post_policy_test.go`:

```go
package storage_test

import (
	"strings"
	"testing"
	"time"

	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/storage"
)

func mustPolicyJSON(t *testing.T, ext string, maxSize int64) string {
	t.Helper()
	ut, ok := enum.LookupUploadType(ext)
	if !ok {
		t.Fatalf("test sozlash xatosi: %s allow-list'da yo'q", ext)
	}
	// Sobit vaqt: test natijasi soatga bog'liq bo'lmasligi kerak.
	expires := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := storage.BuildPostPolicy("slibbucketprivate", "articles/a_1234.pdf", ut, maxSize, expires)
	if err != nil {
		t.Fatalf("BuildPostPolicy xato qaytardi: %v", err)
	}
	return p.String()
}

func TestPolicyPinsContentType(t *testing.T) {
	got := mustPolicyJSON(t, ".pdf", 10<<20)
	if !strings.Contains(got, "application/pdf") {
		t.Errorf("policy Content-Type ni qat'iy belgilashi kerak edi.\npolicy: %s", got)
	}
	if !strings.Contains(got, "Content-Type") {
		t.Errorf("policy shartlarida Content-Type bo'lishi kerak edi.\npolicy: %s", got)
	}
}

func TestPolicyPinsContentLengthRange(t *testing.T) {
	got := mustPolicyJSON(t, ".pdf", 10<<20)
	if !strings.Contains(got, "content-length-range") {
		t.Errorf("policy hajm oralig'ini belgilashi kerak edi.\npolicy: %s", got)
	}
	if !strings.Contains(got, "10485760") {
		t.Errorf("policy'da maksimal hajm (10485760) bo'lishi kerak edi.\npolicy: %s", got)
	}
}

func TestPolicyPinsExactKey(t *testing.T) {
	got := mustPolicyJSON(t, ".pdf", 10<<20)
	if !strings.Contains(got, "articles/a_1234.pdf") {
		t.Errorf("policy obyekt kalitini qat'iy belgilashi kerak edi.\npolicy: %s", got)
	}
}

func TestPolicySetsDispositionFromUploadType(t *testing.T) {
	doc := mustPolicyJSON(t, ".pdf", 10<<20)
	if !strings.Contains(doc, enum.DispositionAttachment) {
		t.Errorf("hujjat uchun attachment kutilgandi.\npolicy: %s", doc)
	}

	ut, _ := enum.LookupUploadType(".png")
	expires := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := storage.BuildPostPolicy("slibbucketpublic", "news/x_1234.png", ut, 10<<20, expires)
	if err != nil {
		t.Fatalf("BuildPostPolicy xato qaytardi: %v", err)
	}
	if !strings.Contains(p.String(), enum.DispositionInline) {
		t.Errorf("rasm uchun inline kutilgandi.\npolicy: %s", p.String())
	}
}

func TestNilUploadTypeIsRejected(t *testing.T) {
	expires := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err := storage.BuildPostPolicy("b", "k.pdf", nil, 10<<20, expires); err == nil {
		t.Error("nil UploadType uchun xato kutilgandi")
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

Run: `go test ./test/infrastructure/storage/... -count=1`
Expected: FAIL — `undefined: storage.BuildPostPolicy`

- [ ] **Step 3: Policy quruvchini yozish**

`src/infrastructure/storage/post_policy.go`:

```go
package storage

import (
	"errors"
	"time"

	"github.com/minio/minio-go/v7"
	"slib.uz/src/core/domain/entity/enum"
)

// BuildPostPolicy yuklash uchun imzolanadigan POST policy quradi.
//
// Policy'ga kiritilgan har bir shart MinIO tomonidan majburlanadi: mijoz
// boshqa Content-Type, boshqa hajm yoki boshqa obyekt kaliti bilan yuklay
// olmaydi — yuklash rad etiladi.
//
// Content-Disposition obyekt metadatasiga yoziladi. Bu ommaviy bucket uchun
// hal qiluvchi: u anonim o'qishga ochiq, ya'ni presigned GET paytidagi
// sarlavha override'i u yerda umuman qo'llanmaydi.
func BuildPostPolicy(
	bucketName string,
	objectName string,
	ut *enum.UploadType,
	maxSize int64,
	expires time.Time,
) (*minio.PostPolicy, error) {
	if ut == nil {
		return nil, errors.New("upload type aniqlanmagan")
	}
	if maxSize <= 0 {
		return nil, errors.New("maksimal hajm musbat bo'lishi kerak")
	}

	p := minio.NewPostPolicy()

	if err := p.SetBucket(bucketName); err != nil {
		return nil, err
	}
	if err := p.SetKey(objectName); err != nil {
		return nil, err
	}
	if err := p.SetContentType(ut.ContentType); err != nil {
		return nil, err
	}
	if err := p.SetContentLengthRange(1, maxSize); err != nil {
		return nil, err
	}
	if err := p.SetContentDisposition(ut.Disposition); err != nil {
		return nil, err
	}
	if err := p.SetExpires(expires); err != nil {
		return nil, err
	}

	return p, nil
}
```

- [ ] **Step 4: Testni ishga tushirib, o'tishini tasdiqlash**

Run: `go test ./test/infrastructure/storage/... -count=1 -v`
Expected: PASS — 5 ta test

- [ ] **Step 5: Portni o'zgartirish**

`src/core/domain/ports/storage/file_storage.go` — 21-qatordagi `PutObjectPresignedUrl` o'rniga:

```go
	PostPolicyPresignedUrl(ctx context.Context, bucket enum.Bucket, objectName string, ut *enum.UploadType, maxSize int64) (*url.URL, map[string]string, error)
```

- [ ] **Step 6: MinIO implementatsiyasini yozish**

`src/infrastructure/storage/minio_storage.go` — 54-56 qatorlardagi `PutObjectPresignedUrl` metodini quyidagiga almashtiring:

```go
func (this *MinioStorage) PostPolicyPresignedUrl(
	ctx context.Context,
	bucket enum.Bucket,
	objectName string,
	ut *enum.UploadType,
	maxSize int64,
) (*url.URL, map[string]string, error) {
	policy, err := BuildPostPolicy(
		this.getBucketName(bucket),
		objectName,
		ut,
		maxSize,
		time.Now().Add(30*time.Minute),
	)
	if err != nil {
		return nil, nil, err
	}
	return this.minio.PresignedPostPolicy(ctx, policy)
}
```

- [ ] **Step 7: Build yiqilishini kuzatish va chaqiruv joyini vaqtincha moslashtirish**

Run: `go build ./...`
Expected: FAIL — `upload_file_usecase.go:25` eski metodni chaqiradi. Bu kutilgan: Task 3 uni to'liq qayta yozadi.

Shoxni yashil holatda saqlash uchun `src/core/application/usecase/uploadusecases/upload_file_usecase.go:25` dagi chaqiruvni vaqtincha moslashtiring — `Execute` ning tanasi Task 3 da butunlay almashtiriladi:

```go
	presignedUrl, _, err := this.storage.PostPolicyPresignedUrl(ctx, bucket, objectName, nil, 1)
	if err != nil {
		return "", "", err
	}
	return presignedUrl.String(), objectName, nil
```

Diqqat: bu vaqtinchalik holat `nil` UploadType uzatadi, ya'ni har bir chaqiruv xato qaytaradi. Bu ataylab — endpoint Task 3 gacha ishlamaydi, lekin **ochiq qolgandan ko'ra ishlamagani xavfsizroq**.

- [ ] **Step 8: Wire va build**

Run: `make wire-build && go build ./... && go vet ./...`
Expected: hammasi toza

- [ ] **Step 9: To'liq to'plam**

Run: `go test ./... -count=1`
Expected: 79 + 5 = 84 test o'tadi

- [ ] **Step 10: Commit**

```bash
git add src/infrastructure/storage/post_policy.go test/infrastructure/storage/post_policy_test.go \
        src/core/domain/ports/storage/file_storage.go src/infrastructure/storage/minio_storage.go \
        src/core/application/usecase/uploadusecases/upload_file_usecase.go cmd/container/container.go
git commit -m "feat(upload): PresignedPutObject o'rniga POST policy (hajm va tur MinIO tomonidan majburlanadi)"
```

---

## Task 3: Presigned yuklash use case, papka→bucket va API kontrakti

**Files:**
- Modify: `src/core/application/usecase/uploadusecases/upload_file_usecase.go` (to'liq qayta yoziladi)
- Test: `test/core/application/usecase/uploadusecases/upload_file_usecase_test.go`
- Modify: `src/entrypoint/presentation/handlers/upload/schema/presigned_put_url_schema.go`
- Modify: `src/entrypoint/presentation/handlers/upload/presigned_put_url_handler.go`

**Interfaces:**
- Consumes: Task 1 dan `enum.LookupUploadType`; Task 2 dan `PostPolicyPresignedUrl`
- Produces:
  - `type PresignedUpload struct { URL string; Fields map[string]string; ObjectName string }`
  - `func (u *UploadFileUseCase) Execute(ctx context.Context, folder enum.StorageFolder, fileName string) (*PresignedUpload, error)`
  - `func BucketForFolder(folder enum.StorageFolder) (enum.Bucket, bool)`

Bu task ekspluatatsiya zanjirining birinchi bo'g'inini uzadi. Ikkita xatti-harakat **majburiy**: (1) ruxsat etilmagan kengaytma uchun storage **umuman chaqirilmaydi**, (2) mijoz yuborgan papka yo'li (`path.Dir`) tashlab yuboriladi.

- [ ] **Step 1: Failing testni yozish**

`test/core/application/usecase/uploadusecases/upload_file_usecase_test.go`:

```go
package uploadusecases_test

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/uploadusecases"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/infrastructure/config"
)

// fakeStorage faqat PostPolicyPresignedUrl ni kuzatadi; qolgan metodlar
// chaqirilsa test panic bilan yiqiladi — bu ataylab.
type fakeStorage struct {
	storage.FileStorage
	calls      int
	gotBucket  enum.Bucket
	gotObject  string
	gotType    *enum.UploadType
	gotMaxSize int64
	err        error
}

func (f *fakeStorage) PostPolicyPresignedUrl(
	_ context.Context, bucket enum.Bucket, objectName string,
	ut *enum.UploadType, maxSize int64,
) (*url.URL, map[string]string, error) {
	f.calls++
	f.gotBucket = bucket
	f.gotObject = objectName
	f.gotType = ut
	f.gotMaxSize = maxSize
	if f.err != nil {
		return nil, nil, f.err
	}
	u, _ := url.Parse("https://storage.example/slibbucketprivate")
	return u, map[string]string{"key": objectName, "Content-Type": ut.ContentType}, nil
}

func newUseCase(s *fakeStorage) *uploadusecases.UploadFileUseCase {
	return uploadusecases.NewUploadFileUseCase(s, &config.Config{UploadedFileMaxSize: 10 << 20})
}

func TestDangerousExtensionNeverReachesStorage(t *testing.T) {
	for _, name := range []string{"evil.html", "evil.svg", "shell.php", "a.exe", "noext"} {
		s := &fakeStorage{}
		_, err := newUseCase(s).Execute(context.Background(), enum.FolderArticles, name)
		if err == nil {
			t.Errorf("%s: xato kutilgandi", name)
		}
		if !errors.Is(err, response.InvalidFileError) {
			t.Errorf("%s: InvalidFileError kutilgandi, %v keldi", name, err)
		}
		if s.calls != 0 {
			t.Errorf("%s: storage chaqirilmasligi kerak edi, %d marta chaqirildi", name, s.calls)
		}
	}
}

func TestAllowedExtensionBuildsRequest(t *testing.T) {
	s := &fakeStorage{}
	got, err := newUseCase(s).Execute(context.Background(), enum.FolderArticles, "maqola.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.calls != 1 {
		t.Fatalf("storage bir marta chaqilishi kerak edi, %d", s.calls)
	}
	if s.gotType == nil || s.gotType.ContentType != "application/pdf" {
		t.Errorf("application/pdf turi uzatilishi kerak edi, %v", s.gotType)
	}
	if s.gotMaxSize != 10<<20 {
		t.Errorf("maksimal hajm uzatilishi kerak edi, %d keldi", s.gotMaxSize)
	}
	if got.Fields == nil {
		t.Error("javobda POST policy form maydonlari bo'lishi kerak edi")
	}
	if got.ObjectName != s.gotObject {
		t.Errorf("javobdagi ObjectName storage'ga uzatilgani bilan mos kelishi kerak")
	}
}

func TestClientDirectoryPathIsStripped(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderArticles, "../../../etc/passwd.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if strings.Contains(s.gotObject, "..") {
		t.Errorf("obyekt kalitida '..' qolmasligi kerak: %q", s.gotObject)
	}
	if !strings.HasPrefix(s.gotObject, string(enum.FolderArticles)+"/") {
		t.Errorf("obyekt kaliti %q papkasi bilan boshlanishi kerak: %q", enum.FolderArticles, s.gotObject)
	}
	if !strings.HasSuffix(s.gotObject, ".pdf") {
		t.Errorf("kengaytma saqlanishi kerak: %q", s.gotObject)
	}
}

func TestObjectNamesAreUnique(t *testing.T) {
	s1, s2 := &fakeStorage{}, &fakeStorage{}
	_, _ = newUseCase(s1).Execute(context.Background(), enum.FolderArticles, "a.pdf")
	_, _ = newUseCase(s2).Execute(context.Background(), enum.FolderArticles, "a.pdf")
	if s1.gotObject == s2.gotObject {
		t.Errorf("bir xil nom uchun kalitlar farq qilishi kerak edi: %q", s1.gotObject)
	}
}

func TestPublicFolderGoesToPublicBucket(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderNews, "muqova.png")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.gotBucket != enum.BucketPublic {
		t.Errorf("news papkasi PUBLIC bucket'ga tushishi kerak, %q keldi", s.gotBucket)
	}
}

func TestPrivateFolderGoesToPrivateBucket(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderApplications, "ariza.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.gotBucket != enum.BucketPrivate {
		t.Errorf("applications papkasi PRIVATE bucket'ga tushishi kerak, %q keldi", s.gotBucket)
	}
}

func TestUnknownFolderIsRejected(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.StorageFolder("../secrets"), "a.pdf")
	if !errors.Is(err, response.InvalidFileError) {
		t.Errorf("InvalidFileError kutilgandi, %v keldi", err)
	}
	if s.calls != 0 {
		t.Errorf("noma'lum papka uchun storage chaqirilmasligi kerak edi")
	}
}

func TestStorageErrorIsPropagated(t *testing.T) {
	boom := errors.New("minio ishlamayapti")
	s := &fakeStorage{err: boom}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderArticles, "a.pdf")
	if !errors.Is(err, boom) {
		t.Errorf("storage xatosi uzatilishi kerak edi, %v keldi", err)
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

Run: `go test ./test/core/application/usecase/uploadusecases/... -count=1`
Expected: FAIL — `NewUploadFileUseCase` ikkita argument qabul qilmaydi, `PresignedUpload` aniqlanmagan

- [ ] **Step 3: Use case'ni qayta yozish**

`src/core/application/usecase/uploadusecases/upload_file_usecase.go` — faylni to'liq almashtiring:

```go
package uploadusecases

import (
	"context"
	"path/filepath"
	"strings"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
)

// PresignedUpload mijozga qaytariladigan yuklash ma'lumoti.
// Fields — POST policy form maydonlari; mijoz ularni faylning OLDIDA
// yuborishi shart (S3 talabi).
type PresignedUpload struct {
	URL        string
	Fields     map[string]string
	ObjectName string
}

type UploadFileUseCase struct {
	storage storage.FileStorage
	maxSize int64
}

// @inject
// Konstruktor *config.Config qabul qiladi, primitiv int64 emas: wire
// primitiv turni o'zi hal qila olmaydi. Bu kodbazadagi mavjud uslub —
// MinioStorage ham shunday quriladi.
func NewUploadFileUseCase(storage storage.FileStorage, env *config.Config) *UploadFileUseCase {
	return &UploadFileUseCase{storage: storage, maxSize: int64(env.UploadedFileMaxSize)}
}

// folderBuckets papkani bucket'ga bog'laydi. Bucket'ni MIJOZ EMAS, server
// tanlaydi — aks holda mijoz maxfiy hujjatni ommaviy bucket'ga yuklay olardi.
var folderBuckets = map[enum.StorageFolder]enum.Bucket{
	enum.FolderNews:                   enum.BucketPublic,
	enum.FolderImages:                 enum.BucketPublic,
	enum.FolderArticles:               enum.BucketPrivate,
	enum.FolderArticle:                enum.BucketPrivate,
	enum.FolderApplications:           enum.BucketPrivate,
	enum.FolderCertificates:           enum.BucketPrivate,
	enum.FolderAntiPlagCertificate:    enum.BucketPrivate,
	enum.FolderArticleExpertConclusion: enum.BucketPrivate,
	enum.FolderTickets:                enum.BucketPrivate,
	enum.FolderSpellCheck:             enum.BucketPrivate,
}

// BucketForFolder papkaga mos bucket'ni qaytaradi.
// Noma'lum papka uchun ("", false) — chaqiruvchi rad etishi shart.
func BucketForFolder(folder enum.StorageFolder) (enum.Bucket, bool) {
	b, ok := folderBuckets[folder]
	return b, ok
}

// Execute yuklash uchun imzolangan POST policy qaytaradi.
//
// Kengaytma allow-list'da bo'lmasa storage UMUMAN chaqirilmaydi — yuklash
// boshlanishidan oldin to'xtatiladi. Bu ekspluatatsiya zanjirining birinchi
// bo'g'inini uzadi.
func (this *UploadFileUseCase) Execute(
	ctx context.Context,
	folder enum.StorageFolder,
	fileName string,
) (*PresignedUpload, error) {
	bucket, ok := BucketForFolder(folder)
	if !ok {
		return nil, response.InvalidFileError
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	uploadType, ok := enum.LookupUploadType(ext)
	if !ok {
		return nil, response.InvalidFileError
	}

	objectName := this.buildObjectName(folder, fileName, ext)

	presignedURL, fields, err := this.storage.PostPolicyPresignedUrl(
		ctx, bucket, objectName, uploadType, this.maxSize,
	)
	if err != nil {
		return nil, err
	}

	return &PresignedUpload{
		URL:        presignedURL.String(),
		Fields:     fields,
		ObjectName: objectName,
	}, nil
}

// buildObjectName obyekt kalitini SERVERDA quradi.
//
// Mijoz yuborgan papka yo'li ataylab tashlab yuboriladi: filepath.Base
// faqat fayl nomini qoldiradi, ya'ni "../../etc/passwd.pdf" -> "passwd.pdf".
// Eski kod path.Dir() ni saqlar edi va shu sababli mijoz obyekt kalitining
// tuzilishini boshqara olardi.
func (this *UploadFileUseCase) buildObjectName(folder enum.StorageFolder, fileName, ext string) string {
	base := filepath.Base(fileName)
	clean := utils.SanitizeFilename(strings.TrimSuffix(base, filepath.Ext(base)))
	if clean == "" || clean == "." || clean == ".." {
		clean = "file"
	}
	return string(folder) + "/" + clean + "_" + utils.RandomHex(4) + ext
}
```

Diqqat: `folderBuckets` xaritasidagi kalit nomlari `src/core/domain/entity/enum/storage.go` dagi konstantalar bilan **aniq** mos kelishi shart — u yerda o'nta papka aniqlangan va xaritada ham o'ntasi bo'lishi kerak, aks holda tushib qolgan papka uchun yuklash jimgina rad etiladi.

- [ ] **Step 4: Testni ishga tushirib, o'tishini tasdiqlash**

Run: `go test ./test/core/application/usecase/uploadusecases/... -count=1 -v`
Expected: PASS — 8 ta test

- [ ] **Step 5: Schema'ni yangilash**

`src/entrypoint/presentation/handlers/upload/schema/presigned_put_url_schema.go` — to'liq almashtiring:

```go
package schema

import "slib.uz/src/core/domain/entity/enum"

// PresignedPutUrlRequest — yuklash so'rovi.
// Bucket maydoni ATAYLAB olib tashlandi: uni endi server papkadan aniqlaydi.
type PresignedPutUrlRequest struct {
	FileName string              `json:"filename"`
	Folder   enum.StorageFolder  `json:"folder"`
}

// PresignedPutUrlResponse — POST policy yuklash ma'lumoti.
// Mijoz Fields maydonlarini faylning OLDIDA yuborishi shart.
type PresignedPutUrlResponse struct {
	URL        string            `json:"url"`
	Fields     map[string]string `json:"fields"`
	ObjectName string            `json:"object_name"`
}

func NewPresignedPutUrlResponse(url string, fields map[string]string, objectName string) *PresignedPutUrlResponse {
	return &PresignedPutUrlResponse{URL: url, Fields: fields, ObjectName: objectName}
}
```

- [ ] **Step 6: Handler'ni yangilash**

`src/entrypoint/presentation/handlers/upload/presigned_put_url_handler.go` — `Handle` metodini va swagger annotatsiyasini almashtiring:

```go
// Handle godoc
// @Tags upload
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body schema.PresignedPutUrlRequest true "Upload File Request"
// @Success 200 {object} schema.PresignedPutUrlResponse
// @Failure 400 {object} response.Response "kengaytma ruxsat etilmagan yoki papka noma'lum"
// @Router /upload/presigned/put-url [post]
func (this *PresignedPutUrlHandler) Handle(c echo.Context) error {
	body, err := context.GetBody[schema.PresignedPutUrlRequest](c)
	if err != nil {
		return err
	}

	result, err := this.uc.Execute(c.Request().Context(), body.Folder, body.FileName)
	if err != nil {
		return err
	}

	return c.JSON(200, schema.NewPresignedPutUrlResponse(result.URL, result.Fields, result.ObjectName))
}
```

`response.Response` — kodbazada `@Failure` uchun ishlatiladigan tur (masalan `handlers/journal/journal_detail_handler.go:28`). Import qo'shish shart emas: swaggo annotatsiyalari izohda yashaydi va kompilyatsiyaga kirmaydi.

- [ ] **Step 7: Wire, docs va to'liq to'plam**

Run:
```bash
make wire-build && make generate-docs && go build ./... && go vet ./... && go test ./... -count=1
```
Expected: build va vet toza, 84 + 8 = 92 test o'tadi

- [ ] **Step 8: Generatorlar idempotentligini tekshirish**

Run: `make wire-build && make generate-docs && git diff --stat`
Expected: bo'sh chiqish

- [ ] **Step 9: Commit**

```bash
git add src/core/application/usecase/uploadusecases/upload_file_usecase.go \
        test/core/application/usecase/uploadusecases/upload_file_usecase_test.go \
        src/entrypoint/presentation/handlers/upload/ \
        cmd/container/container.go src/entrypoint/presentation/docs/
git commit -m "feat(upload): allow-list URL berish paytida, bucket va obyekt kalitini server belgilaydi"
```

---

## Task 4: `/upload/file` — magic-byte va hajm tekshiruvi

**Files:**
- Modify: `src/core/application/usecase/uploadusecases/upload_tempfile_usecase.go`
- Test: `test/core/application/usecase/uploadusecases/upload_tempfile_usecase_test.go`
- Modify: `src/infrastructure/storage/minio_storage.go:74-76`

**Interfaces:**
- Consumes: Task 1 dan `enum.LookupUploadType`, `UploadType.MatchesMagic`, `enum.MagicPrefixLen`
- Produces: `func (u *UploadTempFileUseCase) Execute(file *multipart.FileHeader) (string, error)` — imzo o'zgarmaydi, xatti-harakat o'zgaradi

Bu yo'lda baytlar bizda, shuning uchun tekshiruv presigned yo'ldan kuchliroq: kengaytma **va** imzo mos kelishi shart.

- [ ] **Step 1: Failing testni yozish**

`test/core/application/usecase/uploadusecases/upload_tempfile_usecase_test.go`:

```go
package uploadusecases_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/uploadusecases"
	"slib.uz/src/core/domain/ports/storage"
)

type fakeTempStorage struct {
	storage.FileStorage
	calls int
}

func (f *fakeTempStorage) SaveTempFile(_ *multipart.FileHeader) (string, error) {
	f.calls++
	return "uploads/saqlandi.pdf", nil
}

// makeFileHeader berilgan nom va mazmun bilan haqiqiy multipart.FileHeader yasaydi.
func makeFileHeader(t *testing.T, name string, content []byte) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	r := multipart.NewReader(&buf, w.Boundary())
	form, err := r.ReadForm(int64(len(content)) + 1024)
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}
	return form.File["file"][0]
}

func TestRealPdfIsAccepted(t *testing.T) {
	s := &fakeTempStorage{}
	uc := uploadusecases.NewUploadTempFileUseCase(s)

	fh := makeFileHeader(t, "maqola.pdf", []byte("%PDF-1.7\nhaqiqiy pdf mazmuni"))
	if _, err := uc.Execute(fh); err != nil {
		t.Fatalf("haqiqiy PDF qabul qilinishi kerak edi: %v", err)
	}
	if s.calls != 1 {
		t.Errorf("storage bir marta chaqilishi kerak edi, %d", s.calls)
	}
}

func TestHtmlBytesNamedPdfAreRejected(t *testing.T) {
	s := &fakeTempStorage{}
	uc := uploadusecases.NewUploadTempFileUseCase(s)

	fh := makeFileHeader(t, "zararsiz.pdf", []byte("<html><script>alert(1)</script></html>"))
	_, err := uc.Execute(fh)
	if !errors.Is(err, response.InvalidFileError) {
		t.Errorf("InvalidFileError kutilgandi, %v keldi", err)
	}
	if s.calls != 0 {
		t.Errorf("rad etilgan fayl saqlanmasligi kerak edi, %d marta saqlandi", s.calls)
	}
}

func TestDisallowedExtensionIsRejected(t *testing.T) {
	s := &fakeTempStorage{}
	uc := uploadusecases.NewUploadTempFileUseCase(s)

	fh := makeFileHeader(t, "evil.html", []byte("<html></html>"))
	_, err := uc.Execute(fh)
	if !errors.Is(err, response.InvalidFileError) {
		t.Errorf("InvalidFileError kutilgandi, %v keldi", err)
	}
	if s.calls != 0 {
		t.Errorf("storage chaqirilmasligi kerak edi")
	}
}

func TestTruncatedFileIsRejectedNotPanicked(t *testing.T) {
	s := &fakeTempStorage{}
	uc := uploadusecases.NewUploadTempFileUseCase(s)

	fh := makeFileHeader(t, "qisqa.png", []byte{0x89})
	_, err := uc.Execute(fh)
	if !errors.Is(err, response.InvalidFileError) {
		t.Errorf("InvalidFileError kutilgandi, %v keldi", err)
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

Run: `go test ./test/core/application/usecase/uploadusecases/... -run 'Pdf|Html|Disallowed|Truncated' -count=1`
Expected: FAIL — hozircha hech qanday tekshiruv yo'q, HTML `.pdf` sifatida qabul qilinadi

- [ ] **Step 3: Tekshiruvni implementatsiya qilish**

`src/core/application/usecase/uploadusecases/upload_tempfile_usecase.go` — to'liq almashtiring:

```go
package uploadusecases

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/storage"
)

type UploadTempFileUseCase struct {
	storage storage.FileStorage
}

// @inject
func NewUploadTempFileUseCase(storage storage.FileStorage) *UploadTempFileUseCase {
	return &UploadTempFileUseCase{storage: storage}
}

// Execute faylni tekshiradi va vaqtinchalik saqlaydi.
//
// Bu yo'lda baytlar bizda, shuning uchun tekshiruv presigned yo'ldan
// kuchliroq: kengaytma allow-list'da bo'lishi VA fayl imzosi o'sha turga
// mos kelishi shart. Ikkalasidan biri mos kelmasa fayl saqlanmaydi.
func (this *UploadTempFileUseCase) Execute(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", response.InvalidFileError
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	uploadType, ok := enum.LookupUploadType(ext)
	if !ok {
		return "", response.InvalidFileError
	}

	head, err := readHead(file, enum.MagicPrefixLen)
	if err != nil {
		return "", response.InvalidFileError
	}
	if !uploadType.MatchesMagic(head) {
		return "", response.InvalidFileError
	}

	return this.storage.SaveTempFile(file)
}

// readHead faylning birinchi n baytini o'qiydi.
// Fayl n dan qisqa bo'lsa bor baytlarni qaytaradi — bu holda imzo
// mos kelmaydi va fayl rad etiladi.
func readHead(file *multipart.FileHeader, n int) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buf := make([]byte, n)
	read, err := io.ReadFull(src, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}
	return buf[:read], nil
}
```

- [ ] **Step 4: Testni ishga tushirib, o'tishini tasdiqlash**

Run: `go test ./test/core/application/usecase/uploadusecases/... -count=1 -v`
Expected: PASS — 8 (Task 3) + 4 (Task 4) = 12 ta test

- [ ] **Step 5: Hajm chegarasini jonlantirish**

`src/infrastructure/storage/minio_storage.go:74-76` — izohga olingan blokni tiklang:

```go
	if file.Size > int64(this.env.UploadedFileMaxSize) {
		return "", response.EntityToLargeError
	}
```

Diqqat: eski izohda `config.GetEnv()` yozilgan edi; `MinioStorage` ning `env` maydoni allaqachon mavjud, shuning uchun `this.env` ishlatiladi.

- [ ] **Step 6: To'liq to'plam**

Run: `go build ./... && go vet ./... && go test ./... -count=1`
Expected: 92 + 4 = 96 test o'tadi

- [ ] **Step 7: Commit**

```bash
git add src/core/application/usecase/uploadusecases/upload_tempfile_usecase.go \
        test/core/application/usecase/uploadusecases/upload_tempfile_usecase_test.go \
        src/infrastructure/storage/minio_storage.go
git commit -m "feat(upload): /upload/file magic-byte va hajm bo'yicha tekshiriladi"
```

---

## Task 5: Presigned GET javob sarlavhalari va base64 yo'li

**Files:**
- Modify: `src/infrastructure/storage/minio_storage.go` (`PresignedURL` ~177, `PutObject` ~167)
- Modify: `src/core/application/usecase/uploadusecases/upload_base64_file_usecase.go:67`
- Test: `test/infrastructure/storage/response_params_test.go`

**Interfaces:**
- Consumes: Task 1 dan `enum.LookupUploadType`, `enum.DispositionAttachment`
- Produces: `func ResponseParamsForObject(objectName string) url.Values`

`PresignedURL` ning uchta chaqiruv joyi bor (`presigned_url_usecase.go:34`, `application_article_invoke_usecase.go:48`, `article_file_repository_impl.go:27`) — hammasi `BucketPrivate` bilan. Sarlavhalarni obyekt nomining kengaytmasidan aniqlash **imzoni o'zgartirmaydi**, ya'ni uchala chaqiruv joyiga tegilmaydi.

- [ ] **Step 1: Failing testni yozish**

`test/infrastructure/storage/response_params_test.go`:

```go
package storage_test

import (
	"strings"
	"testing"

	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/storage"
)

func TestKnownDocumentGetsAttachment(t *testing.T) {
	p := storage.ResponseParamsForObject("articles/maqola_1234.pdf")

	if got := p.Get("response-content-type"); got != "application/pdf" {
		t.Errorf("response-content-type: application/pdf kutilgandi, %q keldi", got)
	}
	disp := p.Get("response-content-disposition")
	if !strings.HasPrefix(disp, enum.DispositionAttachment) {
		t.Errorf("attachment kutilgandi, %q keldi", disp)
	}
	if !strings.Contains(disp, "maqola_1234.pdf") {
		t.Errorf("disposition'da fayl nomi bo'lishi kerak edi, %q keldi", disp)
	}
}

func TestKnownImageGetsInline(t *testing.T) {
	p := storage.ResponseParamsForObject("news/muqova_1234.png")

	if got := p.Get("response-content-type"); got != "image/png" {
		t.Errorf("image/png kutilgandi, %q keldi", got)
	}
	if disp := p.Get("response-content-disposition"); !strings.HasPrefix(disp, enum.DispositionInline) {
		t.Errorf("inline kutilgandi, %q keldi", disp)
	}
}

// Eski, allow-list'dan oldin yuklangan fayllar eng xavfsiz holatga tushishi kerak.
func TestUnknownExtensionFallsBackToSafestValues(t *testing.T) {
	for _, name := range []string{"legacy/eski.html", "legacy/noext", "legacy/x.svg"} {
		p := storage.ResponseParamsForObject(name)

		if got := p.Get("response-content-type"); got != "application/octet-stream" {
			t.Errorf("%s: application/octet-stream kutilgandi, %q keldi", name, got)
		}
		if disp := p.Get("response-content-disposition"); !strings.HasPrefix(disp, enum.DispositionAttachment) {
			t.Errorf("%s: attachment kutilgandi, %q keldi", name, disp)
		}
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

Run: `go test ./test/infrastructure/storage/... -count=1`
Expected: FAIL — `undefined: storage.ResponseParamsForObject`

- [ ] **Step 3: Implementatsiya qilish**

`src/infrastructure/storage/post_policy.go` fayliga qo'shing:

```go
// ResponseParamsForObject presigned GET uchun javob sarlavhalarini quradi.
//
// Bu MAXFIY bucket uchun asosiy himoya: brauzer faylni ijro etmasligi
// kafolatlanadi. Ommaviy bucket anonim o'qishga ochiq bo'lgani uchun bu
// parametrlar u yerda qo'llanmaydi — u yerda obyekt metadatasi ishlaydi.
//
// Allow-list'da yo'q kengaytma (allow-list joriy etilgunicha yuklangan
// fayllar) eng xavfsiz holatga tushadi: octet-stream + attachment.
func ResponseParamsForObject(objectName string) url.Values {
	params := make(url.Values)

	contentType := "application/octet-stream"
	disposition := enum.DispositionAttachment

	ext := strings.ToLower(filepath.Ext(objectName))
	if ut, ok := enum.LookupUploadType(ext); ok {
		contentType = ut.ContentType
		disposition = ut.Disposition
	}

	params.Set("response-content-type", contentType)
	params.Set("response-content-disposition",
		disposition+`; filename="`+filepath.Base(objectName)+`"`)

	return params
}
```

Importlarni qo'shing: `"net/url"`, `"path/filepath"`, `"strings"`.

- [ ] **Step 4: Testni ishga tushirib, o'tishini tasdiqlash**

Run: `go test ./test/infrastructure/storage/... -count=1 -v`
Expected: PASS — 5 (Task 2) + 3 (Task 5) = 8 ta test

- [ ] **Step 5: `PresignedURL` ni ulash**

`src/infrastructure/storage/minio_storage.go` — `PresignedURL` ichidagi `reqParams := make(url.Values)` qatorini almashtiring:

```go
	reqParams := ResponseParamsForObject(objectName)
```

- [ ] **Step 6: `PutObject` Content-Type manbasini allow-list'ga o'tkazish**

`minio_storage.go` dagi `PutObject` metodida `ContentType: utils.GetContentType(objectName)` o'rniga:

```go
	contentType := "application/octet-stream"
	disposition := enum.DispositionAttachment
	if ut, ok := enum.LookupUploadType(strings.ToLower(filepath.Ext(objectName))); ok {
		contentType = ut.ContentType
		disposition = ut.Disposition
	}
```

va `minio.PutObjectOptions` ga:

```go
		ContentType:        contentType,
		ContentDisposition: disposition,
```

Sabab: `utils.GetContentType` `mime.TypeByExtension` ga tayanadi, u allow-list'dan kengroq va shuning uchun xavfsizlik qarori uchun yaroqsiz. Bitta yuklash yo'li boshqa qoidada qolmasligi kerak.

- [ ] **Step 7: Base64 yo'lidagi PDF tekshiruvini umumlashtirish**

`src/core/application/usecase/uploadusecases/upload_base64_file_usecase.go:67` dagi qo'lda yozilgan tekshiruvni:

```go
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		return fmt.Errorf("invalid file format: only PDF files are allowed")
	}
```

allow-list'ga tayangan variantga almashtiring:

```go
	pdfType, ok := enum.LookupUploadType(".pdf")
	if !ok || !pdfType.MatchesMagic(data) {
		return response.InvalidFileError
	}
```

`enum` va `response` importlarini qo'shing; `bytes` va `fmt` endi kerak bo'lmasa olib tashlang. Bu PDF imzosining ikkinchi nusxasini yo'q qiladi.

- [ ] **Step 8: To'liq to'plam**

Run: `go build ./... && go vet ./... && go test ./... -count=1`
Expected: 96 + 3 = 99 test o'tadi

- [ ] **Step 9: Commit**

```bash
git add src/infrastructure/storage/post_policy.go src/infrastructure/storage/minio_storage.go \
        src/core/application/usecase/uploadusecases/upload_base64_file_usecase.go \
        test/infrastructure/storage/response_params_test.go
git commit -m "feat(upload): presigned GET javob sarlavhalari va allow-list'ga tayangan Content-Type"
```

---

## Task 6: HTTP xavfsizlik sarlavhalari

**Files:**
- Modify: `src/entrypoint/presentation/app/app.go:242`
- Test: `test/entrypoint/presentation/app/security_headers_test.go`

**Interfaces:**
- Consumes: hech narsa
- Produces: `func SecurityHeadersMiddleware() echo.MiddlewareFunc`

Halol baho: bu backend JSON API qaytaradi, sahifa emas. `nosniff` real foyda beradi — API binar javoblar ham qaytaradi (statistika Excel eksporti, base64 yuklab olish). CSP va `X-Frame-Options` bu yerda kichik qiymatga ega, lekin arzon.

Bu **ekspertiza 2.2.8 (Clickjacking) ni YOPMAYDI** — u frontend sahifalariga tegishli.

- [ ] **Step 1: Failing testni yozish**

`test/entrypoint/presentation/app/security_headers_test.go`:

```go
package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	apppkg "slib.uz/src/entrypoint/presentation/app"
)

func TestSecurityHeadersArePresent(t *testing.T) {
	e := echo.New()
	e.Use(apppkg.SecurityHeadersMiddleware())
	e.GET("/probe", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("200 kutilgandi, %d keldi", rec.Code)
	}

	want := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
	}
	for header, expected := range want {
		if got := rec.Header().Get(header); got != expected {
			t.Errorf("%s: %q kutilgandi, %q keldi", header, expected, got)
		}
	}

	if csp := rec.Header().Get("Content-Security-Policy"); csp == "" {
		t.Error("Content-Security-Policy o'rnatilishi kerak edi")
	}
}

// Eski X-XSS-Protection zamonaviy brauzerlarda o'chirilgan va yoqilganda
// o'zi zaiflik manbai bo'lishi mumkin.
func TestLegacyXssProtectionIsDisabled(t *testing.T) {
	e := echo.New()
	e.Use(apppkg.SecurityHeadersMiddleware())
	e.GET("/probe", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))

	if got := rec.Header().Get("X-XSS-Protection"); got != "" && got != "0" {
		t.Errorf("X-XSS-Protection o'chirilgan bo'lishi kerak edi, %q keldi", got)
	}
}
```

- [ ] **Step 2: Testni ishga tushirib, yiqilishini tasdiqlash**

Run: `go test ./test/entrypoint/presentation/app/... -count=1`
Expected: FAIL — `undefined: app.SecurityHeadersMiddleware`

- [ ] **Step 3: Middleware'ni yozish**

`src/entrypoint/presentation/app/app.go` fayliga qo'shing (eksport qilingan funksiya sifatida — test uni `app_test` paketidan chaqiradi):

```go
// SecurityHeadersMiddleware barcha javoblarga xavfsizlik sarlavhalarini qo'shadi.
//
// nosniff bu yerda eng qimmatlisi: API binar javoblar ham qaytaradi
// (statistika Excel eksporti, base64 yuklab olish) va brauzer ularning
// turini taxmin qilishga urinmasligi kerak.
//
// XSSProtection "0" ATAYLAB: eski X-XSS-Protection sarlavhasi zamonaviy
// brauzerlarda o'chirilgan va yoqilganda o'zi zaiflik manbai bo'lishi mumkin.
func SecurityHeadersMiddleware() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "0",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
	})
}
```

- [ ] **Step 4: Middleware'ni ro'yxatdan o'tkazish**

`initMiddlewares` ichida, `debugLogMiddleware` dan keyin, CORS dan oldin:

```go
	this.echo.Use(SecurityHeadersMiddleware())
```

- [ ] **Step 5: Testni ishga tushirib, o'tishini tasdiqlash**

Run: `go test ./test/entrypoint/presentation/app/... -count=1 -v`
Expected: PASS — 2 ta test

- [ ] **Step 6: To'liq to'plam**

Run: `go build ./... && go vet ./... && go test ./... -count=1`
Expected: 99 + 2 = 101 test o'tadi

- [ ] **Step 7: Commit**

```bash
git add src/entrypoint/presentation/app/app.go \
        test/entrypoint/presentation/app/security_headers_test.go
git commit -m "feat(security): javoblarga nosniff, CSP va X-Frame-Options qo'shildi"
```

---

## Yakuniy tekshiruv

Barcha tasklar tugagach:

- [ ] `go build ./... && go vet ./...` — toza
- [ ] `go test ./... -count=1` — 101 test o'tadi
- [ ] `make wire-build && make generate-docs && git diff --stat` — bo'sh
- [ ] `git diff develop..HEAD -- go.mod go.sum` — bo'sh (yangi bog'liqlik yo'q)
- [ ] `grep -rn "zap" src/core/` — natija yo'q
- [ ] `grep -rn "GetContentType" src/` — faqat `utils/filename.go` dagi ta'rifi qoladi, yuklash yo'llarida chaqiruv yo'q

## Deploy eslatmasi

7-bo'lim (spec) **buzuvchi o'zgarish**ni tavsiflaydi: `PUT` → `multipart POST`, `file` maydoni oxirgi bo'lishi shart. Backend frontend'dan oldin chiqsa, yuklash ishlamay qoladi. Deploy frontend bilan muvofiqlashtirilishi shart.

Deploy oldidan ishlab chiqarishdagi ommaviy bucket policy'si haqiqatan anonim `s3:GetObject` ga ochiqligini tasdiqlang — dizayn shu faktga tayanadi.
