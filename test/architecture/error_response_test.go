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

// responseMethods — javob tanasini JSON qilib serializatsiya qiladigan
// metodlar. JsonResponse — loyihaning o'z Context wrapper'i
// (src/entrypoint/presentation/app/context/context_wrap.go), JSON — echo'niki.
var responseMethods = map[string]bool{
	"JsonResponse": true,
	"JSON":         true,
}

// TestNoBareErrorInJSONResponse xato obyektini javob tanasi sifatida
// berishni taqiqlaydi.
//
// Nima uchun bu buzilish: JsonResponse(status, data) data ni
// encoding/json ga uzatadi. Go'ning `error` interfeysi — bu metod,
// maydon emas. errors.errorString kabi standart turlarda eksport
// qilinadigan maydon umuman yo'q, shuning uchun json.Marshal ularni
// {} deb yozadi. Natijada mijoz {"data":{},"status":400} oladi —
// status kod to'g'ri, lekin sabab yo'qoladi. To'g'ri yozuv —
// err.Error(), u {"data":"social does not belong to user"} beradi.
//
// CHEGARA (ataylab qoldirilgan): qorovul argument NOMI bo'yicha
// ishlaydi, TURI bo'yicha emas — `err` deb nomlangan yalang'och
// identifikatorni qidiradi. Ya'ni:
//
//  1. YOLG'ON SALBIY: `e` yoki `parseErr` deb nomlangan xato o'tib
//     ketadi. Loyihada xato o'zgaruvchisi deyarli har doim `err` deb
//     nomlanadi, shuning uchun amaliy qamrov yetarli. To'liq turga
//     asoslangan tekshiruv go/types va butun paketni yuklashni talab
//     qiladi — narxi foydasidan katta.
//  2. YOLG'ON IJOBIY: agar kimdir `err` deb nomlangan, lekin xato
//     bo'lmagan o'zgaruvchini javob sifatida bersa, qorovul yiqiladi.
//     Bunday holatda to'g'ri yechim — o'zgaruvchini qayta nomlash,
//     qorovulni yumshatish EMAS.
func TestNoBareErrorInJSONResponse(t *testing.T) {
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
			if !ok {
				return true
			}
			if !isBareErrResponse(call) {
				return true
			}
			selector := call.Fun.(*ast.SelectorExpr)
			pos := fset.Position(call.Lparen)
			violations = append(violations, fmt.Sprintf(
				"%s:%d  .%s(..., err) — xato obyekti javob tanasi sifatida berilgan",
				relativeTo(root, pos.Filename), pos.Line, selector.Sel.Name))
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("src/ daraxtini o'qib bo'lmadi: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf(
			"Xato obyekti javob tanasiga berilgan (%d ta joy). json.Marshal uni\n"+
				"{} deb yozadi va mijoz sababni ko'rmaydi. err o'rniga err.Error()\n"+
				"bering, yoki xatoni *response.Response bilan qaytaring — u holda\n"+
				"ResponseMiddleware statusni o'zi to'g'ri qo'yadi.\n\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

// inlineErrConstructors — paket+funksiya nomlari bo'yicha jadval
// (`raw_sql_test.go` dagi `unsafeStringBuilders` naqshi bilan bir xil
// shaklda). Bularning har biri `error` qiymatini joyida quradi — ya'ni
// `err` deb nomlangan o'zgaruvchi bo'lmasa ham, natija baribir xato
// obyekti bo'lib, o'sha `{}` muammosiga tushadi.
var inlineErrConstructors = map[string]map[string]bool{
	"errors": {
		"New":  true,
		"Join": true,
	},
	"fmt": {
		"Errorf": true,
	},
}

// isBareErrResponse chaqiruv .JsonResponse(status, err) yoki
// .JSON(status, err) shaklida ekanligini aniqlaydi. Ikkinchi argument
// ikki shaklda tutiladi: yalang'och `err` identifikatori (nomga
// asoslangan, yuqoridagi CHEGARA'ga qarang) yoki inline xato konstruktori
// — masalan `errors.New("x")` yoki `fmt.Errorf(...)` — bular sof
// sintaktik aniqlanadi, chunki natija turi doim `error`.
//
// CHEGARA: nishon metodlar qat'iy 2-argumentli —
// `JsonResponse(status int, data any)` (context_wrap.go:31) va echo'ning
// `JSON(code int, i interface{})`. Shuning uchun faqat `len(call.Args) < 2`
// tekshiriladi, `> 2` emas — undan ortiq argument bu ikki metod uchun
// umuman kompilyatsiya qilinmaydi.
func isBareErrResponse(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || !responseMethods[selector.Sel.Name] {
		return false
	}
	if len(call.Args) < 2 {
		return false
	}

	switch arg := call.Args[1].(type) {
	case *ast.Ident:
		return arg.Name == "err"
	case *ast.CallExpr:
		inner, ok := arg.Fun.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		pkg, ok := inner.X.(*ast.Ident)
		return ok && inlineErrConstructors[pkg.Name][inner.Sel.Name]
	default:
		return false
	}
}

// TestBareErrResponseDetectorHasTeeth qorovul mantiqining o'zini tekshiradi.
// src/ toza bo'lgach TestNoBareErrorInJSONResponse har doim yashil bo'ladi
// va u endi hech narsani isbotlamaydi — bu test esa detektorning haqiqatan
// ishlashini src/ holatidan qat'i nazar kafolatlaydi.
func TestBareErrResponseDetectorHasTeeth(t *testing.T) {
	const src = `package p

func handle(c ctx, err error, msg string) error {
	if bad {
		return c.JsonResponse(400, err)                  // want: buzilish
	}
	if alsoBad {
		return c.JSON(500, err)                          // want: buzilish
	}
	if inlineNew {
		return c.JsonResponse(400, errors.New("x"))      // want: buzilish
	}
	if inlineJoin {
		return c.JSON(400, errors.Join(err, err))        // want: buzilish
	}
	if inlineErrorf {
		return c.JsonResponse(400, fmt.Errorf("x: %w", err)) // want: buzilish
	}
	c.JsonResponse(400, err.Error())                     // xavfsiz
	c.JsonResponse(400, msg)                             // xavfsiz
	c.JsonResponse(204, nil)                             // xavfsiz
	c.Logger().Error(err)                                // xavfsiz: JSON javob emas
	c.JsonResponse(400, strings.Join(parts, ","))        // xavfsiz: ro'yxatda yo'q funksiya
	return nil
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "detector_teeth.go", src, 0)
	if err != nil {
		t.Fatalf("namuna manbani tahlil qilib bo'lmadi: %v", err)
	}

	var flagged []int
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if ok && isBareErrResponse(call) {
			flagged = append(flagged, fset.Position(call.Lparen).Line)
		}
		return true
	})

	want := []int{5, 8, 11, 14, 17}
	if len(flagged) != len(want) {
		t.Fatalf("qorovul %d ta buzilish topdi, %d ta kutilgan edi: %v",
			len(flagged), len(want), flagged)
	}
	for i, line := range want {
		if flagged[i] != line {
			t.Errorf("%d-buzilish %d-qatorda topildi, %d kutilgan edi",
				i+1, flagged[i], line)
		}
	}
}
