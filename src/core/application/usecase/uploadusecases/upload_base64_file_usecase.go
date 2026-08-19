package uploadusecases

import (
	"encoding/base64"
	"fmt"
	"strings"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/core/utils"
)

type UploadBase64FileUseCase struct {
	storage storage.FileStorage
}

// @inject
func NewUploadBase64FileUseCase(storage storage.FileStorage) *UploadBase64FileUseCase {
	return &UploadBase64FileUseCase{storage: storage}
}

func (this *UploadBase64FileUseCase) Execute(base64Data string, folder enum.StorageFolder, bucket enum.Bucket) (string, error) {
	data, err := this.decodeBase64(base64Data)
	if err != nil {
		return "", response.NewFailResponse(400, err.Error())
	}

	if err := this.validatePDF(data); err != nil {
		return "", response.NewFailResponse(400, err.Error())
	}

	objectName := this.generateObjectName()

	objectPath, err := this.storage.PutObject(bucket, folder, objectName, data)
	if err != nil {
		return "", response.NewFailResponse(500, "failed to upload file to storage")
	}

	return objectPath, nil
}

func (this *UploadBase64FileUseCase) decodeBase64(base64Data string) ([]byte, error) {
	base64Data = strings.TrimSpace(base64Data)

	if strings.HasPrefix(base64Data, "data:") {
		parts := strings.SplitN(base64Data, ",", 2)
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 data: %w", err)
	}

	return data, nil
}

// validatePDF imzoni allow-list orqali tekshiradi.
// Ilgari bu yerda "%PDF-" prefiksi qo'lda yozilgan edi — PDF imzosining
// ikkinchi nusxasi. Endi manba bitta: enum.LookupUploadType.
func (*UploadBase64FileUseCase) validatePDF(data []byte) error {
	pdfType, ok := enum.LookupUploadType(".pdf")
	if !ok || !pdfType.MatchesMagic(data) {
		return response.InvalidFileError
	}

	return nil
}

func (this *UploadBase64FileUseCase) generateObjectName() string {
	prefix := utils.RandomHex(8)
	return fmt.Sprintf("file_%s.pdf", prefix)
}
