package storage

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"
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
