package repository_test

import (
	"testing"

	"slib.uz/src/infrastructure/persistence/repository"
)

func assertFields(t *testing.T, name string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: %v kutilgandi, %v keldi", name, want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s: %v kutilgandi, %v keldi", name, want, got)
		}
	}
}

// Ro'yxatlar handler'lardagi swagger Enums(...) hujjatiga mos bo'lishi kerak.
// Bu test API kontraktini kodga mahkamlaydi: kimdir ro'yxatga maydon qo'shsa
// yoki olib tashlasa, hujjat bilan farq shu yerda ko'rinadi.
func TestArticleSortFieldsMatchSwaggerEnums(t *testing.T) {
	// articles_list_handler.go: Enums(views_count,rating_sum,publication_date)
	assertFields(t, "articles",
		repository.ArticleSortFields.Fields(),
		[]string{"publication_date", "rating_sum", "views_count"})
}

func TestJournalSortFieldsMatchSwaggerEnums(t *testing.T) {
	// journal_list_handler.go: Enums(views_count,rating_sum,established_date)
	assertFields(t, "journals",
		repository.JournalSortFields.Fields(),
		[]string{"established_date", "rating_sum", "views_count"})
}

func TestNewsSortFieldsMatchSwaggerEnums(t *testing.T) {
	// news_list_handler.go: Enums(created_at,-created_at)
	assertFields(t, "news",
		repository.NewsSortFields.Fields(),
		[]string{"created_at"})
}

func TestReportSortFieldsMatchSwaggerEnums(t *testing.T) {
	// report_list_handler.go: Enums(created_at,-created_at)
	assertFields(t, "report",
		repository.ReportSortFields.Fields(),
		[]string{"created_at"})
}

func TestJournalRatingSortFieldsMatchSwaggerEnums(t *testing.T) {
	// journal_rating_list_handler.go: Enums(created_at,-created_at)
	assertFields(t, "journal_ratings",
		repository.JournalRatingSortFields.Fields(),
		[]string{"created_at"})
}

// ReportModel.TableName() "report" qaytaradi — BIRLIKDA. GORM'ning standart
// ko'plik qoidasiga ishonib "reports" deb yozilsa, so'rov ishga tushganda
// relation "reports" does not exist xatosi chiqadi.
func TestReportSortColumnUsesSingularTableName(t *testing.T) {
	got, err := repository.ReportSortFields.Resolve("created_at")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "report.created_at ASC" {
		t.Errorf("%q kutilgandi, %q keldi", "report.created_at ASC", got)
	}
}

func TestSupportDialogSortFieldsMatchSwaggerEnums(t *testing.T) {
	// chat_list_handler.go, support_dialog_list_question_handler.go,
	// support_dialog_list_answer_handler.go: Enums(created_at,-created_at)
	assertFields(t, "support_dialogs",
		repository.SupportDialogSortFields.Fields(),
		[]string{"created_at"})
}

// GetChatsByPaging so'rovi last_msg nomli quyi so'rovni JOIN qiladi.
// Ustun to'liq nom bilan yozilishi shu sababli muhim.
func TestSupportDialogSortColumnIsTableQualified(t *testing.T) {
	got, err := repository.SupportDialogSortFields.Resolve("-created_at")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "support_dialogs.created_at DESC" {
		t.Errorf("%q kutilgandi, %q keldi", "support_dialogs.created_at DESC", got)
	}
}

// ISSNConditions jadvalida ustun nomi emas, BUTUN shart saqlanadi.
// Aks holda repozitoriyda condition + " = ?" birlashtirish bo'lardi va
// qorovul test uni to'xtatardi.
func TestISSNConditionsAreCompleteSQLFragments(t *testing.T) {
	want := map[string]string{
		"issn_paper":  "journals.issn_paper = ?",
		"issn_online": "journals.issn_online = ?",
	}
	if len(repository.ISSNConditions) != len(want) {
		t.Fatalf("%d ta shart kutilgandi, %d keldi", len(want), len(repository.ISSNConditions))
	}
	for key, expected := range want {
		got, ok := repository.ISSNConditions[key]
		if !ok {
			t.Errorf("%q kaliti yo'q", key)
			continue
		}
		if got != expected {
			t.Errorf("%q: %q kutilgandi, %q keldi", key, expected, got)
		}
	}
}

