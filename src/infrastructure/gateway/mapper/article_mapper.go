package mapper

import (
	"encoding/json"
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/gateway/response"
)

func ArticleToMessagingMapper(source *entity.ArticleEntity, file string) ([]byte, error) {
	kafkaDTO := response.KafkaArticleResponse{
		ID:              source.ID,
		Name:            source.Name,
		PublicationDate: source.PublicationDate.Format(time.RFC3339),
		CoAuthorsCount:  source.CoAuthorsCount,
		AccessType:      source.AccessType,
		Annotation:      source.Annotation,
		DOI:             source.DOI,
		ROI:             source.ROI,
		JournalID:       source.JournalID,
		ContentFile:     file,
	}

	if source.LanguageID != 0 && source.Language != nil {
		kafkaDTO.Language = LanguageEntityToResponse(source.Language)
	}

	if len(source.CoAuthors) > 0 {
		kafkaDTO.CoAuthors = make([]*response.AuthorResponse, len(source.CoAuthors))
		for i, author := range source.CoAuthors {
			kafkaDTO.CoAuthors[i] = AuthorEntityToRes(author)
		}
	}

	if len(source.StudyFields) > 0 {
		kafkaDTO.StudyFields = make([]*response.KafkaStudyFieldResponse, len(source.StudyFields))
		for i, field := range source.StudyFields {
			kafkaDTO.StudyFields[i] = StudyFieldToResponse(field)
		}
	}

	if len(source.Tags) > 0 {
		kafkaDTO.Tags = source.Tags
	}

	if source.Journal != nil {
		kafkaDTO.Journal = JournalEntityToResponse(source.Journal)
	}

	return json.Marshal(kafkaDTO)
}

func ROIPublishEntityToArticleEntity(source *entity.ROIPublishEntity) *entity.ArticleEntity {
	publicationDate, _ := time.Parse("2006-01-02 15:04:05", source.PublicationDate)
	return entity.NewArticleEntity(
		0,
		source.Name,
		publicationDate,
		source.CoAuthorsCount,
		nil,
		enum.AccessType(source.AccessType),
		nil,
		0,
		nil,
		source.Annotation,
		"",
		nil,
		source.DOI,
		source.ROI,
		0,
		nil,
		nil,
		0,
		0,
		0,
		0,
		nil,
	)
}
