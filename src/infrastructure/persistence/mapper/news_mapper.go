package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func NewsListModelToEntity(data []*models.NewsModel) []*entity.NewsEntity {
	var result = make([]*entity.NewsEntity, len(data))

	for i, item := range data {
		result[i] = NewsModelToEntity(item)
	}

	return result
}

func NewsModelToEntity(data *models.NewsModel) *entity.NewsEntity {
	return entity.NewNewsEntity(data.ID, FromGormJson[map[string]string](data.Title), FromGormJson[map[string]string](data.Body), data.CategoryID, data.Image, data.ViewsCount, data.CreatedAt)
}
