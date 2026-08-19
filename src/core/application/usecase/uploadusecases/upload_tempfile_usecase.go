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
