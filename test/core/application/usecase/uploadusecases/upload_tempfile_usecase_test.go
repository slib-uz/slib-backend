package uploadusecases_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/uploadusecases"
	"slib.uz/src/core/domain/ports/storage"
)

type fakeTempStorage struct {
	storage.FileStorage
	calls int
}

func (f *fakeTempStorage) SaveTempFile(_ *multipart.FileHeader) (string, error) {
	f.calls++
	return "uploads/saqlandi.pdf", nil
}

// makeFileHeader berilgan nom va mazmun bilan haqiqiy multipart.FileHeader yasaydi.
func makeFileHeader(t *testing.T, name string, content []byte) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	r := multipart.NewReader(&buf, w.Boundary())
	form, err := r.ReadForm(int64(len(content)) + 1024)
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}
	return form.File["file"][0]
}

func TestRealPdfIsAccepted(t *testing.T) {
	s := &fakeTempStorage{}
	uc := uploadusecases.NewUploadTempFileUseCase(s)

	fh := makeFileHeader(t, "maqola.pdf", []byte("%PDF-1.7\nhaqiqiy pdf mazmuni"))
	if _, err := uc.Execute(fh); err != nil {
		t.Fatalf("haqiqiy PDF qabul qilinishi kerak edi: %v", err)
	}
	if s.calls != 1 {
		t.Errorf("storage bir marta chaqilishi kerak edi, %d", s.calls)
	}
}

func TestHtmlBytesNamedPdfAreRejected(t *testing.T) {
	s := &fakeTempStorage{}
	uc := uploadusecases.NewUploadTempFileUseCase(s)

	fh := makeFileHeader(t, "zararsiz.pdf", []byte("<html><script>alert(1)</script></html>"))
	_, err := uc.Execute(fh)
	if !errors.Is(err, response.InvalidFileError) {
		t.Errorf("InvalidFileError kutilgandi, %v keldi", err)
	}
	if s.calls != 0 {
		t.Errorf("rad etilgan fayl saqlanmasligi kerak edi, %d marta saqlandi", s.calls)
	}
}

func TestDisallowedExtensionIsRejected(t *testing.T) {
	s := &fakeTempStorage{}
	uc := uploadusecases.NewUploadTempFileUseCase(s)

	fh := makeFileHeader(t, "evil.html", []byte("<html></html>"))
	_, err := uc.Execute(fh)
	if !errors.Is(err, response.InvalidFileError) {
		t.Errorf("InvalidFileError kutilgandi, %v keldi", err)
	}
	if s.calls != 0 {
		t.Errorf("storage chaqirilmasligi kerak edi")
	}
}

func TestTruncatedFileIsRejectedNotPanicked(t *testing.T) {
	s := &fakeTempStorage{}
	uc := uploadusecases.NewUploadTempFileUseCase(s)

	fh := makeFileHeader(t, "qisqa.png", []byte{0x89})
	_, err := uc.Execute(fh)
	if !errors.Is(err, response.InvalidFileError) {
		t.Errorf("InvalidFileError kutilgandi, %v keldi", err)
	}
}
