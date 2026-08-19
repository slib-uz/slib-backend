package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func InstitutionModelToEntity(data *models.InstitutionModel) *entity.InstitutionEntity {
	return entity.NewInstitutionEntity(data.ID, data.Name, data.Tin, data.Logo)
}

func InstitutionModelToDetailEntity(data *models.InstitutionModel) *entity.InstitutionEntity {
	institution := InstitutionModelToEntity(data)
	institution.Publishers = PublisherModelsToEntities(data.Publishers)
	return institution
}

func PublisherModelsToEntities(publishers []models.PublisherModel) []*entity.PublisherEntity {
	result := make([]*entity.PublisherEntity, len(publishers))
	for i := range publishers {
		result[i] = PublisherModelToEntity(&publishers[i])
	}
	return result
}

func InstitutionEntityToModel(data *entity.InstitutionEntity) *models.InstitutionModel {
	return models.NewInstitutionModel(data.Name, data.Tin, data.Logo)
}
