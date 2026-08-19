package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalIndexingToModel(journalId uint, indexes []*entity.JournalIndexingEntity) []*models.JournalIndexingModel {
	journalIndexingModels := make([]*models.JournalIndexingModel, len(indexes))
	for i, index := range indexes {
		journalIndexingModels[i] = models.NewJournalIndexingModel(index.IndexingType, index.URL, journalId)
	}
	return journalIndexingModels
}

func JournalIndexingModelToEntity(ji *models.JournalIndexingModel) *entity.JournalIndexingEntity {
	return entity.NewJournalIndexingEntity(
		ji.ID,
		ji.IndexingType,
		ji.URL,
		ji.JournalID,
		nil,
	)
}
