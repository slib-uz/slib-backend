package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type IntegrationAffiliationSchema struct {
	ScienceID        string `json:"science_id"`
	OrganizationTin  string `json:"organization_tin"`
	OrganizationName string `json:"organization_name"`
	PositionName     string `json:"position_name"`
}

type IntegrationArticleCreateSchema struct {
	Name                 map[string]string              `json:"name"`
	CoAuthorsCount       uint                           `json:"co_authors_count"`
	StudyFieldsIDs       []uint                         `json:"study_fields_ids"`
	LanguageID           uint                           `json:"language_id"`
	ContentFile          string                         `json:"content_file"`
	Tags                 entity.TagNamesByLang          `json:"tags"`
	DOI                  string                         `json:"doi"`
	ExpertConclusionFile string                         `json:"expert_conclusion_file"`
	Annotation           map[string]string              `json:"annotation"`
	Affiliations         []IntegrationAffiliationSchema `json:"affiliations"`
	References           []string                       `json:"references"`
	PublishedDate        string                         `json:"published_date"`
}

func (this *IntegrationArticleCreateSchema) ToEntity() *entity.ArticlePublishEntity {
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
		nil,
		this.StudyFieldsIDs,
		this.LanguageID,
		this.ContentFile,
		this.Tags,
		doi,
		expertConclusionFile,
		this.Annotation,
		nil,
		this.References,
		&publishedDate,
	)

	affiliations := make([]*entity.ArticleAuthorAffiliationEntity, 0, len(this.Affiliations))
	for _, aff := range this.Affiliations {
		affiliation := entity.NewArticleAuthorAffiliationEntity(
			0,
			nil,
			0,
			nil,
			aff.OrganizationName,
			aff.OrganizationTin,
			aff.PositionName,
		)
		affiliation.ScienceID = aff.ScienceID
		affiliations = append(affiliations, affiliation)
	}
	article.AuthorAffiliations = affiliations

	return article
}
