package architecture_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// forbiddenJSONKeys — bu maydonlar hech qachon JSON javobida chiqmasligi kerak.
// Telefon ataylab yo'q: u UserBasicEntity da (jurnal a'zolari uchun) qoladi va
// use case darajasida so'rovchiga qarab redaksiya qilinadi.
//
// Qorovul chegarasi: telefon (json:"phone_number", json:"phone") va email
// JSON teg darajasida umuman taqiqlanmaydi va bu qorovul ularni tekshirmaydi.
// Sabab — ular ba'zi kontekstlarda ATAYLAB struct darajasida qoladi:
//   - UserBasicEntity.PhoneNumber — jurnal a'zolari ro'yxati, endpoint
//     a'zolikni allaqachon tekshiradi (/journal-manage/{id}/members).
//   - UserProfileEntity.Phone/Email — foydalanuvchi o'z profili, use case
//     ichida so'rovchiga qarab (o'zi/admin) redaksiya qilinadi.
//   - UserEntity.PhoneNumber/Email (ariza detali yo'llari) — use case
//     darajasida RedactUserContact/RedactApplicationContacts bilan
//     redaksiya qilinadi (permissionusecases.RedactContact* oilasi).
//
// Ya'ni telefon/email struct teg darajasida emas, so'rovchiga-qarab
// (deny-by-default) use case mantig'i bilan himoyalanadi — bu global
// go/ast qorovuli bilan ifodalab bo'lmaydigan talab, chunki qorovul
// so'rovchi kontekstini bilmaydi. PINFL/birth_date esa hech qachon hech
// kimga kerak emas (3.1-bo'lim), shuning uchun ular uchun global taqiq
// to'g'ri yechim.
var forbiddenJSONKeys = map[string]bool{
	"pin":        true,
	"pinfl":      true,
	"birth_date": true,
}

// TestNoSensitiveJSONTags maxfiy maydonlarning JSON teg orqali oshkor bo'lishini
// taqiqlaydi (CWE-200). Kelajakda kimdir json:"pin" qo'shsa, test yiqiladi.
func TestNoSensitiveJSONTags(t *testing.T) {
	roots := []string{
		sourceDir(t, "core", "domain", "entity"),
		sourceDir(t, "entrypoint", "presentation", "handlers"),
	}

	var violations []string
	for _, root := range roots {
		walkStructTags(t, root, func(file string, line int, jsonKey string) {
			if forbiddenJSONKeys[jsonKey] {
				violations = append(violations, fmt.Sprintf("%s:%d  json:%q", file, line, jsonKey))
			}
		})
	}

	if len(violations) > 0 {
		t.Errorf(
			"maxfiy maydon JSON teg orqali oshkor bo'lgan (%d ta joy).\n"+
				"PINFL va tug'ilgan sana API javoblarida chiqmasligi kerak;\n"+
				"json:\"-\" ishlating.\n\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

// walkStructTags har bir struct maydonining json teg kalitini callback'ga beradi.
func walkStructTags(t *testing.T, root string, fn func(file string, line int, jsonKey string)) {
	t.Helper()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		f, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			return fmt.Errorf("%s: %w", path, perr)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			field, ok := n.(*ast.Field)
			if !ok || field.Tag == nil {
				return true
			}
			tagValue := strings.Trim(field.Tag.Value, "`")
			jsonTag := reflect.StructTag(tagValue).Get("json")
			if jsonTag == "" {
				return true
			}
			key := strings.Split(jsonTag, ",")[0]
			if key == "" || key == "-" {
				return true
			}
			fn(path, fset.Position(field.Pos()).Line, key)
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("daraxtni o'qib bo'lmadi: %v", err)
	}
}

func sourceDir(t *testing.T, parts ...string) string {
	t.Helper()
	all := append([]string{"..", "..", "src"}, parts...)
	dir, err := filepath.Abs(filepath.Join(all...))
	if err != nil {
		t.Fatalf("yo'lni aniqlab bo'lmadi: %v", err)
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Fatalf("katalog topilmadi: %s", dir)
	}
	return dir
}
