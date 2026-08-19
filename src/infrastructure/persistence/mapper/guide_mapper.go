package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func GuideModelListToEntityList(data []*models.GuideModel) []*entity2.GuideListEntity {
	var result = make([]*entity2.GuideListEntity, len(data))

	for i, item := range data {
		result[i] = GuideModelToEntity(item)
	}

	return result
}

func GuideModelToEntity(data *models.GuideModel) *entity2.GuideListEntity {
	return entity2.NewListGuideEntity(data.ID, FromGormJson[map[string]string](data.Title))
}

func GuideRetrieveModelToEntity(data *models.GuideModel) *entity2.GuideRetrieveEntity {
	return entity2.NewGuideRetrieveEntity(data.ID, FromGormJson[map[string]string](data.Title), FromGormJson[map[string]string](data.Description), data.FilePath, data.VideoUrl)
}
