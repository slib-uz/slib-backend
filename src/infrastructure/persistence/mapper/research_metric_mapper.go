package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ResearchMetricListModelToEntity(models []*models.ResearchMetricModel) []*entity.ResearchMetricEntity {
	entities := make([]*entity.ResearchMetricEntity, len(models))
	for i, model := range models {
		entities[i] = ResearchMetricModelToEntity(model)
	}
	return entities
}

func ResearchMetricEntityToModel(entity *entity.ResearchMetricEntity) *models.ResearchMetricModel {
	return models.NewResearchMetricModel(entity.UserID, entity.ProfileUrl, entity.HIndex, entity.Source)
}

func ResearchMetricModelToEntity(model *models.ResearchMetricModel) *entity.ResearchMetricEntity {
	return entity.NewResearchMetricEntity(model.ID, model.UserID, model.ProfileUrl, model.HIndex, model.Source)
}
