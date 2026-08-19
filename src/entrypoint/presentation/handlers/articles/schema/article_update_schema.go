package schema

import (
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ArticleUpdateRequestSchema struct {
	Name           map[string]string `json:"name"`
	CoAuthorsCount int               `json:"co_authors_count" example:"3"`
	CoAuthorsIDs   []uint            `json:"co_authors_ids" example:"1,2,3"`
	StudyFieldsIDs []uint            `json:"study_fields_ids" example:"1,2"`
	LanguageID     uint              `json:"language_id" example:"1"`
	ContentFile    string            `json:"content_file" example:"temp/file.pdf"`
	Annotation     map[string]string `json:"annotation"`
	AccessType     enum.AccessType   `json:"access_type" example:"1" swaggertype:"integer"`

	DOI                  *string `json:"doi" example:"10.1234/example.doi"`
	ExpertConclusionFile *string `json:"expert_conclusion_file" example:"temp/expert_conclusion.pdf"`

	Pages      string `json:"pages" example:"1-15"`
	WebAddress string `json:"web_address" example:"https://example.com/article"`

	Tags            entity.TagNamesByLang `json:"tags"`
	AffiliationsIDs []uint                `json:"affiliations_ids" example:"10,11"`
	References      []string              `json:"references" example:"[1] Author A. Title"`
	PublishedDate   string                `json:"published_date" example:"2026-01-01"`

	UnconfirmedAuthors []*entity.UnconfirmedAuthorEntity `json:"unconfirmed_authors"`
}

func (this *ArticleUpdateRequestSchema) Validate() (bool, error) {
	if this.DOI != nil && *this.DOI == "" {
		return false, response.NewFailResponse(400, "DOI cannot be empty if provided")
	}
	return true, nil
}

func (this *ArticleUpdateRequestSchema) ToEntity() *entity.ArticleManageUpdateEntity {
	publishedDate, err := time.Parse("2006-01-02", this.PublishedDate)
	if err != nil {
		publishedDate = time.Now()
	}

	article := entity.NewArticleManageUpdateEntity(
		this.Name,
		this.CoAuthorsCount,
		this.CoAuthorsIDs,
		this.StudyFieldsIDs,
		this.LanguageID,
		this.ContentFile,
		this.Annotation,
		this.AccessType,
		this.DOI,
		this.ExpertConclusionFile,
		this.Pages,
		this.WebAddress,
		this.Tags,
		this.AffiliationsIDs,
		this.References,
		&publishedDate,
	)
	article.UnconfirmedAuthors = this.UnconfirmedAuthors
	return article
}
