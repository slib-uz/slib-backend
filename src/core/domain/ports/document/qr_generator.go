package document

type QRGenerator interface {
	GenerateQRBase64(data string) (string, error)
}
