package entity

type ROIArticleUpdateEntity struct {
	ROI              string               `json:"roi"`
	Name             map[string]string    `json:"name"`
	PublicationDate  string               `json:"publication_date"`
	CoAuthorsCount   uint                 `json:"co_authors_count"`
	CoAuthors        []*ROICoAuthorEntity `json:"co_authors"`
	AccessType       int                  `json:"access_type"`
	StudyFieldsCodes []uint               `json:"study_fields_codes"`
	LanguageID       uint                 `json:"language_id"`
	File             string               `json:"file"`
	DOI              *string              `json:"doi"`
	Annotation       map[string]string    `json:"annotation"`
	SourceLink       string               `json:"source_link"`
}

func NewROIArticleUpdateEntity(
	roi string,
	name map[string]string,
	publicationDate string,
	coAuthorsCount uint,
	coAuthors []*ROICoAuthorEntity,
	accessType int,
	studyFieldsCodes []uint,
	languageID uint,
	file string,
	doi *string,
	annotation map[string]string,
	sourceLink string,
) *ROIArticleUpdateEntity {
	return &ROIArticleUpdateEntity{
		ROI:              roi,
		Name:             name,
		PublicationDate:  publicationDate,
		CoAuthorsCount:   coAuthorsCount,
		CoAuthors:        coAuthors,
		AccessType:       accessType,
		StudyFieldsCodes: studyFieldsCodes,
		LanguageID:       languageID,
		File:             file,
		DOI:              doi,
		Annotation:       annotation,
		SourceLink:       sourceLink,
	}
}
