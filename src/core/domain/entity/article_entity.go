package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type ArticleEntity struct {
	ID                       uint                              `json:"id"`
	Name                     map[string]string                 `json:"name"`
	PublicationDate          time.Time                         `json:"publication_date"`
	CoAuthorsCount           int                               `json:"co_authors_count"`
	CoAuthors                []*AuthorEntity                   `json:"co_authors"`
	AccessType               enum.AccessType                   `json:"access_type"`
	StudyFields              []*StudyFieldEntity               `json:"study_fields"`
	LanguageID               uint                              `json:"language_id"`
	Language                 *LanguageEntity                   `json:"language"`
	Annotation               map[string]string                 `json:"annotation"`
	ContentFile              string                            `json:"content_file"`
	Tags                     TagNamesByLang                    `json:"tags"`
	DOI                      *string                           `json:"doi"`
	ROI                      *string                           `json:"roi"`
	JournalID                uint                              `json:"journal_id"`
	Journal                  *JournalEntity                    `json:"journal"`
	ExpertConclusionFile     *string                           `json:"expert_conclusion_file"`
	ViewsCount               int64                             `json:"views_count"`
	RatingCount              int64                             `json:"rating_count"`
	RatingSum                int64                             `json:"rating_sum"`
	RatingAvg                float64                           `json:"rating_avg"`
	ArticleAuthorAffiliation []*ArticleAuthorAffiliationEntity `json:"article_author_affiliation"`
	References               []*ReferenceEntity                `json:"references"`
	UnconfirmedAuthors       []*UnconfirmedAuthorEntity        `json:"unconfirmed_authors"`
}

func NewArticleEntity(ID uint, name map[string]string, publicationDate time.Time, coAuthorsCount int, coAuthors []*AuthorEntity, accessType enum.AccessType, studyFields []*StudyFieldEntity, languageID uint, language *LanguageEntity, annotation map[string]string, contentFile string, tags TagNamesByLang, DOI *string, ROI *string, journalID uint, journal *JournalEntity, expertConclusionFile *string, viewsCount int64, ratingCount int64, ratingSum int64, ratingAvg float64, articleAuthorAffiliation []*ArticleAuthorAffiliationEntity) *ArticleEntity {
	return &ArticleEntity{ID: ID, Name: name, PublicationDate: publicationDate, CoAuthorsCount: coAuthorsCount, CoAuthors: coAuthors, AccessType: accessType, StudyFields: studyFields, LanguageID: languageID, Language: language, Annotation: annotation, ContentFile: contentFile, Tags: tags, DOI: DOI, ROI: ROI, JournalID: journalID, Journal: journal, ExpertConclusionFile: expertConclusionFile, ViewsCount: viewsCount, RatingCount: ratingCount, RatingSum: ratingSum, RatingAvg: ratingAvg, ArticleAuthorAffiliation: articleAuthorAffiliation}
}
