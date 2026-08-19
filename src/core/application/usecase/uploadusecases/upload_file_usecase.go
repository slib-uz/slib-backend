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
	enum.FolderNews:   enum.BucketPublic,
	enum.FolderImages: enum.BucketPublic,
	// Nashr fayli — chop etilgan jurnal soni. U saytda ruxsatsiz ochiladi,
	// shuning uchun PUBLIC. Ilgari mijoz uni `applications` ga yozardi va
	// fayl PRIVATE bucket'ga tushib ko'rinmay qolardi.
	enum.FolderEditions: enum.BucketPublic,
	// Ekspert xulosasi PUBLIC: integration_article_create_usecase.go shu
	// papkaga allaqachon PUBLIC yozadi. Ikki yozuv yo'li bitta papkaga turli
	// bucket bilan yozsa, faylni o'qish paytida qaysi bucket ekani noma'lum
	// bo'lib qoladi.
	enum.FolderArticleExpertConclusion: enum.BucketPublic,
	// Jurnal va OAK sertifikatlari jurnal sahifasida ruxsatsiz ko'rsatiladi.
	enum.FolderCertificates: enum.BucketPublic,
	// Bu ikkalasini server o'zi yozadi va PUBLIC bucket'ga qo'yadi
	// (antiplag_status_update_usecase.go, spellcheck_process_usecase.go).
	// Xarita ularga zid bo'lsa, mijoz yuklagan fayl boshqa bucket'ga tushadi.
	enum.FolderAntiPlagCertificate: enum.BucketPublic,
	enum.FolderSpellCheck:          enum.BucketPublic,
	enum.FolderArticles:            enum.BucketPrivate,
	enum.FolderArticle:             enum.BucketPrivate,
	enum.FolderApplications:        enum.BucketPrivate,
	// Murojaat ilovalarida shaxsiy ma'lumot bo'lishi mumkin. Bu papkaga
	// hozircha hech kim yozmaydi; funksiya qurilganda qayta ko'riladi.
	enum.FolderTickets: enum.BucketPrivate,
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
