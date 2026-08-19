package enum_test

import (
	"testing"

	"slib.uz/src/core/domain/entity/enum"
)

func TestAllowedExtensionsAreFound(t *testing.T) {
	cases := []struct {
		ext         string
		contentType string
		disposition string
	}{
		// PDF inline: brauzer uni o'z viewer'ida ko'rsata oladi, majburiy
		// yuklab olish esa asosiy foydalanish stsenariysini buzardi.
		{".pdf", "application/pdf", enum.DispositionInline},
		{".jpg", "image/jpeg", enum.DispositionInline},
		{".jpeg", "image/jpeg", enum.DispositionInline},
		{".png", "image/png", enum.DispositionInline},
		{".doc", "application/msword", enum.DispositionAttachment},
		{".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", enum.DispositionAttachment},
	}

	for _, c := range cases {
		ut, ok := enum.LookupUploadType(c.ext)
		if !ok {
			t.Fatalf("%s ruxsat etilgan bo'lishi kerak edi, topilmadi", c.ext)
		}
		if ut.ContentType != c.contentType {
			t.Errorf("%s: content-type %q kutilgandi, %q keldi", c.ext, c.contentType, ut.ContentType)
		}
		if ut.Disposition != c.disposition {
			t.Errorf("%s: disposition %q kutilgandi, %q keldi", c.ext, c.disposition, ut.Disposition)
		}
	}
}

func TestLookupIsCaseInsensitive(t *testing.T) {
	for _, ext := range []string{".PDF", ".pdf", ".pDf"} {
		ut, ok := enum.LookupUploadType(ext)
		if !ok {
			t.Fatalf("%s topilishi kerak edi", ext)
		}
		if ut.ContentType != "application/pdf" {
			t.Errorf("%s: application/pdf kutilgandi, %q keldi", ext, ut.ContentType)
		}
	}
}

func TestDangerousExtensionsAreRejected(t *testing.T) {
	// .svg ataylab ro'yxatda yo'q: SVG ichida <script> ishlaydi.
	for _, ext := range []string{".html", ".htm", ".svg", ".exe", ".php", ".js", ".sh", ""} {
		if _, ok := enum.LookupUploadType(ext); ok {
			t.Errorf("%q rad etilishi kerak edi, lekin qabul qilindi", ext)
		}
	}
}

func TestMagicBytesMatch(t *testing.T) {
	pdf, _ := enum.LookupUploadType(".pdf")
	if !pdf.MatchesMagic([]byte("%PDF-1.7\n")) {
		t.Error("haqiqiy PDF boshi qabul qilinishi kerak edi")
	}
	if pdf.MatchesMagic([]byte("<html><scr")) {
		t.Error("HTML baytlari PDF sifatida QABUL QILINMASLIGI kerak edi")
	}

	png, _ := enum.LookupUploadType(".png")
	if !png.MatchesMagic([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		t.Error("haqiqiy PNG boshi qabul qilinishi kerak edi")
	}

	jpg, _ := enum.LookupUploadType(".jpg")
	if !jpg.MatchesMagic([]byte{0xFF, 0xD8, 0xFF, 0xE0}) {
		t.Error("haqiqiy JPEG boshi qabul qilinishi kerak edi")
	}
}

func TestShortInputIsRejectedNotPanicked(t *testing.T) {
	png, _ := enum.LookupUploadType(".png")
	if png.MatchesMagic([]byte{0x89, 0x50}) {
		t.Error("imzodan qisqa kirish rad etilishi kerak edi")
	}
	if png.MatchesMagic(nil) {
		t.Error("bo'sh kirish rad etilishi kerak edi")
	}
}
