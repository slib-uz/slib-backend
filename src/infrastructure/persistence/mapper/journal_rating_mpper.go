package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalRatingModelListToEntityList(models []*models.JournalRatingModel) []*entity2.JournalRatingEntity {
	var result = make([]*entity2.JournalRatingEntity, len(models))
	for i, item := range models {
		result[i] = JournalRatingModelToEntity(item)
	}
	return result
}

func JournalRatingEntityListToModelList(entities []*entity2.JournalRatingEntity) []*models.JournalRatingModel {
	var result = make([]*models.JournalRatingModel, len(entities))
	for i, item := range entities {
		result[i] = JournalRatingEntityToModel(item)
	}
	return result
}

func JournalRatingModelToEntity(model *models.JournalRatingModel) *entity2.JournalRatingEntity {
	var user *entity2.JournalRaterEntity
	if model.User != nil {
		user = JournalRaterModelToEntity(model.User)
	}
	return entity2.NewJournalRatingEntity(
		model.ID,
		model.UserID,
		user,
		model.JournalID,
		model.Stars,
		model.Review,
		model.CreatedAt,
	)
}

func JournalRatingEntityToModel(entity *entity2.JournalRatingEntity) *models.JournalRatingModel {
	return models.NewJournalRatingModel(
		entity.UserID,
		entity.JournalID,
		entity.Stars,
		entity.Review,
	)
}
