package payload

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type PublishArticlePayload struct {
	Name            map[string]string      `json:"name"`
	Annotation      map[string]string      `json:"annotation,omitempty"`
	PublicationDate string                 `json:"publication_date,omitempty"`
	CoAuthorsCount  int                    `json:"co_authors_count"`
	CoAuthors       []*entity.AuthorEntity `json:"co_authors"`
	AccessType      enum.AccessType        `json:"access_type"`
	StudyFieldCodes []*uint                `json:"study_fields_codes"`
	LanguageID      uint                   `json:"language_id"`
	DOI             *string                `json:"doi,omitempty"`
	ISSNPaper       *string                `json:"issn_paper,omitempty"`
	ISSNOnline      *string                `json:"issn_online,omitempty"`
	File            *string                `json:"file,omitempty"`
	SourceLink      string                 `json:"source_link"`
}

func NewPublishArticlePayload(
	name map[string]string,
	annotation map[string]string,
	publicationDate string,
	coAuthorsCount int,
	coAuthors []*entity.AuthorEntity,
	accessType enum.AccessType,
	studyFieldCodes []*uint,
	languageID uint,
	DOI *string,
	ISSNPaper *string,
	ISSNOnline *string,
	file *string,
	sourceLink string,
) *PublishArticlePayload {
	return &PublishArticlePayload{
		Name:            name,
		Annotation:      annotation,
		PublicationDate: publicationDate,
		CoAuthorsCount:  coAuthorsCount,
		CoAuthors:       coAuthors,
		AccessType:      accessType,
		StudyFieldCodes: studyFieldCodes,
		LanguageID:      languageID,
		DOI:             DOI,
		ISSNPaper:       ISSNPaper,
		ISSNOnline:      ISSNOnline,
		File:            file,
		SourceLink:      sourceLink,
	}
}
