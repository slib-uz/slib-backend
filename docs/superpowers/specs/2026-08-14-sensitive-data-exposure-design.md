# Muhim ma'lumotlarni ochiqlanishi (CWE-200) — dizayn spetsifikatsiyasi

**Sana:** 2026-08-14
**Audit topilmasi:** 2.2.5 Muhim turdagi ma'lumotlarni ochiqlanishi, xavflilik darajasi **Yuqori**
**CWE:** CWE-200 — Exposure of Sensitive Information to an Unauthorized Actor
**OWASP Top 2025:** 1-o'rin (Broken Access Control)
**Shox:** `feature/sensitive-data`, `feature/sql-injection` ustiga (zanjirning 4-halqasi)

---

## 1. Zaiflik

Ekspertiza `GET /api/article-application/detail/{applicationId}` endpointini ko'rsatgan.
Oddiy (past huquqli) foydalanuvchi tokeni bilan istalgan `applicationId` ga murojaat
qilib, boshqa foydalanuvchilarning shaxsiy ma'lumotlarini oladi:

- **PINFL** (`pin`) — O'zbekiston shaxsiy identifikatsiya raqami, eng maxfiy
- telefon raqami (`phone_number`)
- tug'ilgan sana (`birth_date`)
- to'liq ism, email

Muallif, co-author'lar va reviewer'lar — hammasining. Hisobotning 20-rasmi
`app_id` ni `1..N` aylantirib avtomatik sizib chiqishni ko'rsatadi.

### 1.1. Bu bitta endpoint emas — naqsh

`UserEntity` (PINFL, telefon, birth_date, email bilan) domenda 10+ entity'da ichma-ich
ishlatiladi. `ApplicationResponseEntity` da u **uch yo'ldan** chiqadi:

- `.User` — ariza egasi
- `.Article.CoAuthors[].User` — co-author'lar (`AuthorEntity.User *UserEntity`)
- `.ReviewStages[].Reviewer` — reviewer'lar (`*UserEntity`)

Ikkinchi oshkor nuqta: `journal_config` handler'ining `CreatorResponse` schemasi
`Pin` maydonini **ataylab** qo'shadi (`journal_config_response.go:58`).

Uchinchi oshkor nuqta: `GET /api/journal/{journalId}/reviewers` — umuman himoyasiz
(anonim kirish), reviewer'larning telefon ro'yxatini qaytaradi.

### 1.2. Sidiruvning tuzatilgan natijasi

Avtomatik sidiruv bir necha endpointni "ochiq API" deb belgiladi — bu **noto'g'ri**.
Sidiruvchi guruh darajasidagi middleware'ni ko'rmadi. Haqiqiy holat (`app.go:284-309`):

| Endpoint | Sidiruvchi dedi | Haqiqat |
|---|---|---|
| `/journal-applications/applications/{id}` | ochiq | `AuthenticatedPermission` — login shart |
| `/users/find` | ochiq | `AuthenticatedPermission` — login shart |
| `/journal/{journalId}/reviewers` | ochiq | **haqiqatan ochiq** — permission yo'q |

Xulosa: avtomatik sidiruvga ko'r-ko'rona ishonmaslik kerak; permission har doim
`app.go` dagi guruh registratsiyasidan tasdiqlanadi.

---

## 2. Muammoning ikki o'lchami

**A. Ma'lumot minimallashtirish** (nima qaytadi) — asosiy himoya.
PINFL/telefon/birth_date/email API javoblarida keraksiz oshkor bo'ladi. Bu bir joyda
(entity + markazlashtirilgan yordamchi) hal qilinadi va hatto avtorizatsiya nomukammal
bo'lsa ham eng maxfiy ma'lumotni to'sadi.

