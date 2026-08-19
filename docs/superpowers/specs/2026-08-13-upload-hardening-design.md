# Fayl yuklash va saqlangan XSS himoyasi (CWE-434 + CWE-79)

**Sana:** 2026-08-13
**Holat:** Tasdiqlangan, amalga oshirishga tayyor
**Zaifliklar:** Ekspertiza hisoboti 2.2.2 — "Saytlararo skript (XSS)", xavflilik darajasi **Yuqori** (1 ta)
va 2.2.9 — "Ixtiyoriy kengaytmadagi faylni yuklash", xavflilik darajasi **O'rta** (5 ta)
**CWE:** 79 (Improper Neutralization of Input During Web Page Generation) va
434 (Unrestricted Upload of File with Dangerous Type)
**OWASP Top 2025:** A05 Injection · A06 Insecure Design

---

## 1. Muammo

Ekspertiza ikkita alohida topilma sifatida yozgan, lekin ular **bitta zanjirning bo'g'inlari**.
Hisobotning o'zi buni ko'rsatadi: XSS 2–4-rasmlarda aynan fayl yuklash oynasi orqali namoyish
qilingan, ixtiyoriy kengaytma esa 30–38-rasmlarda beshta yuklash oynasida takrorlangan.

Zanjir uch bo'g'indan iborat:

```
(1) Mijoz ixtiyoriy kengaytma va MIME turini tanlaydi
      → (2) MinIO faylni aynan o'sha tur bilan saqlaydi
            → (3) Brauzer uni ijro etadi
```

Kodbazani tekshirish har bir bo'g'inning sababini aniqladi:

| Bo'g'in | Joriy holat | Manba |
|---|---|---|
| Kengaytma tekshiruvi | **Umuman yo'q** | `usecase/uploadusecases/upload_file_usecase.go:31` |
| Bucket tanlash | **Mijoz so'rovda belgilaydi** | `handlers/upload/schema/presigned_put_url_schema.go:5` |
| Obyekt kaliti | Mijoz nomi va **papka yo'li** saqlanadi | `upload_file_usecase.go:33,44` |
| Content-Type | Presigned PUT'da cheklanmagan → mijoz belgilaydi | `storage/minio_storage.go:54` |
| Hajm chegarasi | Kod bor, lekin **izohga olingan** | `storage/minio_storage.go:74-76` |
| Magic-byte tekshiruvi | Faqat base64 yo'lida | `upload_base64_file_usecase.go:67` |
| Xizmat qilishda sarlavhalar | `reqParams` **bo'sh** | `storage/minio_storage.go:177` |
| HTTP xavfsizlik sarlavhalari | Birortasi ham yo'q | `app/app.go:242` |

### 1.1. Amaldagi ekspluatatsiya

`PresignedPutObject` faqat metod, bucket, kalit va muddatni imzolaydi — **Content-Type va hajmni
cheklay olmaydi**. Shuning uchun quyidagi ketma-ketlik hozir ishlaydi:

```
1. POST /api/upload/presigned/put-url   { "fileName": "evil.html", "bucket": "PUBLIC" }
   → server imzolaydi, chunki hech narsa tekshirilmaydi
2. PUT <presigned-url>                  Content-Type: text/html
   → MinIO aynan shu turni obyekt metadatasiga yozadi
3. GET https://storage.sciencelib.uz/slibbucketpublic/evil_a1b2.html
   → brauzer HTML ni IJRO ETADI
```

### 1.2. Ta'sir doirasi — halol baholash

Fayllar `storage.sciencelib.uz` subdomenidan xizmat qilinadi, ilova esa `journal.sciencelib.uz`
da joylashgan. Ular **turli origin** bo'lgani uchun yuklangan HTML ilova tokenlariga
(localStorage) yeta olmaydi. Ya'ni bu sessiya o'g'irlashga olib kelmaydi.

Qoladigan real ta'sir — ekspertiza aynan shuni yozgan: **davlat ishonchli domenida ixtiyoriy
HTML/JS joylashtirish**. Sertifikatlar va ilmiy maqolalar bilan bir domenda turgani uchun
fishing xabarlarida ishonch suiiste'mol qilinadi.

