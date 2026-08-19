package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func SupportDialogListModelToEntity(data []*models.SupportDialogModel) []*entity.SupportDialogEntity {
	var result = make([]*entity.SupportDialogEntity, len(data))

	for i, item := range data {
		result[i] = SupportDialogModelToEntity(item)
	}

	return result
}

func SupportDialogModelToEntity(data *models.SupportDialogModel) *entity.SupportDialogEntity {
	return entity.NewSupportDialogEntity(data.ID, data.MessageType, data.OwnerID, SupportDialogOwnerModelToEntity(data.Owner), data.Message, data.ChatID, data.IsRead, data.CreatedAt)
}

func SupportDialogEntityToModel(data *entity.SupportDialogEntity) *models.SupportDialogModel {
	return models.NewSupportDialogModel(data.MessageType, data.OwnerID, data.Message, data.ChatID, data.IsRead)
}
