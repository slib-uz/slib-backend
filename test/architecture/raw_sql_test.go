// Package architecture_test butun kod bazasiga tegishli qoidalarni tekshiradi.
package architecture_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// guardedMethods — GORM'ning birinchi argumenti SQL matni bo'lgan metodlari.
// Qolgan argumentlar ? ga bog'lanadigan qiymatlar, ular xavfsiz.
//
// CHEGARA (yolg'on ijobiy xavfi): qorovul metod NOMI bo'yicha ishlaydi, tur
// bo'yicha emas — chaqiruv joyida qaysi qabul qiluvchining metodi
// chaqirilayotgani tekshirilmaydi. Bir nechta shu ro'yxatdagi nomlar GORM'ga
// tegishli bo'lmagan holda ham loyihada ishlatiladi: masalan echo'ning
// marshrutlash `Group(prefix string) *echo.Group` metodi
// (src/entrypoint/presentation/groups/ da ~40 chaqiruv) xuddi shu nom bilan
// mavjud. Xuddi shu xavf `Select`/`Table`/`Not`/`Or`/`Exec` kabi umumiy
// nomlar uchun ham bor. Bugun bu chaqiruvlarning barchasi literal qator
// bilan qilinadi, shuning uchun yolg'on ijobiy yo'q — lekin kimdir ertaga
// `group.Group("/v" + version)` deb yozsa, qorovul buni "SQL inyeksiya" deb
// yiqitadi, garchi echo'ning Group'i SQL bilan aloqasi yo'q bo'lsa ham.
// Bunday holatda to'g'ri yechim SQL bo'lagini (yoki marshrut prefiksini)
// nomlangan konstantaga chiqarish, qorovulni yumshatish EMAS.
var guardedMethods = map[string]bool{
	"Order":    true,
	"Where":    true,
	"Or":       true,
	"Not":      true,
	"Having":   true,
	"Group":    true,
	"Select":   true,
	"Table":    true,
	"Joins":    true,
	"Raw":      true,
	"Exec":     true,
	"Pluck":    true,
	"Distinct": true,
}

