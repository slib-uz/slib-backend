package entity

import (
	"encoding/json"
	"testing"
)

func TestTagNamesByLang_unmarshalsNestedMap(t *testing.T) {
	var tags TagNamesByLang
	err := json.Unmarshal([]byte(`{"uz":["A","B"],"en":["AI"]}`), &tags)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags["uz"]) != 2 || tags["uz"][0] != "A" || tags["en"][0] != "AI" {
		t.Fatalf("got %#v", tags)
	}
}

func TestTagNamesByLang_unmarshalsLegacyStringArrayAsUz(t *testing.T) {
	var tags TagNamesByLang
	err := json.Unmarshal([]byte(`["Sun'iy intellekt","ML"]`), &tags)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags["uz"]) != 2 || tags["uz"][0] != "Sun'iy intellekt" || tags["uz"][1] != "ML" {
		t.Fatalf("got %#v", tags)
	}
}

func TestTagNamesByLang_nullIsEmpty(t *testing.T) {
	var tags TagNamesByLang
	if err := json.Unmarshal([]byte(`null`), &tags); err != nil {
		t.Fatal(err)
	}
	if tags != nil {
		t.Fatalf("got %#v", tags)
	}
}

func TestTagNamesByLangFromEntities_groupsByLang(t *testing.T) {
	got := TagNamesByLangFromEntities([]*TagEntity{
		{ID: 1, Name: "A", Lang: "uz"},
		{ID: 2, Name: "B", Lang: "en"},
		{ID: 3, Name: "C", Lang: ""},
	})
	if got["uz"][0] != "A" || got["uz"][1] != "C" || got["en"][0] != "B" {
		t.Fatalf("got %#v", got)
	}
}