### 1.3. Ommaviy bucket policy'si — dizaynga hal qiluvchi ta'sir

Ommaviy bucket anonim o'qishga ochiq:

```json
{ "Effect": "Allow", "Principal": { "AWS": ["*"] },
  "Action": ["s3:GetObject"], "Resource": ["arn:aws:s3:::slibbucketpublic/*"] }
```

Bundan kelib chiqadigan xulosa dizayn uchun markaziy: **ommaviy bucket'dagi faylga presigned
URL'siz murojaat qilinadi**, demak presigned GET paytidagi har qanday sarlavha override'i
chetlab o'tiladi. Ommaviy bucket uchun yagona ishlaydigan himoya — metadatani **yuklash paytida
obyektning o'ziga yozish**.

---

## 2. Qamrov

**Kiradi:**
- Kengaytma va MIME allow-list'i — `src/core/domain` da yagona manba
- `PresignedPutObject` → `PresignedPostPolicy` (hajm, Content-Type, kalit MinIO tomonidan majburlanadi)
- Bucket tanlashni mijozdan olib qo'yish
- `POST /api/upload/file` uchun magic-byte tekshiruvi va hajm chegarasi
- Presigned GET javob sarlavhalari override'i (maxfiy bucket uchun)
- Echo'ga HTTP xavfsizlik sarlavhalari

**Kirmaydi** (alohida ishlar):
- **CORS origin whitelist** — `app.go:246` da `AllowOrigins` izohga olingan, `CORS_ALLOW_ORIGIN`
  env mavjud va `required` bo'lsa-da ulanmagan. Foydalanuvchi qarori bilan bu ishdan chiqarildi
- **Antivirus skanerlash** — ekspertiza tavsiya qilgan, lekin bu alohida infratuzilma ishi
- **Bucket'da allaqachon yotgan fayllar** — ommaviy bucket test muhitiniki, mazmuni ahamiyatsiz
- **Ekspertiza 2.2.8 (Clickjacking)** — u frontend sahifalariga tegishli, bu backend ularni
  xizmat qilmaydi. `X-Frame-Options` bu yerda uni yopmaydi
- **Yuklashdan keyin baytlarni MinIO'dan qayta o'qib tekshirish.** XSS ni yopish uchun kerak
  emas — Content-Type MinIO tomonidan majburlangani sababli brauzer HTML ni ijro etmaydi.
  Lekin u "PDF deb atalgan bajariluvchi fayl" ssenariysini yopardi; 10-bo'limga qarang

---

## 3. Qabul qilingan qarorlar

| Savol | Qaror | Sabab |
|---|---|---|
| Ruxsat etilgan turlar | PDF, JPEG, PNG, DOC, DOCX | Foydalanuvchi belgilagan mahsulot ehtiyoji |
| SVG | **Rad etiladi** | SVG ichida `<script>` ishlaydi — u XSS vektori |
| Excel yuklash | **Rad etiladi** | Faqat eksport uchun ishlatiladi, yuklash oqimi yo'q |
| Presign usuli | `PresignedPostPolicy` | Hajm + Content-Type + kalit — uchalasini MinIO majburlaydi |
| Frontend ta'siri | **Buzuvchi o'zgarish qabul qilindi** | To'liq majburlash muhimroq deb topildi |
| Allow-list joyi | Kodda konstanta, env'da emas | Xavfsizlik chegarasi konfiguratsiya bilan bo'shatilmasligi kerak |
| Rasm disposition | `inline` | Majburlangan `image/*` turi yetarli; `attachment` UX ni buzardi |
| PDF disposition | `inline` | **2026-08-18 da o'zgartirildi** — quyidagi tuzatishga qarang |
| Office hujjat disposition | `attachment` | Brauzer `.doc`/`.docx` ni render qila olmaydi — `inline` amalda yuklab olishga tushadi |
| Xatolik holati | **Rad etish** | Tur aniqlanmasa yoki shubha bo'lsa — yuklash rad etiladi |

---

## 4. Arxitektura

