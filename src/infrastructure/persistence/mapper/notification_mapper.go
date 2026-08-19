package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func NotificationEntityListToModelList(data []*entity.NotificationEntity) []*models.NotificationModel {
	var modelsList []*models.NotificationModel
	for _, item := range data {
		modelsList = append(modelsList, NotificationEntityToModel(item))
	}
	return modelsList
}

func NotificationEntityToModel(data *entity.NotificationEntity) *models.NotificationModel {
	return models.NewNotificationModel(
		data.UserID,
		ToGormJson(data.Title),
		ToGormJson(data.Body),
		data.Topic,
		ToGormJson(data.ExtraData),
		data.IsEmail,
		data.IsSms,
		data.IsBroadcast,
	)
}

func NotificationModelToEntity(data *models.NotificationModel) *entity.NotificationEntity {
	return entity.NewNotificationEntity(
		data.ID,
		data.UserID,
		FromGormJson[map[string]string](data.Title),
		FromGormJson[map[string]string](data.Body),
		data.Topic,
		FromGormJson[map[string]string](data.ExtraData),
		data.IsEmail,
		data.IsSms,
		data.IsBroadcast,
		data.IsRead,
		data.CreatedAt,
	)
}
