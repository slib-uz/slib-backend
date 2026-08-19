package entity

import "testing"

func TestFlattenStudyFieldCatalog_includesChildren(t *testing.T) {
	parentID := uint(1)
	fields := []*StudyFieldEntity{
		{
			ID:   1,
			Name: map[string]string{"uz": "Ota"},
			Children: []*StudyFieldEntity{
				{ID: 2, Name: map[string]string{"uz": "Bola"}, ParentID: &parentID},
			},
		},
	}

	got := FlattenStudyFieldCatalog(fields)
	if len(got) != 2 {
		t.Fatalf("len=%d want 2", len(got))
	}
	if got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("ids=%v", got)
	}
}

func TestArticleMetadataExtraction_NormalizeFillsLangKeys(t *testing.T) {
	meta := &ArticleMetadataExtraction{
		ArticleName: map[string]string{"uz": "Nom"},
		Tags:        TagNamesByLang{"uz": []string{"a"}},
	}
	meta.Normalize(nil)
	if meta.ArticleName["ru"] != "" || meta.ArticleName["en"] != "" {
		t.Fatalf("name=%v", meta.ArticleName)
	}
	if meta.Tags["en"] == nil || meta.StudyFieldIDs == nil || meta.References == nil {
		t.Fatalf("got %#v", meta)
	}
}

func TestNormalizeKeepsExtraLanguage(t *testing.T) {
	meta := &ArticleMetadataExtraction{
		ArticleName: map[string]string{"uz": "Nom", "de": "Name"},
	}
	meta.Normalize([]string{"uz", "ru", "en", "de"})
	if meta.ArticleName["de"] != "Name" {
		t.Fatalf("name=%v", meta.ArticleName)
	}
}

func TestMetadataLangs(t *testing.T) {
	uz := MetadataLangs("uz")
	if len(uz) != 3 {
		t.Fatalf("uz=%v", uz)
	}
	de := MetadataLangs(" DE ")
	if len(de) != 4 || de[3] != "de" {
		t.Fatalf("de=%v", de)
	}
}
