package storage_test

import (
	"strings"
	"testing"
	"time"

	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/storage"
)

func mustPolicyJSON(t *testing.T, ext string, maxSize int64) string {
	t.Helper()
	ut, ok := enum.LookupUploadType(ext)
	if !ok {
		t.Fatalf("test sozlash xatosi: %s allow-list'da yo'q", ext)
	}
	// Sobit vaqt: test natijasi soatga bog'liq bo'lmasligi kerak.
	expires := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := storage.BuildPostPolicy("slibbucketprivate", "articles/a_1234"+ext, ut, maxSize, expires)
	if err != nil {
		t.Fatalf("BuildPostPolicy xato qaytardi: %v", err)
	}
	return p.String()
}

func TestPolicyPinsContentType(t *testing.T) {
	got := mustPolicyJSON(t, ".pdf", 10<<20)
	if !strings.Contains(got, "application/pdf") {
		t.Errorf("policy Content-Type ni qat'iy belgilashi kerak edi.\npolicy: %s", got)
	}
	if !strings.Contains(got, "Content-Type") {
		t.Errorf("policy shartlarida Content-Type bo'lishi kerak edi.\npolicy: %s", got)
	}
}

func TestPolicyPinsContentLengthRange(t *testing.T) {
	got := mustPolicyJSON(t, ".pdf", 10<<20)
	if !strings.Contains(got, "content-length-range") {
		t.Errorf("policy hajm oralig'ini belgilashi kerak edi.\npolicy: %s", got)
	}
	if !strings.Contains(got, "10485760") {
		t.Errorf("policy'da maksimal hajm (10485760) bo'lishi kerak edi.\npolicy: %s", got)
	}
}

func TestPolicyPinsExactKey(t *testing.T) {
	got := mustPolicyJSON(t, ".pdf", 10<<20)
	if !strings.Contains(got, "articles/a_1234.pdf") {
		t.Errorf("policy obyekt kalitini qat'iy belgilashi kerak edi.\npolicy: %s", got)
	}
}

func TestPolicySetsDispositionFromUploadType(t *testing.T) {
	doc := mustPolicyJSON(t, ".docx", 10<<20)
	if !strings.Contains(doc, enum.DispositionAttachment) {
		t.Errorf("hujjat uchun attachment kutilgandi.\npolicy: %s", doc)
	}

	// PDF ommaviy bucket'da ham brauzerda ochilishi kerak — u yerda
	// aynan shu policy metadatasi obyektga muhrlanadi.
	if pdf := mustPolicyJSON(t, ".pdf", 10<<20); !strings.Contains(pdf, enum.DispositionInline) {
		t.Errorf("PDF uchun inline kutilgandi.\npolicy: %s", pdf)
	}

	ut, _ := enum.LookupUploadType(".png")
	expires := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := storage.BuildPostPolicy("slibbucketpublic", "news/x_1234.png", ut, 10<<20, expires)
	if err != nil {
		t.Fatalf("BuildPostPolicy xato qaytardi: %v", err)
	}
	if !strings.Contains(p.String(), enum.DispositionInline) {
		t.Errorf("rasm uchun inline kutilgandi.\npolicy: %s", p.String())
	}
}

func TestNilUploadTypeIsRejected(t *testing.T) {
	expires := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err := storage.BuildPostPolicy("b", "k.pdf", nil, 10<<20, expires); err == nil {
		t.Error("nil UploadType uchun xato kutilgandi")
	}
}
