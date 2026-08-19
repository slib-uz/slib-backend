package utils

import (
	"errors"
	"testing"
)

func TestNormalizeTagNamesByLang_trimsAndLowercasesLang(t *testing.T) {
	got, err := NormalizeTagNamesByLang(map[string][]string{
		" UZ ": {"  Sun'iy intellekt  ", ""},
		"en":   {"AI", "AI"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %#v", got)
	}
	byLang := map[string]string{}
	for _, p := range got {
		byLang[p.Lang] = p.Name
	}
	if byLang["uz"] != "Sun'iy intellekt" || byLang["en"] != "AI" {
		t.Fatalf("got %#v", got)
	}
}

func TestNormalizeTagNamesByLang_rejectsLongName(t *testing.T) {
	_, err := NormalizeTagNamesByLang(map[string][]string{
		"uz": {"this-name-is-definitely-longer-than-32"},
	})
	if !errors.Is(err, ErrTagNameTooLong) {
		t.Fatalf("err = %v, want ErrTagNameTooLong", err)
	}
}

func TestNormalizeTagNamesByLang_rejectsLongLanguageCode(t *testing.T) {
	_, err := NormalizeTagNamesByLang(map[string][]string{
		"uz-lat": {"ok"},
	})
	if !errors.Is(err, ErrLanguageCodeTooLong) {
		t.Fatalf("err = %v, want ErrLanguageCodeTooLong", err)
	}
}
