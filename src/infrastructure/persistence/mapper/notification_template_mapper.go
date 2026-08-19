package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func NotificationTemplateModelToEntity(data *models.NotificationTemplateModel) *entity.NotificationTemplateEntity {
	return entity.NewNotificationTemplateEntity(
		data.ID,
		data.Key,
		FromGormJson[map[string]string](data.Title),
		FromGormJson[map[string]string](data.Body),
	)
}
