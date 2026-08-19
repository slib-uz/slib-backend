package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ApplicationModelListToEntityList(models []*models.ArticleApplicationModel) []*entity.ApplicationEntity {
	entities := make([]*entity.ApplicationEntity, len(models))
	for i, model := range models {
		entities[i] = ApplicationModelToEntity(model)
	}
	return entities
}

func ApplicationModelToEntity(model *models.ArticleApplicationModel) *entity.ApplicationEntity {
	var article *entity.ArticleEntity
	if model.Article != nil {
		article = ArticleModelToEntity(model.Article)
	}

	var user *entity.UserEntity
	if model.User != nil {
		user = UserModelToEntity(model.User)
	}

	var journal *entity.JournalEntity
	if model.Journal != nil {
		journal = JournalModelToEntity(model.Journal)
	}

	var reviewStages []*entity.ReviewStageEntity
	if model.ReviewStages != nil {
		reviewStages = make([]*entity.ReviewStageEntity, len(model.ReviewStages))
		for i, stage := range model.ReviewStages {
			reviewStages[i] = ReviewStageModelToEntity(stage)
		}
	}

	return entity.NewApplicationEntity(
		model.ID,
		model.Number,
		model.ArticleID,
		article,
		model.JournalID,
		journal,
		model.UserID,
		user,
		model.IsPublished,
		reviewStages,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
}
