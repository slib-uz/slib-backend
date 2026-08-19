package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func AntiPlagResulModelListToEntityList(data []*models.AntiPlagResultModel) []*entity2.AntiPlagResultEntity {
	var result = make([]*entity2.AntiPlagResultEntity, len(data))

	for i, item := range data {
		result[i] = AntiPlagResultModelToEntity(item)
	}

	return result
}

func AntiPlagResultEntityToModel(data *entity2.AntiPlagResultEntity) *models.AntiPlagResultModel {

	return models.NewAntiPlagResultModel(
		data.ID,
		data.ReviewStageID,
		data.ApplicationID,
		data.ArticleID,
		data.JournalID,
		data.ExternalID,
		data.Status,
		data.StatusDisplay,
		data.PlagiarismPercent,
		data.LegalPercent,
		data.SelfCitePercent,
		data.UnknownPercent,
		data.ShortReportURL,
		data.FullReportURL,
		data.ExternalCreatedAt,
	)
}

func AntiPlagResultModelToEntity(data *models.AntiPlagResultModel) *entity2.AntiPlagResultEntity {

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

	return entity2.NewAntiPlagResultEntity(
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
		data.Status,
		data.StatusDisplay,
		data.PlagiarismPercent,
		data.LegalPercent,
		data.SelfCitePercent,
		data.UnknownPercent,
		data.ShortReportURL,
		data.FullReportURL,
		data.ExternalCreatedAt,
		&data.CreatedAt,
		data.Certificate,
	)
}