// Maqolalar so'rovi qidiruv parametrlari berilganda journals jadvalini JOIN
// qiladi, va views_count/rating_sum ustunlari IKKALA jadvalda ham mavjud.
// Ustunlar to'liq nom bilan yozilmasa, PostgreSQL "ambiguous column" xatosini
// beradi — ya'ni hujjatlashtirilgan qiymat ishlamaydi.
func TestArticleSortColumnsAreTableQualified(t *testing.T) {
	for _, field := range repository.ArticleSortFields.Fields() {
		got, err := repository.ArticleSortFields.Resolve(field)
		if err != nil {
			t.Fatalf("%s: xatolik kutilmagandi: %v", field, err)
		}
		if len(got) < 9 || got[:9] != "articles." {
			t.Errorf("%s: ustun \"articles.\" bilan boshlanishi kerak, %q keldi", field, got)
		}
	}
}

// Fields() faqat kalitlarni tasdiqlaydi — ustunlarning haqiqiy SQL
// qiymatlarini emas. "reprot.created_at DESC" kabi imlo xatosi Fields()
// testidan bemalol o'tib ketadi va faqat ishga tushganda "relation does not
// exist" bilan portlaydi. Bu test har bir ro'yxatning standart tartibini
// (Resolve("")) so'zma-so'z tasdiqlaydi — spec §5 aynan shundan qo'rqadi.
func TestDefaultOrderExpressions(t *testing.T) {
	cases := []struct {
		name string
		list sortResolver
		want string
	}{
		{"articles", repository.ArticleSortFields, "articles.publication_date DESC"},
		{"journals", repository.JournalSortFields, "journals.established_date DESC"},
		{"news", repository.NewsSortFields, "news.created_at DESC"},
		{"report", repository.ReportSortFields, "report.created_at DESC"},
		{"journal_ratings", repository.JournalRatingSortFields, "journal_ratings.created_at DESC"},
		{"support_dialogs", repository.SupportDialogSortFields, "support_dialogs.created_at DESC"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.list.Resolve("")
			if err != nil {
				t.Fatalf("standart tartib uchun xatolik kutilmagandi: %v", err)
			}
			if got != tc.want {
				t.Errorf("%q kutilgandi, %q keldi", tc.want, got)
			}
		})
	}
}

// sortResolver — sorting.Whitelist ning Resolve metodini o'zida saqlaydigan
// har qanday qiymat uchun minimal interfeys. Faqat shu test faylida,
// jadval qatorlarini bir xil turga keltirish uchun ishlatiladi.
type sortResolver interface {
	Resolve(ordering string) (string, error)
}

// Har bir ro'yxatning har bir maydoni to'g'ri jadval prefiksi bilan
// boshlanishini tasdiqlaydi. TestArticleSortColumnsAreTableQualified va
// TestSupportDialogSortColumnIsTableQualified alohida-alohida shuni
// tekshiradi; bu test qolgan to'rtta ro'yxatni ham qamraydi.
func TestAllSortFieldsAreTableQualified(t *testing.T) {
	cases := []struct {
		name   string
		list   sortResolver
		fields []string
		prefix string
	}{
		{"journals", repository.JournalSortFields, repository.JournalSortFields.Fields(), "journals."},
		{"news", repository.NewsSortFields, repository.NewsSortFields.Fields(), "news."},
		{"report", repository.ReportSortFields, repository.ReportSortFields.Fields(), "report."},
		{"journal_ratings", repository.JournalRatingSortFields, repository.JournalRatingSortFields.Fields(), "journal_ratings."},
		{"support_dialogs", repository.SupportDialogSortFields, repository.SupportDialogSortFields.Fields(), "support_dialogs."},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, field := range tc.fields {
				got, err := tc.list.Resolve(field)
				if err != nil {
					t.Fatalf("%s: xatolik kutilmagandi: %v", field, err)
				}
				if len(got) < len(tc.prefix) || got[:len(tc.prefix)] != tc.prefix {
					t.Errorf("%s: ustun %q bilan boshlanishi kerak, %q keldi", field, tc.prefix, got)
				}
			}
		})
	}
}
