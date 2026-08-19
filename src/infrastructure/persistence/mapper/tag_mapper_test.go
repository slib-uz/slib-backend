package mapper

import (
	"testing"

	"slib.uz/src/infrastructure/persistence/models"
)

func TestTagModelToEntity_mapsNameAndLang(t *testing.T) {
	model := &models.TagModel{ID: 15, Name: "Sun'iy intellekt", Lang: "uz"}
	got := TagModelToEntity(model)
	if got.ID != 15 || got.Name != "Sun'iy intellekt" || got.Lang != "uz" {
		t.Fatalf("got %#v", got)
	}
}

func TestTagModelsToNamesByLang_groupsByLang(t *testing.T) {
	got := TagModelsToNamesByLang([]*models.TagModel{
		{ID: 1, Name: "A", Lang: "uz"},
		{ID: 2, Name: "B", Lang: "uz"},
		{ID: 3, Name: "C", Lang: "en"},
	})
	if len(got["uz"]) != 2 || got["uz"][0] != "A" || got["uz"][1] != "B" {
		t.Fatalf("uz = %#v", got["uz"])
	}
	if len(got["en"]) != 1 || got["en"][0] != "C" {
		t.Fatalf("en = %#v", got["en"])
	}
}