// TestNoConcatenatedSQL SQL matniga foydalanuvchi kiritishini birlashtirishni
// taqiqlaydi (CWE-89).
//
// Nima uchun aynan birinchi argument: GORM'da birinchi argument har doim SQL
// matni. Shuning uchun Where("name ILIKE ?", "%"+name+"%") o'tadi — u yerdagi
// + qiymat argumentida va ? orqali parametrlashtiriladi. Where("name = '" +
// name + "'") esa yiqiladi.
//
// Qorovulning ikkita chegarasi bor, va ikkalasi ham ataylab qoldirilgan —
// murakkablikni oshirish foydasidan katta bo'lardi:
//
// 1. PROVENANCE: qorovul SINTAKTIK. Chaqiruv joyidagi ifodani ko'radi,
// o'zgaruvchining qayerdan kelganini emas. journal_repository_impl.go dagi
// GetJournalStatisticsV2 filterClause ni += bilan quradi va Raw ga
// o'zgaruvchi sifatida beradi — qorovul buni o'tkazib yuboradi. O'sha kod
// xavfsiz, lekin buni qorovul ISBOTLAMAYDI. Xavfsiz yo'lning haqiqatan
// xavfsizligini sorting paketining testlari kafolatlaydi.
//
// 2. NOM TO'QNASHUVI: qorovul guardedMethods'dagi metod NOMI bo'yicha
// ishlaydi, tur bo'yicha emas. GORM'ga tegishli bo'lmagan bir xil nomli
// metodlar (masalan, echo'ning marshrutlash Group(prefix) metodi) yolg'on
// ijobiy berishi mumkin — guardedMethods e'lonidagi izohga qarang.
//
// Himoya shu ikki chegara bilan birga ikki qismdan iborat va ikkalasi ham
// kerak: qorovul xavfli SHAKLni to'sadi, sorting testlari xavfsiz yo'lning
// haqiqatan xavfsiz ekanini isbotlaydi.
func TestNoConcatenatedSQL(t *testing.T) {
	root := sourceRoot(t)

	var violations []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		file, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return fmt.Errorf("%s: %w", path, parseErr)
		}

		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || !guardedMethods[selector.Sel.Name] {
				return true
			}
			if reason := unsafeSQLExpr(call.Args[0]); reason != "" {
				pos := fset.Position(call.Lparen)
				violations = append(violations, fmt.Sprintf(
					"%s:%d  .%s(...) birinchi argumentida %s",
					relativeTo(root, pos.Filename), pos.Line, selector.Sel.Name, reason))
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("src/ daraxtini o'qib bo'lmadi: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf(
			"SQL matni foydalanuvchi kiritishi bilan birlashtirilgan (%d ta joy).\n"+
				"Tartiblash uchun src/infrastructure/persistence/sorting paketidan\n"+
				"foydalaning; boshqa hollarda SQL bo'lagini konstanta jadvalidan oling.\n\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

// unsafeStringBuilders — paket+funksiya nomlari bo'yicha jadval. Bularning
// har biri qatorlarni ish vaqtida yig'adi (formatlaydi, birlashtiradi,
// takrorlaydi, almashtiradi) — ya'ni foydalanuvchi kiritishini SQL matniga
// aralashtirish uchun aynan shu funksiyalar ishlatiladi.
var unsafeStringBuilders = map[string]map[string]bool{
	"fmt": {
		"Sprintf":  true,
		"Sprint":   true,
		"Sprintln": true,
	},
	"strings": {
		"Join":       true,
		"Repeat":     true,
		"Replace":    true,
		"ReplaceAll": true,
	},
}

// unsafeSQLExpr ifoda ichida SQL matnini yig'ish belgilarini qidiradi.
// Bo'sh qator qaytsa — ifoda xavfsiz.
func unsafeSQLExpr(arg ast.Expr) string {
	reason := ""
	ast.Inspect(arg, func(node ast.Node) bool {
		if reason != "" {
			return false
		}
		switch typed := node.(type) {
		case *ast.BinaryExpr:
			if typed.Op == token.ADD {
				reason = "qatorlarni + bilan birlashtirish bor"
				return false
			}
		case *ast.CallExpr:
			if selector, ok := typed.Fun.(*ast.SelectorExpr); ok {
				pkg, isIdent := selector.X.(*ast.Ident)
				if isIdent && unsafeStringBuilders[pkg.Name][selector.Sel.Name] {
					reason = fmt.Sprintf("%s.%s chaqiruvi bor", pkg.Name, selector.Sel.Name)
					return false
				}
			}
		}
		return true
	})
	return reason
}

// sourceRoot src/ katalogining yo'lini qaytaradi. Test test/architecture/
// katalogida ishlaydi, shuning uchun ikki pog'ona yuqoriga chiqiladi.
func sourceRoot(t *testing.T) string {
	t.Helper()

	root, err := filepath.Abs(filepath.Join("..", "..", "src"))
	if err != nil {
		t.Fatalf("src/ yo'lini aniqlab bo'lmadi: %v", err)
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		t.Fatalf("src/ katalogi topilmadi: %s", root)
	}
	return root
}

func relativeTo(root, path string) string {
	rel, err := filepath.Rel(filepath.Dir(root), path)
	if err != nil {
		return path
	}
	return rel
}

// TestUnsafeSQLExprDetectsKnownDangerousShapes qorovulning o'zini sinaydi:
// unsafeSQLExpr bu paketda e'lon qilingan, shuning uchun to'g'ridan-to'g'ri
// chaqirish mumkin. Kimdir uni buzsa (masalan bir shoxni olib tashlasa),
// TestNoConcatenatedSQL o'zi buni sezmaydi — src/ hozircha toza bo'lgani
// uchun test yashil qolaveradi. Bu test qorovulning tishini alohida
// tekshiradi, src/ holatidan mustaqil.
func TestUnsafeSQLExprDetectsKnownDangerousShapes(t *testing.T) {
	cases := []struct {
		name   string
		expr   string
		unsafe bool
	}{
		{"binary string concat", `a + b`, true},
		{"fmt.Sprintf", `fmt.Sprintf("x %s", y)`, true},
		{"strings.Join", `strings.Join(parts, " ")`, true},
		{"literal + variable concat", `"lit" + v`, true},

		{"string literal", `"literal"`, false},
		{"bare identifier", `orderExpr`, false},
		{"non-listed function call", `joinOr(conditions)`, false},
		{"enum selector", `enum.Something`, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tc.expr)
			if err != nil {
				t.Fatalf("ifodani tahlil qilib bo'lmadi %q: %v", tc.expr, err)
			}

			reason := unsafeSQLExpr(expr)
			gotUnsafe := reason != ""

			if gotUnsafe != tc.unsafe {
				t.Errorf("%q: xavfli=%v kutilgandi, xavfli=%v keldi (sabab: %q)",
					tc.expr, tc.unsafe, gotUnsafe, reason)
			}
		})
	}
}
