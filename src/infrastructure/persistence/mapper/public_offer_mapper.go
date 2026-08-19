package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func PublicOfferModelListToEntity(data []*models.PublicOfferModel) []*entity.PublicOfferEntity {
	var result = make([]*entity.PublicOfferEntity, len(data))

	for i, item := range data {
		result[i] = PublicOfferModelToEntity(item)
	}

	return result
}

func PublicOfferModelToEntity(data *models.PublicOfferModel) *entity.PublicOfferEntity {
	return entity.NewPublicOfferEntity(data.ID, FromGormJson[map[string]string](data.Description))
}
