package entity

import "strings"

type StudyFieldCatalogItem struct {
	ID   uint              `json:"id"`
	Name map[string]string `json:"name"`
}

type ArticleMetadataExtraction struct {
	ArticleName     map[string]string `json:"article_name"`
	ArticleLanguage string            `json:"article_language"`
	StudyFieldIDs   []uint            `json:"study_field_ids"`
	DOI             string            `json:"doi"`
	Tags            TagNamesByLang    `json:"tags"`
	References      []string          `json:"references"`
	Annotation      map[string]string `json:"annotation"`
}

func MetadataLangs(articleLanguage string) []string {
	langs := []string{"uz", "ru", "en"}
	code := strings.ToLower(strings.TrimSpace(articleLanguage))
	if code == "" {
		return langs
	}
	for _, lang := range langs {
		if lang == code {
			return langs
		}
	}
	return append(langs, code)
}

func EmptyArticleMetadataExtraction(langs []string) *ArticleMetadataExtraction {
	langs = normalizeLangList(langs)
	return &ArticleMetadataExtraction{
		ArticleName:     emptyLocalizedString(langs),
		ArticleLanguage: "",
		StudyFieldIDs:   []uint{},
		DOI:             "",
		Tags:            emptyLocalizedTags(langs),
		References:      []string{},
		Annotation:      emptyLocalizedString(langs),
	}
}

func (this *ArticleMetadataExtraction) Normalize(langs []string) {
	if this == nil {
		return
	}
	langs = normalizeLangList(langs)
	this.ArticleName = ensureLocalizedString(this.ArticleName, langs)
	this.Annotation = ensureLocalizedString(this.Annotation, langs)
	this.Tags = ensureLocalizedTags(this.Tags, langs)
	if this.StudyFieldIDs == nil {
		this.StudyFieldIDs = []uint{}
	}
	if this.References == nil {
		this.References = []string{}
	}
}

func FlattenStudyFieldCatalog(fields []*StudyFieldEntity) []StudyFieldCatalogItem {
	var out []StudyFieldCatalogItem
	var walk func([]*StudyFieldEntity)
	walk = func(items []*StudyFieldEntity) {
		for _, item := range items {
			if item == nil {
				continue
			}
			out = append(out, StudyFieldCatalogItem{ID: item.ID, Name: item.Name})
			if len(item.Children) > 0 {
				walk(item.Children)
			}
		}
	}
	walk(fields)
	if out == nil {
		return []StudyFieldCatalogItem{}
	}
	return out
}

func CatalogIDSet(items []StudyFieldCatalogItem) map[uint]struct{} {
	out := make(map[uint]struct{}, len(items))
	for _, item := range items {
		out[item.ID] = struct{}{}
	}
	return out
}

func FilterStudyFieldIDs(ids []uint, allowed map[uint]struct{}) []uint {
	if len(ids) == 0 {
		return []uint{}
	}
	var out []uint
	for _, id := range ids {
		if _, ok := allowed[id]; ok {
			out = append(out, id)
		}
	}
	if out == nil {
		return []uint{}
	}
	return out
}

func normalizeLangList(langs []string) []string {
	if len(langs) == 0 {
		return []string{"uz", "ru", "en"}
	}
	return langs
}

func emptyLocalizedString(langs []string) map[string]string {
	out := make(map[string]string, len(langs))
	for _, lang := range langs {
		out[lang] = ""
	}
	return out
}

func emptyLocalizedTags(langs []string) TagNamesByLang {
	out := make(TagNamesByLang, len(langs))
	for _, lang := range langs {
		out[lang] = []string{}
	}
	return out
}

func ensureLocalizedString(in map[string]string, langs []string) map[string]string {
	out := emptyLocalizedString(langs)
	for _, lang := range langs {
		if in != nil {
			out[lang] = in[lang]
		}
	}
	return out
}

func ensureLocalizedTags(in TagNamesByLang, langs []string) TagNamesByLang {
	out := emptyLocalizedTags(langs)
	for _, lang := range langs {
		if in != nil && in[lang] != nil {
			out[lang] = in[lang]
		}
	}
	return out
}
