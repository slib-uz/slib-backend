package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func NotificationSendResultEntityToModel(data *entity.NotificationSendResultEntity) *models.NotificationSendResultModel {
	return models.NewNotificationSendResultModel(
		data.NotificationID,
		data.SuccessCount,
		data.FailureCount,
		ToGormJson(data.FailedTokens),
		ToGormJson(data.Errors),
	)
}