Himoya uch qatlamli. Muhimi: **har bir qatlam mustaqil ravishda hujumni to'xtata oladi**, va
qatlamlarning qiymati bucket turiga qarab farq qiladi.

| Qatlam | Qayerda | Maxfiy bucket | Ommaviy bucket |
|---|---|---|---|
| L1 — allow-list | URL berish paytida | Asosiy | Asosiy |
| L2 — obyekt metadatasi | Yuklash paytida (POST policy) | Qo'shimcha | **Yagona himoya** |
| L3 — javob override'i | GET paytida (presigned) | **Asosiy** | Ishlamaydi (anonim kirish) |

L1 eng qimmatli bo'g'in, chunki u **yuklash boshlanishidan oldin** ishlaydi: ruxsat etilmagan
kengaytma uchun presigned URL umuman berilmaydi.

### 4.1. Nega L3 ommaviy bucket uchun ishlamaydi

`response-content-disposition` — presigned URL'ning **so'rov parametri**. Anonim GET so'rovida
bunday parametr bo'lmaydi, MinIO obyektning saqlangan metadatasini qaytaradi. Shuning uchun
ommaviy bucket'da faqat L2 himoya qiladi va uni "qo'shimcha qatlam" deb hisoblash xato bo'lardi.

---

## 5. Komponentlar

### 5.1. Yangi: allow-list — `src/core/domain/entity/enum/upload_type.go`

Yagona haqiqat manbai. Uchta narsani bir joyda bog'laydi: kengaytma → MIME turi → magic-byte
imzosi. Uchtasi uch joyda yashasa, ular vaqt o'tishi bilan uzoqlashadi — IDOR ishida
`JournalMemberPermissionUseCase` bilan aynan shu muammo bo'lgan.

```go
type UploadType struct {
    Extension   string   // ".pdf" — nuqta bilan, kichik harfda
    ContentType string   // "application/pdf"
    Magic       [][]byte // ruxsat etilgan bosh baytlar; bo'sh bo'lsa tekshirilmaydi
    Disposition string   // "attachment" | "inline"
}

// LookupUploadType kengaytma bo'yicha turni qaytaradi.
// Kengaytma katta-kichik harfdan qat'i nazar topiladi (.PDF == .pdf).
// Topilmasa (nil, false) — chaqiruvchi RAD ETISHI shart.
func LookupUploadType(ext string) (*UploadType, bool)
```

Boshlang'ich ro'yxat:

| Kengaytma | Content-Type | Magic-byte | Disposition |
|---|---|---|---|
| `.pdf` | `application/pdf` | `%PDF-` | inline |
| `.jpg`, `.jpeg` | `image/jpeg` | `FF D8 FF` | inline |
| `.png` | `image/png` | `89 50 4E 47 0D 0A 1A 0A` | inline |
| `.doc` | `application/msword` | `D0 CF 11 E0 A1 B1 1A E1` | attachment |
| `.docx` | `application/vnd.openxmlformats-officedocument.wordprocessingml.document` | `50 4B 03 04` | attachment |

Diqqat: `.docx` ZIP konteyneri, ya'ni `PK\x03\x04` imzosi `.docx` ni boshqa ZIP asosidagi
formatlardan ajratmaydi. Bu ongli kelishuv — ZIP ichini ochib tekshirish bu spec doirasidan
tashqarida, va allow-list'ning o'zi asosiy himoyani beradi.

### 5.2. O'zgaradi: `UploadFileUseCase` — presigned yuklash

`src/core/application/usecase/uploadusecases/upload_file_usecase.go`

Hozirgi `generateUniqueObjectName` **olib tashlanadi**. Sababi `upload_file_usecase.go:33` dagi
`path.Dir(original)`: u mijozning papka yo'lini saqlaydi, ya'ni mijoz obyekt kalitining
tuzilishini boshqara oladi.

Yangi imzo bucket'ni ham, papkani ham serverga o'tkazadi:

