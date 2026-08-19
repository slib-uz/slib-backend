package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type ArticlePublishEntity struct {
	Name                        map[string]string
	CoAuthorsCount              uint
	CoAuthorsIDs                []uint
	StudyFieldsIDs              []uint
	LanguageID                  uint
	ContentFile                 string
	Tags                        TagNamesByLang
	DOI                         *string
	ExpertConclusionFile        *string
	Annotation                  map[string]string
	ArticleAuthorAffiliationIDs []uint
	AuthorAffiliations          []*ArticleAuthorAffiliationEntity
	References                  []string
	PublishedDate               *time.Time
	UnconfirmedAuthors          []*UnconfirmedAuthorEntity
	AccessType                  enum.AccessType
}

func NewArticlePublishEntity(
	name map[string]string,
	coAuthorsCount uint,
	coAuthorsIDs []uint,
	studyFieldsIDs []uint,
	languageID uint,
	file string,
	tags TagNamesByLang,
	DOI *string,
	expertConclusionFile *string,
	annotation map[string]string,
	articleAuthorAffiliationIDs []uint,
	references []string,
	publishedDate *time.Time,
) *ArticlePublishEntity {
	return &ArticlePublishEntity{
		Name:                        name,
		CoAuthorsCount:              coAuthorsCount,
		CoAuthorsIDs:                coAuthorsIDs,
		StudyFieldsIDs:              studyFieldsIDs,
		LanguageID:                  languageID,
		ContentFile:                 file,
		Tags:                        tags,
		DOI:                         DOI,
		ExpertConclusionFile:        expertConclusionFile,
		Annotation:                  annotation,
		ArticleAuthorAffiliationIDs: articleAuthorAffiliationIDs,
		References:                  references,
		PublishedDate:               publishedDate,
	}
}
