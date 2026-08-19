package document

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
	"slib.uz/src/core/domain/ports/document"
)

type QRGeneratorImpl struct {
}

// @inject
func NewQRGenerator() document.QRGenerator {
	return &QRGeneratorImpl{}
}

func (this *QRGeneratorImpl) GenerateQRBase64(data string) (string, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)

	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}
