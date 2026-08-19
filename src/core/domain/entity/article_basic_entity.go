package entity

import "time"

type ArticleBasicEntity struct {
	ID                   uint                `json:"id"`
	Name                 map[string]string   `json:"name"`
	PublicationDate      time.Time           `json:"publication_date"`
	CoAuthors            []*AuthorEntity     `json:"co_authors"`
	StudyFields          []*StudyFieldEntity `json:"study_fields"`
	ContentFile          string              `json:"content_file"`
	JournalID            uint                `json:"journal_id"`
	Journal              *JournalEntity      `json:"journal"`
	DOI                  *string             `json:"doi"`
	ROI                  *string             `json:"roi"`
	RatingCount          int64               `json:"rating_count"`
	RatingSum            int64               `json:"rating_sum"`
	RatingAvg            float64             `json:"rating_avg"`
	ExpertConclusionFile *string             `json:"expert_conclusion_file"`
}

func NewArticleBasicEntity(ID uint, name map[string]string, publicationDate time.Time, coAuthors []*AuthorEntity, studyFields []*StudyFieldEntity, contentFile string, journalID uint, journal *JournalEntity, DOI *string, ROI *string, ratingCount int64, ratingSum int64, ratingAvg float64, expertConclusionFile *string) *ArticleBasicEntity {
	return &ArticleBasicEntity{ID: ID, Name: name, PublicationDate: publicationDate, CoAuthors: coAuthors, StudyFields: studyFields, ContentFile: contentFile, JournalID: journalID, Journal: journal, DOI: DOI, ROI: ROI, RatingCount: ratingCount, RatingSum: ratingSum, RatingAvg: ratingAvg, ExpertConclusionFile: expertConclusionFile}
}
