package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func AiDetectResultEntityToModel(data *entity2.AiDetectResultEntity) *models.AiDetectResultModel {
	return models.NewAiDetectResultModel(
		data.ID,
		data.ReviewStageID,
		data.ApplicationID,
		data.ArticleID,
		data.JournalID,
		data.ExternalID,
		data.WordsCount,
		data.Status,
		data.StatusDisplay,
		data.HumanPercent,
		data.WarnPercent,
		data.AiPercent,
		data.ReportURL,
		data.ExternalCreatedAt,
	)
}

func AiDetectResultModelToEntity(data *models.AiDetectResultModel) *entity2.AiDetectResultEntity {
	var application *entity2.ApplicationEntity
	var journal *entity2.JournalEntity
	var reviewStage *entity2.ReviewStageEntity
	var article *entity2.ArticleEntity

	if data.Application != nil {
		application = ApplicationModelToEntity(data.Application)
	}
	if data.Journal != nil {
		journal = JournalModelToEntity(data.Journal)
	}
	if data.Article != nil {
		article = ArticleModelToEntity(data.Article)
	}
	if data.ReviewStage != nil {
		reviewStage = ReviewStageModelToEntity(data.ReviewStage)
	}

	return entity2.NewAiDetectResultEntity(
		data.ID,
		data.ReviewStageID,
		reviewStage,
		data.ApplicationID,
		application,
		data.ArticleID,
		article,
		data.JournalID,
		journal,
		data.ExternalID,
		data.WordsCount,
		data.Status,
		data.StatusDisplay,
		data.HumanPercent,
		data.WarnPercent,
		data.AiPercent,
		data.ReportURL,
		data.ExternalCreatedAt,
		&data.CreatedAt,
	)
}
