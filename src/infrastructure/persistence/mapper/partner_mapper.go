package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func PartnerModelListToEntityList(data []*models.PartnerModel) []*entity.PartnerEntity {
	var result = make([]*entity.PartnerEntity, len(data))

	for i, item := range data {
		result[i] = PartnerModelToEntity(item)
	}

	return result
}

func PartnerModelToEntity(data *models.PartnerModel) *entity.PartnerEntity {
	return entity.NewPartnerEntity(data.ID, data.Title, data.LogoPath, data.Link)
}