**B. Avtorizatsiya** (kim ko'radi) — qo'shimcha himoya.
Bitta aniq ochiq teshik: `/journal/{id}/reviewers`.

**Qamrovdan tashqarida** (foydalanuvchi qarori): `/article-application/detail` va
`/journal-applications/{id}` ning chuqur workflow-egalik guard'i. Ikkalasi ham review
workflow bilimini talab qiladi va ikkalasidan ham PINFL/birth_date olib tashlangandan
keyin maxfiy ma'lumot sizmaydi.

---

## 3. Yechim — A o'lcham

### 3.1. PINFL va birth_date'ni butunlay olib tashlash

Tekshiruv: `Pin` va `BirthDate` hech qanday frontend ekranida kerak emas.
`UserProfileEntity` (foydalanuvchi o'z profili) ularni o'z ichiga olmaydi. Backend
mantiqda `.Pin` faqat ichki qidiruvda (`user_repository_impl.go:74` `GetByPin`) va
mapper'larda ishlatiladi — hech biri JSON round-trip qilmaydi.

**O'zgarishlar:**

1. `src/core/domain/entity/user_entity.go` — `Pin` va `BirthDate` maydonlariga
   `json:"-"` tegi. Bir joyda, `UserEntity` chiqadigan **barcha** javoblarni
   (ariza detali, review bosqichlari, co-author'lar, `my-application`) bir vaqtda yopadi.

2. `src/entrypoint/presentation/handlers/journal_config/schema/journal_config_response.go` —
   `CreatorResponse.Pin` maydoni va uni to'ldiruvchi kod (`:57-63`) butunlay olib
   tashlanadi. `CreatorResponse` faqat `ID` va `FullName` qoladi.

**Nega `json:"-"`, mapper'da `nil` emas:** teg deklarativ va bir joyda turadi;
qorovul test bilan mahkamlanadi. Mapper'da `Pin = nil` qilish 10+ joyda takrorlanadi
va drift beradi — SQL/upload ishlarida qochgan sinf.

**Reja bosqichida tasdiqlanadi:** `.Pin` va `.BirthDate` ning `src/` bo'ylab barcha
o'qilishlari JSON round-trip qilmasligini `go build` va qo'lda tekshiruv bilan
tasdiqlash (hozircha faqat schema va mapper topilgan — ikkalasi ham xavfsiz).

### 3.2. Telefon va email'ni so'rovchiga qarab redaksiya

Telefon uch entity'da chiqadi, himoya ehtiyoji farq qiladi:

| Entity | Endpoint | Himoya |
|---|---|---|
| `UserBasicEntity.PhoneNumber` | `/journal-manage/{id}/members` | Endpoint allaqachon a'zolikni tekshiradi (PR #18) — telefon qoladi |
| `ReviewerEntity.PhoneNumber` | `/journal/{id}/reviewers` | **Olib tashlanadi** (3.3-bo'lim) |
| `UserEntity.PhoneNumber` / `.Email` | ariza detallari | **Redaksiya kerak** — endpoint egalikni tekshirmaydi |

Faqat `UserEntity` (ariza detalida ichma-ich) telefon/email'i maydon darajasida
redaksiya qilinadi, chunki uni har qanday login qilgan foydalanuvchi ko'radi.

**Markazlashtirilgan yordamchi** (`src/core/application/usecase/userusecases` yoki
mos paketda):

```go
// RedactContactInfo so'rovchi egasi yoki privilegiyalangan bo'lmasa
// aloqa ma'lumotlarini bo'shatadi. PINFL va birth_date bu yerda kerak emas —
// ular json:"-" bilan baribir chiqmaydi.
func RedactContactInfo(user *entity.UserEntity, requesterID uint, isPrivileged bool) {
    if user == nil || user.ID == requesterID || isPrivileged {
        return
    }
    user.PhoneNumber = ""
    user.Email = ""
}
```

**`isPrivileged`:** admin (`permissionusecases.IsAdmin`). Jurnal muharriri tekshiruvi
ataylab kiritilmaydi — u workflow-egalik bilan bog'liq va qamrovdan tashqarida.
Natijada ariza detalida telefon/email faqat so'rovchining o'z yozuvida ko'rinadi
(masalan co-author sifatida o'zi) yoki admin uchun. Bu eng xavfsiz (deny-by-default)
tomon; muharrirlarga kerak bo'lsa keyin kengaytiriladi.

**Ta'siri:**

- `application_detail_usecase.Execute(applicationId)` →
  `Execute(requester *entity.UserBasicEntity, applicationId uint)` (IDOR ishidagi
  imzo o'zgarishi naqshi). Handler so'rovchini kontekstdan uzatadi.
- Use case javobdagi **har bir** `UserEntity` ustida yordamchini chaqiradi:
  `.User`, har bir `.Article.CoAuthors[].User`, har bir `.ReviewStages[].Reviewer`.
- `myapplication_detail_usecase` da so'rovchi = ega (`GetUserAppByID(applicationId, userId)`
  allaqachon egalikni tekshiradi), shuning uchun eganing o'z telefoni ko'rinadi;
  co-author/reviewer'lar redaksiya qilinadi.

**Aniq use case ro'yxati reja bosqichida tasdiqlanadi:** `ApplicationResponseEntity`
qaytaradigan ikki use case (`application_detail`, `myapplication_detail`) asosiy.
`UserEntity` qaytaradigan boshqa use case'lar (`user_list`, `user_find`, `user_update`,
`journal_members_list`, publisher/institution admin) sidiruv bilan tekshiriladi —
ularning ko'pchiligi admin-gated yoki o'z-ma'lumot, lekin har biri qat'iy ko'riladi.

### 3.3. Reviewer telefon ro'yxatini olib tashlash

`ReviewerEntity.PhoneNumber` → `json:"-"`. Telefon reviewer ro'yxatidan butunlay
chiqadi (foydalanuvchi qarori: login qilganlar ismlarni ko'radi, telefon hech kimga).

---

## 4. Yechim — B o'lcham

### 4.1. `/journal/{id}/reviewers` ni login orqasiga qo'yish

Hozir (`journal_group.go:48`):

```go
group.GET("/:journalId/reviewers", this.reviewersList.Handle)
```

`/journal` guruhi permission'siz (`app.go:293`), marshrut middleware'siz — anonim
kirish. O'zgarish: `permissions.AuthenticatedPermission` bilan o'raladi:

```go
group.GET("/:journalId/reviewers", this.reviewersList.Handle, permissions.AuthenticatedPermission)
```

Telefon 3.3-bo'limda allaqachon olib tashlangani uchun, bu ikki qatlam birga: login
shart + telefon yo'q.

---

## 5. Qorovul va testlar

### 5.1. Maxfiy JSON teg qorovuli

`test/architecture/sensitive_json_test.go` — `go/ast` bilan `src/core/domain/entity/`
va handler schema kataloglaridagi barcha struct teglarini skanerlaydi.

**Qoida:** hech qanday struct maydonida `json:"pin"`, `json:"pinfl"`, yoki
`json:"birth_date"` bo'lmasligi kerak. Kimdir kelajakda PINFL'ni JSON'ga qaytarsa
yoki yangi entity'da qo'shsa, test `fayl:qator` bilan yiqiladi.

**Nega faqat pin/birth_date, telefon emas:** telefon `UserBasicEntity` da ataylab
qoladi (jurnal a'zolari uchun). Global qorovul uni yolg'on ushlaydi. Shuning uchun
telefon maydon-darajasidagi redaksiya bilan, PINFL/birth_date esa qorovul + `json:"-"`
bilan himoyalanadi — ikki xil talab, ikki xil himoya.

### 5.2. Testlar

- `RedactContactInfo` — o'zi ko'radi / admin ko'radi / begona ko'rmaydi (telefon+email
  bo'sh) / `nil` user
- `UserEntity` JSON marshal — `Pin` va `BirthDate` chiqmaydi (teg tasdig'i, `encoding/json`
  bilan marshal qilib, natijada bu kalitlar yo'qligini tekshirish)
- `application_detail_usecase` — begona so'rovchi uchun muallif, co-author va reviewer
  telefon/email'i bo'sh; so'rovchi = ega bo'lganda o'z yozuvi ko'rinadi
- `ReviewerEntity` JSON marshal — `phone_number` yo'q
- Qorovulning o'z tishi — `sensitive_json_test` mantiqiga yolg'on maxfiy teg
  (`json:"pin"`) berib, ushlashini tasdiqlash

**Ma'lumotlar bazasi talab qiladigan repozitoriylar testlanmaydi** — loyihada baza
bilan test yo'q. Butun qaror use case va entity darajasida, ular to'liq testlanadi.

---

## 6. Qamrovdan tashqarida

- **`/article-application/detail` va `/journal-applications/{id}` egalik guard'i** —
  chuqur workflow-egalik mantiqi. PINFL/birth_date olib tashlangandan keyin maxfiy
  ma'lumot sizmaydi; qolgan IDOR (ariza metadatasi) alohida ish.
- **Jurnal muharriri "privileged" sifatida** — 3.2 da faqat admin. Muharrirlarga
  ariza detalida telefon kerak bo'lsa, keyingi ishda.
- **`UserBasicEntity` telefoni** (`/journal-manage/{id}/members`) — endpoint a'zolikni
  tekshiradi (PR #18), o'zgartirilmaydi.
- **Baza bilan integratsiya testlari** — alohida infratuzilma ishi.

---

## 7. Frontend uchun ta'sir

- **PINFL, birth_date** — barcha API javoblaridan yo'qoladi. Agar biror ekran ularni
  ishlatayotgan bo'lsa (kutilmaydi — profil entity'sida yo'q edi), u ekran buziladi.
- **Reviewer telefoni** (`/journal/{id}/reviewers`) — yo'qoladi; endpoint endi login
  talab qiladi (anonim kirish to'xtaydi).
- **Ariza detalidagi telefon/email** — begona so'rovchi uchun bo'sh keladi; ariza egasi
  o'z yozuvini ko'radi.

**Deploy oldidan:** frontend jamoasi bilan bu uchta o'zgarishni solishtirib chiqish —
ayniqsa reviewer endpointining anonim → login o'zgarishi va PINFL/birth_date'ning
yo'qolishi.

---

## 8. Shox va zanjir

```
develop ← PR #18  (feature/security-hardening)  CWE-613 + CWE-639
        ← PR #21  (feature/upload-hardening)    CWE-434 + CWE-79
        ← PR #22  (feature/sql-injection)       CWE-89
        ← yangi   (feature/sensitive-data)      CWE-200
```

Merge tartibi majburiy: #18 → #21 → #22 → yangi. `feature/sql-injection` ustiga
quriladi, chunki `test/` katalogi va qorovul naqshi shu zanjirda keladi.
