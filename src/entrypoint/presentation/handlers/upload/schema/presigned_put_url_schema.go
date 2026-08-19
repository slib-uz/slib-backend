package schema

import "slib.uz/src/core/domain/entity/enum"

// PresignedPutUrlRequest — yuklash so'rovi.
// Bucket maydoni ATAYLAB olib tashlandi: uni endi server papkadan aniqlaydi,
// aks holda mijoz maxfiy hujjatni ommaviy bucket'ga yuklay olardi.
type PresignedPutUrlRequest struct {
	FileName string             `json:"filename"`
	Folder   enum.StorageFolder `json:"folder"`
}

// PresignedPutUrlResponse — POST policy yuklash ma'lumoti.
// Mijoz Fields maydonlarini faylning OLDIDA yuborishi shart (S3 talabi).
type PresignedPutUrlResponse struct {
	URL        string            `json:"url"`
	Fields     map[string]string `json:"fields"`
	ObjectName string            `json:"object_name"`
}

func NewPresignedPutUrlResponse(url string, fields map[string]string, objectName string) *PresignedPutUrlResponse {
	return &PresignedPutUrlResponse{URL: url, Fields: fields, ObjectName: objectName}
}
