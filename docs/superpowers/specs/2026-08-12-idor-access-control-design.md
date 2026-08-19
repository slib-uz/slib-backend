# Resursga bog'langan ruxsat nazorati (IDOR, CWE-639)

**Sana:** 2026-08-12
**Holat:** Tasdiqlangan, amalga oshirishga tayyor
**Zaiflik:** Ekspertiza hisoboti 2.2.3 — "Xavfsiz bo'lmagan to'g'ridan-to'g'ri obyekt manbalari", xavflilik **Yuqori**, 2 ta holat
**CWE:** 639 (Authorization Bypass Through User-Controlled Key) · **OWASP Top 2025:** A01 Broken Access Control

---

## 1. Muammo

Ekspertiza ikkita endpointda IDOR aniqladi: past huquqli foydalanuvchi so'rov parametrini
o'zgartirib begona ma'lumotga kirgan.

| Endpoint | Nima ochilgan |
|---|---|
| `GET /api/user/{id}` | Istalgan foydalanuvchi profili: F.I.SH., science ID, ilmiy daraja/unvon, ORCID, ish joyi va **roles massivi** (qaysi jurnal/nashriyotda qanday lavozim) |
| `GET /api/journal-manage/{journalId}/members` | Istalgan jurnalning tahririyat a'zolari |

