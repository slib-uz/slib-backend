package sorting_test

import (
	"errors"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/infrastructure/persistence/sorting"
)

// newsFields — testlar uchun namunaviy ro'yxat. Haqiqiy ro'yxatlar
// repozitoriylarda e'lon qilinadi; bu yerda faqat xatti-harakat sinaladi.
func newsFields() sorting.Whitelist {
	return sorting.New("news.created_at DESC", map[string]string{
		"created_at": "news.created_at",
		"views":      "news.views_count",
	})
}

func TestResolveAllowedFieldAscending(t *testing.T) {
	got, err := newsFields().Resolve("created_at")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "news.created_at ASC" {
		t.Errorf("%q kutilgandi, %q keldi", "news.created_at ASC", got)
	}
}

func TestResolveMinusPrefixMeansDescending(t *testing.T) {
	got, err := newsFields().Resolve("-created_at")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "news.created_at DESC" {
		t.Errorf("%q kutilgandi, %q keldi", "news.created_at DESC", got)
	}
}

func TestResolveEmptyReturnsDefault(t *testing.T) {
	got, err := newsFields().Resolve("")
	if err != nil {
		t.Fatalf("bo'sh qiymat xatolik bermasligi kerak: %v", err)
	}
	if got != "news.created_at DESC" {
		t.Errorf("standart tartib kutilgandi, %q keldi", got)
	}
}

func TestResolveUnknownFieldIsRejected(t *testing.T) {
	got, err := newsFields().Resolve("password_hash")
	if !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("InvalidSortFieldError kutilgandi, %v keldi", err)
	}
	if got != "" {
		t.Errorf("xatolik bilan birga bo'sh natija kutilgandi, %q keldi", got)
	}
}

func TestResolveBareMinusIsRejected(t *testing.T) {
	if _, err := newsFields().Resolve("-"); !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("InvalidSortFieldError kutilgandi, %v keldi", err)
	}
}

// Auditning aynan payload'i. Bu test zaiflikni to'g'ridan-to'g'ri mahkamlaydi.
func TestResolveRejectsAuditPayload(t *testing.T) {
	payload := "(SELECT(*)FROM(*)pg_sleep(10))"
	got, err := newsFields().Resolve(payload)
	if !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("hujum satri rad etilishi kerak edi, xatolik: %v", err)
	}
	if got != "" {
		t.Fatalf("hujum satri natijaga tushdi: %q", got)
	}
}

func TestResolvePairDirections(t *testing.T) {
	cases := map[string]string{
		"":     "news.created_at ASC",
		"asc":  "news.created_at ASC",
		"ASC":  "news.created_at ASC",
		"desc": "news.created_at DESC",
		"DESC": "news.created_at DESC",
	}
	for direction, want := range cases {
		got, err := newsFields().ResolvePair("created_at", direction)
		if err != nil {
			t.Errorf("%q: xatolik kutilmagandi: %v", direction, err)
			continue
		}
		if got != want {
			t.Errorf("%q: %q kutilgandi, %q keldi", direction, want, got)
		}
	}
}

func TestResolvePairRejectsUnknownDirection(t *testing.T) {
	if _, err := newsFields().ResolvePair("created_at", "asc; DROP TABLE news"); !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("InvalidSortFieldError kutilgandi, %v keldi", err)
	}
}

func TestResolvePairRejectsUnknownField(t *testing.T) {
	if _, err := newsFields().ResolvePair("pg_sleep(10)", "desc"); !errors.Is(err, response.InvalidSortFieldError) {
		t.Fatalf("InvalidSortFieldError kutilgandi, %v keldi", err)
	}
}

func TestResolvePairEmptyFieldReturnsDefault(t *testing.T) {
	got, err := newsFields().ResolvePair("", "desc")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if got != "news.created_at DESC" {
		t.Errorf("standart tartib kutilgandi, %q keldi", got)
	}
}

func TestFieldsIsSorted(t *testing.T) {
	got := newsFields().Fields()
	want := []string{"created_at", "views"}
	if len(got) != len(want) {
		t.Fatalf("%v kutilgandi, %v keldi", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%v kutilgandi, %v keldi", want, got)
		}
	}
}
