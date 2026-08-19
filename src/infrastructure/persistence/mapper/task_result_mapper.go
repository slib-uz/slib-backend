package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func TaskResultModelToEntity(model *models.TaskResultModel) *entity.TaskResultEntity {

	return entity.NewTaskResultEntity(
		model.TaskID,
		model.Task,
		model.Payload,
		model.Status,
		model.Error,
		model.RetryCount,
		&model.CreatedAt,
		&model.UpdatedAt,
		model.CompletedAt,
	)
}

func TaskResultEntityToModel(entity *entity.TaskResultEntity) models.TaskResultModel {
	return models.NewTaskResultModel(
		entity.TaskID,
		entity.Task,
		entity.Payload,
		entity.Status,
		entity.Error,
		entity.RetryCount,
		entity.CompletedAt,
	)
}
