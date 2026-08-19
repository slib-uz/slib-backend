package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func AboutModelToEntity(data *models.AboutModel) *entity2.AboutEntity {
	return entity2.NewAboutEntity(data.ID, FromGormJson[map[string]string](data.Content))
}

func AboutModelListToEntityList(data []*models.AboutModel) []*entity2.AboutEntity {
	var result = make([]*entity2.AboutEntity, len(data))
	for i, item := range data {
		result[i] = AboutModelToEntity(item)
	}
	return result
}
