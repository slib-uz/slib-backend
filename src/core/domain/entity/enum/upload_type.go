package enum

import (
	"bytes"
	"strings"
)

// Content-Disposition qiymatlari.
const (
	DispositionAttachment = "attachment"
	DispositionInline     = "inline"
)

// MagicPrefixLen — imzoni tekshirish uchun yetarli bosh baytlar soni.
// Ro'yxatdagi eng uzun imzo 8 bayt (PNG).
const MagicPrefixLen = 8

// UploadType ruxsat etilgan bitta fayl turini tavsiflaydi.
// Kengaytma, MIME turi, imzo va disposition bir joyda turadi — ular
// uch xil faylda yashasa, vaqt o'tishi bilan uzoqlashib ketadi.
type UploadType struct {
	Extension   string
	ContentType string
	Disposition string
	// Magic — ruxsat etilgan bosh baytlar variantlari.
	// Bo'sh bo'lsa imzo tekshirilmaydi (hozircha bunday tur yo'q).
	Magic [][]byte
}

var uploadTypes = map[string]*UploadType{
	".pdf": {
		Extension:   ".pdf",
		ContentType: "application/pdf",
		// Inline: brauzer PDF ni o'z viewer'ida ko'rsata oladi va kutubxona
		// uchun asosiy stsenariy aynan shu — o'qish, yuklab olish emas.
		// Xavfsizlik majburlangan `application/pdf` turiga tayanadi: bu tur
		// bilan brauzer HTML/skriptni ijro etmaydi, PDF esa sandbox'da ochiladi.
		Disposition: DispositionInline,
		Magic:       [][]byte{[]byte("%PDF-")},
	},
	".jpg": {
		Extension:   ".jpg",
		ContentType: "image/jpeg",
		Disposition: DispositionInline,
		Magic:       [][]byte{{0xFF, 0xD8, 0xFF}},
	},
	".jpeg": {
		Extension:   ".jpeg",
		ContentType: "image/jpeg",
		Disposition: DispositionInline,
		Magic:       [][]byte{{0xFF, 0xD8, 0xFF}},
	},
	".png": {
		Extension:   ".png",
		ContentType: "image/png",
		Disposition: DispositionInline,
		Magic:       [][]byte{{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
	},
	".doc": {
		Extension:   ".doc",
		ContentType: "application/msword",
		Disposition: DispositionAttachment,
		Magic:       [][]byte{{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
	},
	".docx": {
		Extension:   ".docx",
		ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		Disposition: DispositionAttachment,
		// PK\x03\x04 — ZIP konteyneri. Bu imzo .docx ni boshqa ZIP asosidagi
		// formatlardan ajratmaydi; konteyner ichini ochish spec doirasidan tashqarida.
		Magic: [][]byte{{0x50, 0x4B, 0x03, 0x04}},
	},
}

// LookupUploadType kengaytma bo'yicha ruxsat etilgan turni qaytaradi.
// Katta-kichik harf ahamiyatsiz. Topilmasa (nil, false) — chaqiruvchi
// yuklashni RAD ETISHI shart.
func LookupUploadType(ext string) (*UploadType, bool) {
	ut, ok := uploadTypes[strings.ToLower(ext)]
	return ut, ok
}

// MatchesMagic head baytlari shu turning imzolaridan biriga mos kelishini tekshiradi.
// head imzodan qisqa bo'lsa false qaytaradi — panic emas.
func (u *UploadType) MatchesMagic(head []byte) bool {
	if len(u.Magic) == 0 {
		return true
	}
	for _, sig := range u.Magic {
		if bytes.HasPrefix(head, sig) {
			return true
		}
	}
	return false
}
