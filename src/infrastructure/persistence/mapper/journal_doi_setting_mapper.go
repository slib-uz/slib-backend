package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalDoiSettingEntityToModel(it *entity.JournalDoiSettingEntity) *models.JournalDoiSettingModel {
	return models.NewJournalDoiSettingModel(it.JournalID, it.JournalName, it.Username, it.Password, it.DOIPrefix, it.DOISuffix)
}

func JournalDoiSettingModelToEntity(it *models.JournalDoiSettingModel) *entity.JournalDoiSettingEntity {
	return entity.NewJournalDoiSettingEntity(it.ID, it.JournalID, it.JournalName, it.Username, it.Password, it.DOIPrefix, it.DOISuffix, it.CreatedAt, it.UpdatedAt)
}