> **Tuzatish (amalga oshirishdan keyin).** Spec'ning dastlabki tahriri `GET /api/user/{id}`
> javobida PINFL va pasport ma'lumotlari bor deb yozgan edi — bu noto'g'ri. `UserDetailEntity`
> (`src/core/domain/entity/user_detail_entity.go`, so'rov `user_repository_impl.go:525-565`)
> da na PINFL, na pasport, na tug'ilgan sana, na email, na telefon bor. Cheklov baribir
> o'rinli: roles massivi tashkiliy tuzilmani ochadi. Haqiqiy bog'lanish ma'lumotlari
> sizishi boshqa yo'nalishda edi — `GET /api/user/profile/{id}` (email va telefon).
> U kirish nazorati bilan emas, javobdan maydonlarni olib tashlash bilan yopildi
> (9-bo'limga qarang), chunki bu endpoint ommaviy muallif kartochkasini to'ldiradi.

Ikkala usecase resurs ID'sidan boshqa hech narsa qabul qilmaydi — ya'ni tekshirish uchun
ma'lumotning o'zi ham yo'q:

```go
// userusecases/user_detail_usecase.go:17
func (this *UserDetailUseCase) Execute(id uint) (*entity2.UserDetailEntity, error) {
	return this.repository.GetDetailByID(id)          // kim so'rayotgani noma'lum
}

// journalusecases/journal_members_list_usecase.go:19
func (this *JournalMembersListUsecase) Execute(journalID uint) (...) {
	_members, err := this.repository.GetByJournalID(journalID)
}
```

### 1.1. Kodbazani tekshirishda topilgan kengroq muammo

Loyihada `permissionusecases/` paketida 9 ta ruxsat tekshiruvi bor va aksariyati scope'ni
to'g'ri tekshiradi. Lekin eng ko'p ishlatiladigani —
`JournalMemberPermissionUseCase` — ichida `TODO` bilan birga xato yotibdi:

```go
// permissionusecases/journal_member_permission_usecase.go
// TODO: Check if the user has a role in the journal and if that role is one of the allowed roles for the action.
if role.Role == enum.RolePublisherAdmin {
	return true                    // QAYSI nashriyot ekani tekshirilmaydi
}
```

`RolePublisherAdmin` roliga ega istalgan foydalanuvchi **istalgan jurnalning** boshqaruv
ma'lumotlariga kiradi. Bu funksiya **20 ta joyda** chaqiriladi, jumladan taqrizlash, ariza
ro'yxatlari, fayl yuklab olish havolalari va tahririyat CRUD amallari.

Ya'ni ekspertiza 2 ta IDOR topgan, ammo bu bitta `TODO` yana 20 ta chaqiruv joyini ochiq
qoldirgan.

---

## 2. Qamrov

**Kiradi:**
- `GET /api/user/{id}` — o'zi yoki admin tekshiruvi
- `GET /api/journal-manage/{journalId}/members` — jurnal boshqaruvi tekshiruvi
- `JournalMemberPermissionUseCase` dagi publisher-admin xatosini tuzatish
- Yagona `IsAdmin` yordamchisi (faqat yangi kod uchun)

**Kirmaydi** (alohida ishlar):
- Deny-by-default avtorizatsiya: 302 route'dan 246 tasida guard yo'q
- `permissions.RolePermission` interceptorining scope'siz ekani
- PINFL/pasport ma'lumotlarini javobdan olib tashlash (ekspertiza 2.2.5)
- Mavjud uchta "admin" tekshiruvini yagona `IsAdmin` ga ko'chirish

---

## 3. Qabul qilingan qarorlar

| Savol | Qaror | Sabab |
|---|---|---|
| Qamrov | 2 ta endpoint + buzuq markaziy tekshiruv | Bitta funksiya tuzatilishi 20 ta joyni yopadi — eng yaxshi ta'sir/xarajat nisbati |
| `/user/{id}` siyosati | **Faqat o'zi va admin** (4.4-bo'limdagi `IsAdmin`: `IsAdmin` bayrog'i yoki `RoleAdmin` roli) | Ommaviy muallif ma'lumoti uchun `/api/author/...` va `/api/users/find` allaqachon mavjud |
| Publisher-admin semantikasi | Repozitoriy orqali to'g'ri tekshirish | Qisqa yo'lni olib tashlash qonuniy ish oqimini buzardi |
| `/journal-manage/.../members` | `JournalMemberPermissionUseCase` | Endpoint boshqaruv bo'limida; yangi mantiq yozilmaydi |
| Tekshiruv qatlami | **Usecase qatlami** | Loyihada o'rnatilgan uslub; repozitoriyli tekshiruv tabiiy joylashadi |
| Xato semantikasi | **Fail-closed** | Ruxsat tekshiruvi hech qachon xatoni ruxsatga aylantirmaydi |

---

## 4. Komponentlar

### 4.1. `JournalMemberPermissionUseCase` — markaziy tuzatish

`src/core/application/usecase/permissionusecases/journal_member_permission_usecase.go`

```go
type JournalMemberPermissionUseCase struct {
	journalRepository repository.JournalRepository      // yangi bog'liqlik
}

// Execute foydalanuvchining jurnal boshqaruviga kirish huquqini tekshiradi.
func (this *JournalMemberPermissionUseCase) Execute(
	userRoles []*entity.UserRoleEntity, journalID uint,
) (bool, error)                                        // imzo o'zgardi
```

Qoidalar:

| Rol | Shart |
|---|---|
| `RoleAdmin` | Har doim ruxsat |
| `RoleChiefEditor`, `RoleSecretary` | Faqat `role.JournalID != nil && *role.JournalID == journalID` |
| `RolePublisherAdmin` | `role.PublisherID != nil` va u jurnalning egasiga teng bo'lsa |
| Qolganlari | Rad |

Barcha ko'rsatkichli maydonlar (`JournalID`, `PublisherID`) ishlatilishdan oldin `nil` ga
tekshiriladi. Bu shunchaki ehtiyot chorasi emas: `permissions/role_based_permission.go:53`
da aynan shu tekshiruv yo'qligi sababli `*item.PublisherID` panic berishi mumkin — o'sha
fayl bu spec doirasiga kirmaydi, ammo bu yerda xato takrorlanmasligi kerak.

Publisher-admin uchun mavjud `JournalRepository.GetPublisherIdByJournalId(journalID)` ishlatiladi.

**Dangasa qidiruv:** DB so'rovi faqat foydalanuvchida `RolePublisherAdmin` roli bo'lsa **va**
arzon tekshiruvlar natija bermagan bo'lsagina bajariladi. Bosh muharrir yoki kotib uchun
qo'shimcha so'rov hech qachon ketmaydi.

### 4.2. `UserDetailUseCase`

`src/core/application/usecase/userusecases/user_detail_usecase.go`

```go
func (this *UserDetailUseCase) Execute(
	user *entity.UserBasicEntity, id uint,
) (*entity2.UserDetailEntity, error) {
	if user.ID != id && !permissionusecases.IsAdmin(user) {
		return nil, response.PermissionDeniedError
	}
	return this.repository.GetDetailByID(id)
}
```

### 4.3. `JournalMembersListUsecase`

`src/core/application/usecase/journalusecases/journal_members_list_usecase.go`

`JournalMemberPermissionUseCase` in'ektsiya qilinadi va chaqiriladi. Yangi mantiq yozilmaydi.

```go
func (this *JournalMembersListUsecase) Execute(
	user *entity.UserBasicEntity, journalID uint,
) ([]*entity.UserRoleWithBasicUserEntity, error)
```

### 4.4. Yangi `IsAdmin` yordamchisi

`src/core/application/usecase/permissionusecases/is_admin.go`

Loyihada hozir uch xil "admin" tushunchasi bor va ular kelishmaydi:

| Joy | Tekshiradi |
|---|---|
| `permissions.AdminPermission` | Faqat `c.User.IsAdmin` bayrog'i |
| `UserListUseCase` | `IsAdmin` **yoki** `RoleAdmin` |
| `JournalMemberPermissionUseCase` | Faqat `RoleAdmin` |

```go
func IsAdmin(user *entity.UserBasicEntity) bool {
	if user == nil {
		return false
	}
	if user.IsAdmin {
		return true
	}
	for _, r := range user.Roles {
		if r.Role == enum.RoleAdmin {
			return true
		}
	}
	return false
}
```

Yangi kod shuni ishlatadi. **Mavjud uchta joy ko'chirilmaydi** — bu ongli chala qoldirish:
ko'chirish xatti-harakatni kutilmagan joylarda o'zgartirishi mumkin va bu spec doirasidan
chiqadi.

### 4.5. Handlerlar

`user_detail_handler.go` va `journal_members_list_handler.go` `c.User` ni usecase'ga uzatadi.
Ikkalasi ham `AuthenticatedPermission` bilan himoyalangan guruhda, shuning uchun `c.User`
`nil` bo'lmaydi.

---

## 5. Oqimlar

### 5.1. Ruxsat tekshiruvi

```
1. Rollar bo'ylab: RoleAdmin bormi?                        → (true, nil)
2. Rollar bo'ylab: ChiefEditor/Secretary va jurnal mosmi?  → (true, nil)
3. Foydalanuvchida RolePublisherAdmin bormi?
   yo'q → (false, nil)                                      ← DB so'rovi yo'q
   ha   → GetPublisherIdByJournalId(journalID)
          xato → (false, err)
          mos  → (true, nil)
4.                                                          → (false, nil)
```

### 5.2. `GET /api/user/{id}`

`user.ID == id` yoki `IsAdmin(user)` → ruxsat; aks holda `403`.

### 5.3. `GET /api/journal-manage/{journalId}/members`

Guruh darajasidagi `AuthenticatedPermission` dan keyin usecase ichida
`JournalMemberPermissionUseCase` chaqiriladi.

---

## 6. Imzo o'zgarishining to'lqini — 20 ta chaqiruv joyi

`Execute` `bool` o'rniga `(bool, error)` qaytaradi. **16 tasi** odatiy shaklda:

```go
// hozir
if !this.memberPermission.Execute(user.Roles, edition.JournalID) {
	return response.NewFailResponse(403, "...")
}

// bo'ladi
allowed, err := this.memberPermission.Execute(user.Roles, edition.JournalID)
if err != nil {
	return err
}
if !allowed {
	return response.NewFailResponse(403, "...")
}
```

Odatiy shakldagi joylar:

1. `editionusecases/edition_create_usecase.go:21`
2. `editionusecases/edition_update_usecase.go:28`
3. `editionusecases/edition_delete_usecase.go:37`
4. `editionusecases/edition_attach_articles_usecase.go:28`
5. `editionusecases/edition_detach_articles_usecase.go:34`
6. `journaleditorialusecases/journal_editorial_create_usecase.go:21`
7. `journaleditorialusecases/journal_editorial_update_usecase.go:25`
8. `journaleditorialusecases/journal_editorial_delete_usecase.go:25`
9. `journalnewsusecases/journal_news_create_usecase.go:21`
10. `journalnewsusecases/journal_news_update_usecase.go:25`
11. `journalnewsusecases/journal_news_delete_usecase.go:25`
12. `article_applications_usecases/peer_review_usecase.go:45`
13. `article_applications_usecases/technical_reivew_usecase.go:43`
14. `article_applications_usecases/application_publish_usecase.go:51`
15. `articleusecases/direct_article_create_usecase.go:49`
16. `roiusecase/article_publish_roi_usecase.go:42`

**4 tasi nostandart shaklda** va alohida e'tibor talab qiladi:

| Joy | Hozirgi shakl | Nima qilish kerak |
|---|---|---|
| `uploadusecases/presigned_url_usecase.go:55` | `if Execute(...) { return true, nil }` | Xatoni yutib yubormasdan ijobiy tarmoqni saqlash |
| `article_applications_usecases/applications_list_usecase.go:62` | `if Execute(...) { return nil }` | Xuddi shunday |
| `deadlineusecases/extend_deadline_usecase.go:67` | `if x == y && Execute(...)` — qo'shma shart | Shartni ikkiga bo'lish; `&&` qisqa tutashuvi xatoni chalkashtiradi |
| `permissionusecases/application_reviewer_permission.go:23` | `return Execute(...), nil` | `return Execute(...)` ga soddalashadi |

---

## 7. Xatolarni qayta ishlash

| Holat | Natija |
|---|---|
| Ruxsat yo'q | `403 permission denied` |
| Jurnal topilmadi (publisher qidiruvida) | Xato yuqoriga uzatiladi |
| Repozitoriy xatosi | Xato yuqoriga uzatiladi |

Ruxsat tekshiruvi **hech qachon xatoni `true` ga aylantirmaydi** (fail-closed). Bu sessiya
bekor qilish spec'idagi fail-open'ning ataylab teskarisi: u yerda savol "sayt ishlab
tursinmi", bu yerda esa "begona ma'lumot ochilib ketsinmi".

---

## 8. Testlar

Sessiya spec'i bilan bir xil uslub: standart `testing` paketi, embedded interfeys bilan
soxta repozitoriy, yangi bog'liqlik qo'shilmaydi.

| Test | Nimani tekshiradi |
|---|---|
| `JournalMemberPermission` | Bosh muharrir faqat o'z jurnaliga kiradi |
| `JournalMemberPermission` | **Boshqa nashriyot admini rad etiladi** — hozir aynan shu buzuq |
| `JournalMemberPermission` | O'z nashriyoti jurnaliga publisher admin kiradi |
| `JournalMemberPermission` | `RoleAdmin` har doim kiradi |
| `JournalMemberPermission` | Publisher-admin roli yo'q bo'lsa DB so'rovi bajarilmaydi |
| `JournalMemberPermission` | Repozitoriy xatosi `(false, err)` bo'lib qaytadi |
| `UserDetailUseCase` | O'zini ko'radi; begonani ko'rmaydi; admin hammasini ko'radi |
| `JournalMembersListUsecase` | Begona jurnal a'zolari `403` bilan rad etiladi |

---

## 9. Bog'liqlik va qoldiq risklar

**Sessiya spec'i bilan bog'liqlik.** `docs/superpowers/specs/2026-08-12-session-revocation-design.md`
loyihadagi birinchi testlarni va `Makefile` ga `test` target'ini kiritadi. Ikkala ish
mustaqil; qaysi biri avval chiqsa, o'sha target'ni qo'shadi. Reja buni shartli qadam
sifatida yozadi.

**Yangi endpointda tekshiruv unutilishi mumkin.** Usecase qatlamidagi tekshiruv deklarativ
emas — yangi endpoint yozgan dasturchi uni qo'shishni unutishi mumkin. Buni interceptor ham
hal qilmaydi (u ham qo'shilishi kerak edi). Haqiqiy yechim — deny-by-default ishi, u alohida
spec sifatida qoladi.

**Uch xil "admin" tushunchasi saqlanib qoladi.** Yangi kod `IsAdmin` ni ishlatadi, lekin
mavjud uchta joy o'z mantig'ida qoladi. Bu ongli qaror, 4.4-bo'limga qarang.

**~~PINFL javobda qoladi.~~ (Noto'g'ri edi — tuzatildi.)** `/user/{id}` javobida
(`UserDetailEntity`) PINFL ham, pasport ham, tug'ilgan sana ham, email ham, telefon ham
yo'q; u {id, F.I.SH., science ID, rasm, roles, maqolalar soni} va ish joyi, ilmiy
ko'rsatkichlar, ilmiy daraja/unvon, ORCID dan iborat. Endpoint faqat o'ziga va adminga
cheklandi va bu qaror kuchida qoladi, lekin sababi boshqa: roles massivi tashkiliy
tuzilmani ochadi.

**Bog'lanish ma'lumotlari `/user/profile/{id}` da edi.** Ekspertiza ko'zda tutgan
"begona foydalanuvchi ma'lumotlarini o'qish" aslida shu yo'nalishda ochiq turgan edi:
`UserProfileEntity` email va telefonni o'z ichiga oladi. Endpoint kirish nazorati bilan
yopilmadi — u ommaviy muallif kartochkasini to'ldiradi, o'ziga/adminga cheklash boshqa
mualliflarni ko'rishni buzardi. Buning o'rniga so'rovchi egasi yoki admin bo'lmasa,
javobdan `Email` va `Phone` olib tashlanadi. Javobning qolgan qismi o'zgarmaydi.
Shaxsiy ma'lumotlarni javoblardan tizimli qisqartirish — ekspertiza 2.2.5 zaifligi,
alohida ish bo'lib qoladi.

---

## 10. Muvaffaqiyat mezonlari

1. Foydalanuvchi `GET /api/user/{begona_id}` so'rasa `403` oladi
2. Foydalanuvchi `GET /api/user/{o'z_id}` so'rasa `200` oladi
3. `RoleAdmin` istalgan foydalanuvchi profilini ko'ra oladi
4. A jurnalning bosh muharriri `GET /api/journal-manage/{B}/members` so'rasa `403` oladi
5. **A nashriyoti admini B nashriyoti jurnaliga tegishli amallarda `403` oladi** — hozir `200` oladi
6. O'z nashriyoti jurnallarida publisher admin ilgarigidek ishlaydi
7. `go build ./...` va `go test ./...` toza o'tadi
