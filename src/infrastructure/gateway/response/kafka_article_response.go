package response

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type KafkaArticleResponse struct {
	ID              uint                       `json:"id"`
	Name            map[string]string          `json:"name"`
	PublicationDate string                     `json:"publication_date"`
	CoAuthorsCount  int                        `json:"co_authors_count"`
	CoAuthors       []*AuthorResponse          `json:"co_authors,omitempty"`
	AccessType      enum.AccessType            `json:"access_type"`
	StudyFields     []*KafkaStudyFieldResponse `json:"study_fields,omitempty"`
	LanguageID      uint                       `json:"language_id"`
	Language        *LanguageResponse          `json:"language,omitempty"`
	Annotation      map[string]string          `json:"annotation,omitempty"`
	Tags            entity.TagNamesByLang      `json:"tags,omitempty"`
	DOI             *string                    `json:"doi,omitempty"`
	ROI             *string                    `json:"roi,omitempty"`
	JournalID       uint                       `json:"journal_id"`
	Journal         *KafkaJournalResponse      `json:"journal,omitempty"`
	ContentFile     string                     `json:"content_file"`
}
