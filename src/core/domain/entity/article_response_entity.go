package entity

type ArticleResponseEntity struct {
	ID              uint                `json:"id"`
	Name            map[string]string   `json:"name"`
	PublicationDate string              `json:"publication_date"`
	CoAuthors       []*UserSharedEntity `json:"co_authors"`
	StudyFields     []*StudyFieldEntity `json:"study_fields"`
	DOI             *string             `json:"doi,omitempty"`
	ROI             *string             `json:"roi,omitempty"`
	Percent         float64             `json:"percent"`
}

func NewArticleResponseEntity(ID uint, name map[string]string, publicationDate string, coAuthors []*UserSharedEntity, studyFields []*StudyFieldEntity, DOI *string, ROI *string, percent float64) *ArticleResponseEntity {
	return &ArticleResponseEntity{ID: ID, Name: name, PublicationDate: publicationDate, CoAuthors: coAuthors, StudyFields: studyFields, DOI: DOI, ROI: ROI, Percent: percent}
}
