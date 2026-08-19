package storage_test

import (
	"strings"
	"testing"

	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/storage"
)

// PDF brauzerda ochilishi kerak — majburiy yuklab olish emas.
func TestKnownPdfGetsInline(t *testing.T) {
	p := storage.ResponseParamsForObject("articles/maqola_1234.pdf")

	if got := p.Get("response-content-type"); got != "application/pdf" {
		t.Errorf("response-content-type: application/pdf kutilgandi, %q keldi", got)
	}
	disp := p.Get("response-content-disposition")
	if !strings.HasPrefix(disp, enum.DispositionInline) {
		t.Errorf("inline kutilgandi, %q keldi", disp)
	}
	if !strings.Contains(disp, "maqola_1234.pdf") {
		t.Errorf("disposition'da fayl nomi bo'lishi kerak edi, %q keldi", disp)
	}
}

// Brauzer render qila olmaydigan hujjatlar attachment bo'lib qoladi.
func TestKnownOfficeDocumentGetsAttachment(t *testing.T) {
	p := storage.ResponseParamsForObject("articles/maqola_1234.docx")

	disp := p.Get("response-content-disposition")
	if !strings.HasPrefix(disp, enum.DispositionAttachment) {
		t.Errorf("attachment kutilgandi, %q keldi", disp)
	}
}

func TestKnownImageGetsInline(t *testing.T) {
	p := storage.ResponseParamsForObject("news/muqova_1234.png")

	if got := p.Get("response-content-type"); got != "image/png" {
		t.Errorf("image/png kutilgandi, %q keldi", got)
	}
	if disp := p.Get("response-content-disposition"); !strings.HasPrefix(disp, enum.DispositionInline) {
		t.Errorf("inline kutilgandi, %q keldi", disp)
	}
}

// Eski, allow-list'dan oldin yuklangan fayllar eng xavfsiz holatga tushishi kerak.
func TestUnknownExtensionFallsBackToSafestValues(t *testing.T) {
	for _, name := range []string{"legacy/eski.html", "legacy/noext", "legacy/x.svg"} {
		p := storage.ResponseParamsForObject(name)

		if got := p.Get("response-content-type"); got != "application/octet-stream" {
			t.Errorf("%s: application/octet-stream kutilgandi, %q keldi", name, got)
		}
		if disp := p.Get("response-content-disposition"); !strings.HasPrefix(disp, enum.DispositionAttachment) {
			t.Errorf("%s: attachment kutilgandi, %q keldi", name, disp)
		}
	}
}
