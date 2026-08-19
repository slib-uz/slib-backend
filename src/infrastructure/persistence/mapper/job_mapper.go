package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JobModelToEntity(model *models.JobModel) *entity.JobEntity {
	entity := entity.NewJobEntity(
		model.OrganizationID,
		model.UserID,
		model.OrganizationTin,
		model.OrganizationName,
		model.PositionName,
	)
	entity.ID = model.ID
	return entity
}

func JobEntityToModel(userID uint, entity *entity.JobEntity) *models.JobModel {
	return models.NewJobModel(
		entity.OrganizationID,
		userID,
		entity.OrganizationTin,
		entity.OrganizationName,
		entity.PositionName,
	)
}

func JobUpdateMapper(existing, job *models.JobModel) *models.JobModel {
	existing.OrganizationID = job.OrganizationID
	existing.UserID = job.UserID
	existing.OrganizationTin = job.OrganizationTin
	existing.PositionName = job.PositionName

	return existing
}

func JobModelListToEntity(models []*models.JobModel) []*entity.JobEntity {
	entities := make([]*entity.JobEntity, len(models))
	for i, model := range models {
		entities[i] = JobModelToEntity(model)
	}
	return entities
}
