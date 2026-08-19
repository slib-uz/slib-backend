package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ArticleCreateSchema struct {
	Name                        map[string]string                 `json:"name"`
	CoAuthorsCount              uint                              `json:"co_authors_count"`
	CoAuthorsIDs                []uint                            `json:"co_authors_ids"`
	StudyFieldsIDs              []uint                            `json:"study_fields_ids"`
	LanguageID                  uint                              `json:"language_id"`
	ContentFile                 string                            `json:"content_file"`
	Tags                        entity.TagNamesByLang             `json:"tags"`
	DOI                         string                            `json:"doi"`
	ExpertConclusionFile        string                            `json:"expert_conclusion_file"`
	Annotation                  map[string]string                 `json:"annotation"`
	ArticleAuthorAffiliationIDs []uint                            `json:"article_author_affiliation_ids"`
	References                  []string                          `json:"references"`
	PublishedDate               string                            `json:"published_date"`
	UnconfirmedAuthors          []*entity.UnconfirmedAuthorEntity `json:"unconfirmed_authors"`
	AccessType                  enum.AccessType                   `json:"access_type"`
}

func (this *ArticleCreateSchema) ToEntity() *entity.ArticlePublishEntity {
	var doi *string
	if this.DOI != "" {
		doi = &this.DOI
	}

	var expertConclusionFile *string
	if this.ExpertConclusionFile != "" {
		expertConclusionFile = &this.ExpertConclusionFile
	}

	publishedDate, err := time.Parse("2006-01-02", this.PublishedDate)
	if err != nil {
		publishedDate = time.Now()
	}

	article := entity.NewArticlePublishEntity(
		this.Name,
		this.CoAuthorsCount,
		this.CoAuthorsIDs,
		this.StudyFieldsIDs,
		this.LanguageID,
		this.ContentFile,
		this.Tags,
		doi,
		expertConclusionFile,
		this.Annotation,
		this.ArticleAuthorAffiliationIDs,
		this.References,
		&publishedDate,
	)
	article.UnconfirmedAuthors = this.UnconfirmedAuthors
	article.AccessType = this.AccessType
	return article
}
