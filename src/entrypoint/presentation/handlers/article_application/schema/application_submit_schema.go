package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type ApplicationSubmitRequest struct {
	Name                        map[string]string     `json:"name" validate:"required"`
	CoAuthorsCount              uint                  `json:"co_authors_count" validate:"required"`
	CoAuthorsIDs                []uint                `json:"co_authors_ids"`
	StudyFieldsIDs              []uint                `json:"study_fields_ids"`
	LanguageID                  uint                  `json:"language_id" validate:"required"`
	ContentFile                 string                `json:"file" validate:"required"`
	Tags                        entity.TagNamesByLang `json:"tags"`
	DOI                         *string               `json:"doi"`
	ExpertConclusionFile        *string               `json:"expert_conclusion_file"`
	Annotation                  map[string]string     `json:"annotation" validate:"required"`
	ArticleAuthorAffiliationIDs []uint                `json:"article_author_affiliation_ids"`
	References                  []string              `json:"references"`
	PublishedDate               string                `json:"published_date"`
}

func (this *ApplicationSubmitRequest) ToEntity() *entity.ArticlePublishEntity {

	publishedDate, err := time.Parse("2006-01-02", this.PublishedDate)
	if err != nil {
		publishedDate = time.Now()
	}

	return entity.NewArticlePublishEntity(
		this.Name,
		this.CoAuthorsCount,
		this.CoAuthorsIDs,
		this.StudyFieldsIDs,
		this.LanguageID,
		this.ContentFile,
		this.Tags,
		this.DOI,
		this.ExpertConclusionFile,
		this.Annotation,
		this.ArticleAuthorAffiliationIDs,
		this.References,
		&publishedDate,
	)
}
