package schema

import (
	"slib.uz/src/core/domain/entity"
)

type ResubmitForPeerReviewRequest struct {
	SubmissionID uint   `json:"submission_id" validate:"required"`
	File         string `json:"file" validate:"required"`
}

type ArticleUpdateRequest struct {
	Name           *map[string]string `json:"name"`
	CoAuthorsCount *int               `json:"co_authors_count"`
	CoAuthorsIDs   *[]uint            `json:"co_authors_ids"`
	StudyFieldsIDs *[]uint            `json:"study_fields_ids"`
	LanguageID     *uint              `json:"language_id"`
	ContentFile    *string            `json:"file"`
	Annotation     *map[string]string `json:"annotation"`

	DOI      *string `json:"doi"`
	DOIClear *bool   `json:"doi_clear"`

	ExpertConclusionFile  *string `json:"expert_conclusion_file"`
	ExpertConclusionClear *bool   `json:"expert_conclusion_file_clear"`

	AffiliationsIDs *[]uint `json:"article_author_affiliation_ids"`

	Tags entity.TagNamesByLang `json:"tags"`
}

func (r *ArticleUpdateRequest) ToEntity() *entity.ArticleUpdateEntity {
	e := entity.NewArticleUpdateEntity(
		r.Name,
		r.CoAuthorsCount,
		r.CoAuthorsIDs,
		r.StudyFieldsIDs,
		r.LanguageID,
		r.ContentFile,
		r.Annotation,
		r.DOI,
		r.DOIClear,
		r.ExpertConclusionFile,
		r.ExpertConclusionClear,
		r.AffiliationsIDs,
	)
	e.Tags = r.Tags
	return e
}