```go
func (this *UploadFileUseCase) Execute(
    ctx context.Context,
    folder enum.StorageFolder,   // server belgilaydi, mijoz emas
    fileName string,             // faqat kengaytmani olish uchun ishlatiladi
) (*PresignedUpload, error)

type PresignedUpload struct {
    URL        string
    Fields     map[string]string  // POST policy form maydonlari
    ObjectName string
}
```

Mantiq:
1. `filepath.Ext(fileName)` → kichik harfga keltiriladi
2. `LookupUploadType(ext)` → topilmasa `response.InvalidFileError` (**400**), URL berilmaydi
3. Obyekt kaliti serverda quriladi: `{folder}/{sanitized}_{randomHex}{ext}` — mijoz yo'lining
   `path.Dir` qismi **tashlab yuboriladi**
4. Bucket `folder` dan kelib chiqadi (5.4-bo'limga qarang)
5. POST policy quriladi va imzolanadi

### 5.3. O'zgaradi: `FileStorage` porti va MinIO implementatsiyasi

`PutObjectPresignedUrl` o'rniga:

```go
// PostPolicyPresignedUrl yuklash uchun imzolangan POST policy qaytaradi.
// Content-Type, hajm oralig'i va obyekt kaliti policy ichida qat'iy belgilanadi —
// ularni mijoz o'zgartira olmaydi, MinIO mos kelmagan yuklashni rad etadi.
PostPolicyPresignedUrl(ctx context.Context, bucket enum.Bucket, objectName string,
    uploadType *enum.UploadType, maxSize int64) (*url.URL, map[string]string, error)
```

`minio-go v7.0.100` da mavjud (`post-policy.go`da tekshirildi):

| Chaqiruv | Nima beradi |
|---|---|
| `SetBucket`, `SetKey` | Kalitni server belgilaydi |
| `SetContentType` | Mijoz `text/html` yubora olmaydi — policy mos kelmaydi |
| `SetContentLengthRange(1, maxSize)` | Hajmni MinIO majburlaydi |
| `SetContentDisposition` | Obyekt metadatasiga yoziladi — **anonim GET'da ham qaytariladi** |
| `SetExpires` | 30 daqiqa (hozirgi qiymat saqlanadi) |

### 5.4. Bucket tanlash serverga o'tadi

`PresignedPutUrlRequest` dan `Bucket` maydoni **olib tashlanadi**. Uning o'rniga mijoz maqsadni
(`folder`) yuboradi, bucket esa papkadan kelib chiqadi:

| Papka | Bucket | Sabab |
|---|---|---|
| `news`, `images` | PUBLIC | Saytda ochiq ko'rsatiladi |
| `articles`, `applications`, `certificates`, `antiplag_certificate`, `expert_conclusion`, `tickets`, `spellcheck`, `article` | PRIVATE | Ilmiy va shaxsiy hujjatlar |

Noma'lum papka → **400**. Bu ro'yxat kodda konstanta bo'ladi.

### 5.5. O'zgaradi: `POST /api/upload/file` — server orqali yuklash

`storage/minio_storage.go:63` `SaveTempFile`. Bu yerda baytlar bizda, shuning uchun tekshiruv
kuchliroq bo'ladi:

1. Kengaytma allow-list'da bo'lishi shart
2. **Magic-byte** e'lon qilingan turga mos kelishi shart — birinchi baytlar o'qiladi va
   `UploadType.Magic` bilan solishtiriladi
3. Hajm chegarasi — `minio_storage.go:74-76` dagi izohga olingan kod jonlantiriladi

2-qadam `upload_base64_file_usecase.go:67` dagi `%PDF-` mantig'ini umumlashtiradi. O'sha joy ham
yangi umumiy funksiyaga o'tkaziladi, ya'ni PDF tekshiruvi ikki nusxada qolmaydi.

Base64 yo'lidagi `PutObject` chaqiruvi (`minio_storage.go:167`) hozir Content-Type ni
`utils.GetContentType` dan oladi. U ham allow-list'ga o'tkaziladi — aks holda bitta yuklash yo'li
`mime.TypeByExtension` ga, qolganlari allow-list'ga tayangan bo'lib qolardi.

### 5.6. O'zgaradi: presigned GET — `minio_storage.go:177`

`reqParams` hozir bo'sh. Ikkita parametr qo'shiladi:

```go
reqParams.Set("response-content-disposition", disposition)  // attachment; filename="..."
reqParams.Set("response-content-type", contentType)         // saqlangan turni bekor qiladi
```

Ikkalasi ham S3 standarti va `minio-go` tomonidan ruxsat etilgan (`utils.go:552,555`).

Disposition va Content-Type obyekt kalitining kengaytmasidan `LookupUploadType` orqali olinadi.
Kengaytma allow-list'da bo'lmasa (eski fayllar) — `application/octet-stream` va `attachment`,
ya'ni **noma'lum tur eng xavfsiz holatga tushadi**.

### 5.7. Yangi: HTTP xavfsizlik sarlavhalari — `app/app.go:242`

```go
this.echo.Use(middleware.SecureWithConfig(middleware.SecureConfig{
    XSSProtection:         "0",
    ContentTypeNosniff:    "nosniff",
    XFrameOptions:         "DENY",
    ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
}))
```

Qiymati haqida halol baho: bu backend JSON API qaytaradi, sahifa emas. `nosniff` real foyda
beradi — API binar javoblar ham qaytaradi (statistika Excel eksporti
`handlers/journal/journal_statistics_handler.go:56`, base64 yuklab olish). CSP va
`X-Frame-Options` bu yerda kichik qiymatga ega, lekin arzon va zarari yo'q.

`XSSProtection: "0"` ataylab: eski `X-XSS-Protection` sarlavhasi zamonaviy brauzerlarda
o'chirilgan va yoqilganda o'zi zaiflik manbai bo'lishi mumkin.

### 5.8. O'zgarmaydi

`SanitizeFilename` (`utils/filename.go:13`) hozirgi holida qoladi — u kirill harflarini ruxsat
etadi, bu maqsadli. `GetContentType` saqlanadi, lekin yuklash yo'llarida endi ishlatilmaydi:
uning `mime.TypeByExtension` ga tayanishi allow-list'dan kengroq va shuning uchun xavfsizlik
qarori uchun yaroqsiz.

---

## 6. Oqimlar

### 6.1. Presigned yuklash (yangi)

```
1. POST /api/upload/presigned/put-url   { "fileName": "maqola.pdf", "folder": "articles" }
2. Kengaytma ajratiladi → ".pdf" → LookupUploadType → topildi
   ├─ topilmasa → 400, TAMOM (yuklash boshlanmaydi)
3. Papka → bucket (articles → PRIVATE)
4. Obyekt kaliti: articles/maqola_a1b2c3d4.pdf
5. POST policy: ContentType=application/pdf, Length=[1,max],
                Disposition=inline, Key=<yuqoridagi>
6. Javob: { url, fields{...}, objectName }
7. Mijoz: POST <url> multipart/form-data (avval fields, oxirida file)
8. MinIO policy'ga solishtiradi → mos kelmasa RAD ETADI
```

### 6.2. Xizmat qilish

```
Maxfiy bucket:  presigned GET + response-content-disposition/type override
Ommaviy bucket: anonim GET → obyektning saqlangan metadatasi qaytariladi
                (5-qadamda yozilgani)
```

---

## 7. API kontrakti o'zgarishi — frontend bilan muvofiqlashtirish

Bu **buzuvchi o'zgarish**. Backend frontend'dan oldin chiqsa, yuklash ishlamay qoladi.

**Hozir:**
```
POST /api/upload/presigned/put-url   { "fileName": "a.pdf", "bucket": "PRIVATE" }
→ { "url": "...", "objectName": "..." }
→ mijoz: PUT <url>, body = fayl baytlari
```

**Keyin:**
```
POST /api/upload/presigned/put-url   { "fileName": "a.pdf", "folder": "articles" }
→ { "url": "...", "fields": { ... }, "objectName": "..." }
→ mijoz: POST <url>, multipart/form-data
         AVVAL barcha fields maydonlari, ENG OXIRIDA "file" maydoni
```

Maydonlar tartibi S3 talabi — `file` oxirgi bo'lishi shart, aks holda MinIO policy'ni
qo'llay olmaydi.

Yangi xatolar:

| Holat | Status | Xato |
|---|---|---|
| Kengaytma allow-list'da yo'q | 400 | `InvalidFileError` (`response.go:42`) |
| Noma'lum papka | 400 | `InvalidFileError` |
| Magic-byte mos emas (`/upload/file`) | 400 | `InvalidFileError` |
| Hajm chegaradan katta (`/upload/file`) | 413 | `EntityToLargeError` (`response.go:43`) |
| Hajm chegaradan katta (presigned) | MinIO'dan `EntityTooLarge` | Mijoz ushlashi kerak |
| Content-Type policy'ga mos emas | MinIO'dan `AccessDenied` | Mijoz noto'g'ri yubordi |

---

## 8. Testlash

Uslub oldingi ikki rejadagidek: yangi bog'liqlik yo'q, standart `testing`, qo'lda yozilgan
fake'lar, testlar loyiha rootidagi `test/` katalogida paket tuzilishini takrorlaydi
(`package <dir>_test`).

**Allow-list**
- Har bir ruxsat etilgan kengaytma topiladi va to'g'ri Content-Type qaytaradi
- Katta-kichik harf: `.PDF`, `.pdf`, `.pDf` — bir xil natija
- Ruxsat etilmagan: `.html`, `.svg`, `.exe`, `.php`, `.js` — `(nil, false)`
- Kengaytmasiz nom → `(nil, false)`

**Presigned yuklash**
- `evil.html` → 400, storage **umuman chaqirilmaydi** (chaqiruv hisoblagichi bilan)
- `maqola.pdf` → policy'da ContentType, hajm oralig'i va kalit to'g'ri o'rnatilgan
- Mijoz `../../etc/passwd` yuborsa → kalit `{folder}/passwd_xxxx` bo'ladi, `..` yo'qoladi
- Papka → bucket xaritasi: har bir papka to'g'ri bucket'ga tushadi
- Noma'lum papka → 400

**Magic-byte**
- PDF baytlari `.pdf` sifatida → qabul
- HTML baytlari `.pdf` nomi bilan → **rad**
- JPEG/PNG imzolari to'g'ri tekshiriladi
- Fayl imzodan qisqa (masalan 2 bayt) → rad, panic emas

**Presigned GET**
- `reqParams` da ikkala parametr o'rnatiladi
- Noma'lum kengaytmali eski obyekt → `octet-stream` + `attachment`

---

## 9. Deploy tartibi

1. **Oldindan tekshirish:** ishlab chiqarishdagi ommaviy bucket policy'si haqiqatan anonim
   `s3:GetObject` ga ochiqligini tasdiqlash. Agar yopiq bo'lsa — dizayn faqat kuchayadi, lekin
   buni taxmin qilib qo'ymaslik kerak
2. `UPLOADED_FILE_MAX_SIZE` qiymatini kelishish. **Yangi env o'zgaruvchisi kerak emas** —
   maydon `env.go:51` da allaqachon mavjud va `required`, ya'ni qiymat to'ldirilgan; faqat
   `minio_storage.go:74` da izohga olingan tekshiruv jonlantiriladi
3. **Frontend bilan birga chiqarish** — 7-bo'limga qarang. Backend oldin chiqsa yuklash siniydi
4. Chiqargandan keyin har bir yuklash oqimini qo'lda tekshirish: maqola PDF, jurnal muqovasi,
   yangilik rasmi, ekspert xulosasi, avatar

---

## 10. Qoldiq risklar

**Baytlar presigned yo'lda tekshirilmaydi.** MinIO Content-Type'ni majburlaydi, lekin baytlarni
o'qimaydi. Mijoz HTML baytlarini `application/pdf` yorlig'i bilan yuklay oladi. XSS uchun bu
muhim emas — saqlangan tur `application/pdf` bo'lgani uchun brauzer HTML'ni ijro etmaydi. Lekin
"ishonchli domendan zararli fayl tarqatish" ssenariysining bir qismi ochiq qoladi: PDF deb
atalgan bajariluvchi fayl yuklanishi mumkin. To'liq yopish uchun yuklashdan keyin baytlarni
qayta o'qish kerak — ongli ravishda keyingi ishga qoldirildi.

**`.docx` imzosi ZIP bilan bir xil.** `PK\x03\x04` `.docx` ni boshqa ZIP formatlaridan
ajratmaydi. Konteyner ichini ochib tekshirish bu spec doirasidan tashqarida.

**CORS hamon barcha domenlarga ochiq.** `app.go:246` da `AllowOrigins` izohga olingan.
Foydalanuvchi qarori bilan bu ishdan chiqarildi, lekin u yuqoridagi himoyalarni zaiflashtiradi
va alohida ish sifatida qolmoqda.

**Antivirus yo'q.** Ekspertiza tavsiya qilgan, lekin bu infratuzilma ishi.

**Ommaviy bucket'da eski fayllar.** Test muhitiniki bo'lgani uchun hisobga olinmadi. Agar
ishlab chiqarish bucket'ida ham xavfsiz bo'lmagan metadatali fayllar bo'lsa, ular L2 dan
foyda ko'rmaydi — inventarizatsiya alohida ish bo'ladi.

---

## 11. Muvaffaqiyat mezonlari

1. `evil.html` uchun presigned URL **berilmaydi** — 400, storage chaqirilmaydi
2. Ruxsat etilgan kengaytma bilan olingan URL'ga boshqa Content-Type bilan yuklashga urinish
   MinIO tomonidan **rad etiladi**
3. Hajm chegarasidan katta fayl MinIO tomonidan **rad etiladi**
4. Obyekt kaliti mijoz yuborgan papka yo'lini **saqlamaydi**
5. `POST /api/upload/file` ga PDF nomi bilan HTML baytlari yuborilsa **rad etiladi**
6. Ommaviy bucket'ga yuklangan `.doc`/`.docx` anonim GET'da
   `Content-Disposition: attachment` bilan qaytariladi; `.pdf` — `inline`
7. Maxfiy bucket'dagi noma'lum kengaytmali obyekt presigned GET'da `attachment` bilan qaytariladi
8. Barcha API javoblarida `X-Content-Type-Options: nosniff` mavjud
9. Mavjud yuklash oqimlari (maqola, muqova, yangilik, xulosa, avatar) ishlashda davom etadi

---

## 12. Tuzatish — 2026-08-18: PDF `inline`

**Muammo.** Chiqarilgandan keyin PDF'lar brauzerda umuman ochilmay qoldi: har bir havola
darhol yuklab olishga o'tardi, PDF viewer kengaytmalari ham ishlamasdi. Sabab —
`Content-Disposition: attachment` top-level navigatsiyada brauzer uchun so'zsiz buyruq:
u resursni ichki viewer'ga ham, kengaytmaga ham bermaydi.

**Qaror.** `.pdf` uchun disposition `inline` ga o'zgartirildi
(`enum/upload_type.go`). `.doc`/`.docx` va noma'lum kengaytmalar `attachment` bo'lib qoladi.

**Xavfsizlikka ta'siri.** Asosiy himoya o'zgarmadi: MinIO `Content-Type` ni majburlaydi,
ya'ni `.pdf` sifatida yuklangan HTML baytlari `application/pdf` yorlig'i bilan qaytariladi
va brauzer uni HTML sifatida ijro etmaydi. PDF'ning o'zi brauzer viewer'ida sandbox'da
ochiladi. 10-bo'limdagi "PDF deb atalgan bajariluvchi fayl" qoldiq riski o'zgarishsiz —
u `attachment` bilan ham to'liq yopilmagan edi.

**Migratsiya qilinmadi.** Ommaviy bucket'da allaqachon yotgan PDF'lar metadatasida
`attachment` saqlanib qoladi va yuklab olishga tushaveradi. Foydalanuvchi qarori bilan
bu ish doiradan chiqarildi (10-bo'limdagi "Ommaviy bucket'da eski fayllar" bandiga mos). Agar keyinchalik
kerak bo'lsa — `CopyObject` + `REPLACE` bilan metadata qayta yoziladi.
