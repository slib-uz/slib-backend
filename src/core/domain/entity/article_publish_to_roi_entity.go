package entity

type ROIPublishEntity struct {
	Name            map[string]string       `json:"name"`
	PublicationDate string                  `json:"publication_date"`
	CoAuthorsCount  int                     `json:"co_authors_count"`
	AccessType      int                     `json:"access_type"`
	CoAuthors       []*ROICoAuthorEntity    `json:"co_authors"`
	StudyFields     []*ROIStudyFieldsEntity `json:"study_fields"`
	LanguageID      uint                    `json:"language_id"`
	Annotation      map[string]string       `json:"annotation"`
	DOI             *string                 `json:"doi"`
	ROI             *string                 `json:"roi"`
	URL             string                  `json:"url"`
	Journal         *ROIJournalEntity       `json:"journal"`
	Publisher       *ROIPublisherEntity     `json:"publisher"`
	Citations       interface{}             `json:"citations"`
}

func NewROIPublishEntity(name map[string]string, publicationDate string, coAuthorsCount int, accessType int, coAuthors []*ROICoAuthorEntity, studyFields []*ROIStudyFieldsEntity, languageID uint, annotation map[string]string, doi *string, roi *string, url string, journal *ROIJournalEntity, publisher *ROIPublisherEntity, citations interface{}) *ROIPublishEntity {
	return &ROIPublishEntity{
		Name:            name,
		PublicationDate: publicationDate,
		CoAuthorsCount:  coAuthorsCount,
		AccessType:      accessType,
		CoAuthors:       coAuthors,
		StudyFields:     studyFields,
		LanguageID:      languageID,
		Annotation:      annotation,
		DOI:             doi,
		ROI:             roi,
		URL:             url,
		Journal:         journal,
		Publisher:       publisher,
		Citations:       citations,
	}
}

type ROICoAuthorEntity struct {
	FullName  string `json:"full_name"`
	ScienceID string `json:"science_id"`
}

func NewROICoAuthorEntity(fullName string, scienceID string) *ROICoAuthorEntity {
	return &ROICoAuthorEntity{
		FullName:  fullName,
		ScienceID: scienceID,
	}
}

type ROIStudyFieldsEntity struct {
	Name map[string]string `json:"name"`
	Code uint              `json:"code"`
}

func NewROIStudyFieldsEntity(name map[string]string, code uint) *ROIStudyFieldsEntity {
	return &ROIStudyFieldsEntity{
		Name: name,
		Code: code,
	}
}

type ROIJournalEntity struct {
	Name            map[string]string `json:"name"`
	ShortName       map[string]string `json:"short_name"`
	EstablishedDate string            `json:"established_date"`
	ISSNPaper       string            `json:"issn_paper"`
	ISSNOnline      string            `json:"issn_online"`
	PublisherID     uint              `json:"publisher_id"`
}

func NewROIJournalEntity(name map[string]string, shortName map[string]string, establishedDate string, issnPaper string, issnOnline string, publisherID uint) *ROIJournalEntity {
	return &ROIJournalEntity{
		Name:            name,
		ShortName:       shortName,
		EstablishedDate: establishedDate,
		ISSNPaper:       issnPaper,
		ISSNOnline:      issnOnline,
		PublisherID:     publisherID,
	}
}

type ROIPublisherEntity struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	TIN       string `json:"tin"`
}

func NewROIPublisherEntity(name string, shortName string, tin string) *ROIPublisherEntity {
	return &ROIPublisherEntity{
		Name:      name,
		ShortName: shortName,
		TIN:       tin,
	}
}
