package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type ArticleManageUpdateEntity struct {
	Name           map[string]string `json:"name,omitempty"`
	CoAuthorsCount int               `json:"co_authors_count,omitempty"`
	CoAuthorsIDs   []uint            `json:"co_authors_ids,omitempty"`
	StudyFieldsIDs []uint            `json:"study_fields_ids,omitempty"`
	LanguageID     uint              `json:"language_id,omitempty"`
	ContentFile    string            `json:"content_file,omitempty"`
	Annotation     map[string]string `json:"annotation,omitempty"`
	AccessType     enum.AccessType   `json:"access_type,omitempty"`

	DOI                  *string `json:"doi,omitempty"`
	ExpertConclusionFile *string `json:"expert_conclusion_file,omitempty"`

	Pages      string `json:"pages,omitempty"`
	WebAddress string `json:"web_address,omitempty"`

	Tags            TagNamesByLang `json:"tags,omitempty"`
	AffiliationsIDs []uint         `json:"affiliations_ids,omitempty"`
	References      []string       `json:"references,omitempty"`
	PublishedDate   *time.Time     `json:"published_date,omitempty"`

	UnconfirmedAuthors []*UnconfirmedAuthorEntity `json:"unconfirmed_authors,omitempty"`
}

func NewArticleManageUpdateEntity(
	name map[string]string,
	coAuthorsCount int,
	coAuthorsIDs []uint,
	studyFieldsIDs []uint,
	languageID uint,
	contentFile string,
	annotation map[string]string,
	accessType enum.AccessType,
	DOI *string,
	expertConclusionFile *string,
	pages string,
	webAddress string,
	tags TagNamesByLang,
	affiliationsIDs []uint,
	references []string,
	publishedDate *time.Time,
) *ArticleManageUpdateEntity {
	return &ArticleManageUpdateEntity{
		Name:                 name,
		CoAuthorsCount:       coAuthorsCount,
		CoAuthorsIDs:         coAuthorsIDs,
		StudyFieldsIDs:       studyFieldsIDs,
		LanguageID:           languageID,
		ContentFile:          contentFile,
		Annotation:           annotation,
		AccessType:           accessType,
		DOI:                  DOI,
		ExpertConclusionFile: expertConclusionFile,
		Pages:                pages,
		WebAddress:           webAddress,
		Tags:                 tags,
		AffiliationsIDs:      affiliationsIDs,
		References:           references,
		PublishedDate:        publishedDate,
	}
}
