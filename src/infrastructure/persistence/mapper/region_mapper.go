package mapper

import (
	"encoding/json"

	"gorm.io/datatypes"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func RegionModelToEntity(m *models.RegionModel) *entity.RegionEntity {
	return entity.NewRegionEntity(m.ID, jsonNameToMap(m.Name), m.Soato)
}

func RegionModelsToEntities(ms []models.RegionModel) []*entity.RegionEntity {
	entities := make([]*entity.RegionEntity, len(ms))
	for i := range ms {
		entities[i] = RegionModelToEntity(&ms[i])
	}
	return entities
}

func MapRegionDistrictFromModel(
	regionID, districtID *uint,
	region *models.RegionModel,
	district *models.DistrictModel,
) (*uint, *uint, *entity.RegionEntity, *entity.DistrictEntity) {
	var regionEntity *entity.RegionEntity
	var districtEntity *entity.DistrictEntity
	if region != nil {
		regionEntity = RegionModelToEntity(region)
	}
	if district != nil {
		districtEntity = DistrictModelToEntity(district)
	}
	return regionID, districtID, regionEntity, districtEntity
}

func jsonNameToMap(name datatypes.JSON) map[string]string {
	var result map[string]string
	if err := json.Unmarshal(name, &result); err != nil {
		return map[string]string{}
	}
	return result
}
