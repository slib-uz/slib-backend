package mapper

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
	"slib.uz/src/infrastructure/gateway/response/payload"
)

func ArticleEntityToPayload(data *entity.ArticleEntity, uploadedPath *string, sourceLink string) *payload.PublishArticlePayload {
	codes := make([]*uint, len(data.StudyFields))
	for i, field := range data.StudyFields {
		codes[i] = field.Code
	}

	return payload.NewPublishArticlePayload(
		data.Name,
		data.Annotation,
		data.PublicationDate.Format(time.DateOnly),
		data.CoAuthorsCount,
		data.CoAuthors,
		data.AccessType,
		codes,
		data.LanguageID,
		data.DOI,
		data.Journal.ISSNPaper,
		data.Journal.ISSNOnline,
		uploadedPath,
		sourceLink,
	)
}

func RoiBodyToEntity(data *response.ROIResponse) *entity.RoiSimilarityEntity {
	return entity.NewRoiSimilarityEntity(
		data.Data.Name,
		data.Data.ROI,
		data.Data.Percent,
		false,
	)
}

func ArticleSearchEntityToPayload(data *entity.ArticleEntity) *payload.PublishArticlePayload {
	codes := make([]*uint, len(data.StudyFields))
	for i, field := range data.StudyFields {
		codes[i] = field.Code
	}

	return payload.NewPublishArticlePayload(
		data.Name,
		data.Annotation,
		data.PublicationDate.Format(time.DateOnly),
		data.CoAuthorsCount,
		data.CoAuthors,
		data.AccessType,
		codes,
		data.LanguageID,
		data.DOI,
		data.Journal.ISSNPaper,
		data.Journal.ISSNOnline,
		nil,
		"",
	)
}

func ROIResponseToEntity(data *entity.ArticleEntity, bodyData *response.PublishArticleResponse) *entity.ArticleResponseEntity {
	authors := make([]*entity.UserSharedEntity, len(data.CoAuthors))
	for i, author := range data.CoAuthors {
		authors[i] = entity.NewUserSharedEntity(
			author.ID,
			author.ScienceID,
			author.FullName,
		)
	}
	return entity.NewArticleResponseEntity(
		bodyData.ID,
		data.Name,
		data.PublicationDate.Format(time.DateOnly),
		authors,
		data.StudyFields,
		data.DOI,
		bodyData.ROI,
		bodyData.Percent,
	)
}

func ROISuccessResToEntity(data *response.ROIDetailResponse) *entity.ArticleResponseEntity {
	coAuthors := make([]*entity.UserSharedEntity, len(data.CoAuthors))
	for i, author := range data.CoAuthors {
		coAuthors[i] = entity.NewUserSharedEntity(
			author.ID,
			author.ScienceID,
			author.FullName,
		)
	}
	studyFields := make([]*entity.StudyFieldEntity, len(data.StudyFields))
	for i, field := range data.StudyFields {
		studyFields[i] = entity.NewStudyFieldEntity(
			field.ID,
			field.Name,
			nil,
			field.Code,
			nil,
		)
	}

	return entity.NewArticleResponseEntity(
		data.ID,
		data.Name,
		data.PublicationDate,
		coAuthors,
		studyFields,
		data.DOI,
		data.ROI,
		data.Percent,
	)
}
