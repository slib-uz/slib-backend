package utils_test

import (
	"testing"

	"slib.uz/src/core/utils"
)

// TestIsValidPhoneNumber formatning chegaralarini qulflaydi. Qoida: aynan
// "998" + 9 ta raqam. Normalizatsiya yo'q — "+", probel va defis rad etiladi.
func TestIsValidPhoneNumber(t *testing.T) {
	cases := []struct {
		name  string
		phone string
		want  bool
	}{
		{"to'g'ri raqam", "998901234567", true},
		{"istalgan 2 xonali prefiks qabul qilinadi", "998331234567", true},
		{"bir raqam ortiq", "9989012345678", false},
		{"bir raqam kam", "99890123456", false},
		{"998 prefiksi yo'q", "901234567", false},
		{"plus bilan", "+998901234567", false},
		{"probelli", "998 90 123 45 67", false},
		{"defisli", "998-90-123-45-67", false},
		{"boshqa mamlakat kodi", "997901234567", false},
		{"harf aralashgan", "99890123456a", false},
		{"butunlay harf", "abc", false},
		{"bo'sh satr", "", false},
		{"oldida qo'shimcha matn", "x998901234567", false},
		{"orqasida qo'shimcha matn", "998901234567x", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := utils.IsValidPhoneNumber(tc.phone); got != tc.want {
				t.Errorf("IsValidPhoneNumber(%q) = %v, %v kutilgandi", tc.phone, got, tc.want)
			}
		})
	}
}
