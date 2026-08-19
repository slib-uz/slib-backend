package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func DistrictModelToEntity(m *models.DistrictModel) *entity.DistrictEntity {
	return entity.NewDistrictEntity(m.ID, jsonNameToMap(m.Name), m.Soato, m.RegionID)
}

func DistrictModelsToEntities(ms []models.DistrictModel) []*entity.DistrictEntity {
	entities := make([]*entity.DistrictEntity, len(ms))
	for i := range ms {
		entities[i] = DistrictModelToEntity(&ms[i])
	}
	return entities
}
