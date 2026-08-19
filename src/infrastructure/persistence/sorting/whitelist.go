// Package sorting so'rov parametridan kelgan tartiblash nomini ruxsat etilgan
// SQL ustun ifodasiga aylantiradi.
//
// Nima uchun bu kerak: GORM'ning Order(string) metodi qiymatni
// clause.Column{Raw: true} qilib qo'yadi, ya'ni matn SQL'ga o'zgarishsiz
// tushadi. GORM'ning parametrlashtirish himoyasi faqat ? ga bog'lanadigan
// qiymatlarga tegishli; ustun nomlari undan tashqarida qoladi.
//
// Paket persistence qatlamida turadi, chunki ustun nomlari ma'lumotlar
// bazasining tushunchasi — domenning emas.
package sorting

import (
	"sort"
	"strings"

	"slib.uz/src/core/application/response"
)

const (
	ascending  = " ASC"
	descending = " DESC"
)

// Whitelist API'da ko'rinadigan tartiblash nomini SQL ifodasiga bog'laydi.
//
// columns qiymatlari va defaultOrder faqat kod ichidagi konstantalar bo'lishi
// SHART — ular SQL matniga o'zgarishsiz tushadi. Foydalanuvchi kiritishi
// natijaga hech qachon qo'shilmaydi: u faqat map kaliti sifatida ishlatiladi.
type Whitelist struct {
	columns      map[string]string
	defaultOrder string
}

// New ro'yxat yaratadi. defaultOrder — parametr berilmaganda ishlatiladigan
// to'liq ORDER BY ifodasi (masalan "articles.publication_date DESC").
func New(defaultOrder string, columns map[string]string) Whitelist {
	return Whitelist{columns: columns, defaultOrder: defaultOrder}
}

// Fields ruxsat etilgan API nomlarini alifbo tartibida qaytaradi.
// Kontrakt testlari uchun: ro'yxat swagger hujjatiga mosligini mahkamlaydi.
func (w Whitelist) Fields() []string {
	fields := make([]string, 0, len(w.columns))
	for field := range w.columns {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}

// Resolve bitta parametrni ishlaydi: "-created_at" -> "<ustun> DESC".
// Bo'sh qiymat standart tartibni beradi. Ro'yxatda yo'q nom — xatolik.
func (w Whitelist) Resolve(ordering string) (string, error) {
	if ordering == "" {
		return w.defaultOrder, nil
	}

	field := ordering
	direction := ascending
	if strings.HasPrefix(ordering, "-") {
		field = strings.TrimPrefix(ordering, "-")
		direction = descending
	}

	column, ok := w.columns[field]
	if !ok {
		return "", response.InvalidSortFieldError
	}

	// Ikkala bo'lak ham konstanta: column — jadvaldan, direction — yuqoridagi
	// const bloklardan. Foydalanuvchi kiritishi bu yerga tushmaydi.
	return column + direction, nil
}

// ResolvePair maydon va yo'nalish alohida parametr bo'lganda ishlatiladi
// (jurnallar ro'yxatidagi sort_by + order).
func (w Whitelist) ResolvePair(field, direction string) (string, error) {
	if field == "" {
		return w.defaultOrder, nil
	}

	column, ok := w.columns[field]
	if !ok {
		return "", response.InvalidSortFieldError
	}

	switch strings.ToUpper(direction) {
	case "", "ASC":
		return column + ascending, nil
	case "DESC":
		return column + descending, nil
	default:
		return "", response.InvalidSortFieldError
	}
}
