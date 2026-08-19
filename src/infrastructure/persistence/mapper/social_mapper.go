package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func SocialEntityToModel(entity *entity.SocialEntity) *models.SocialModel {
	return models.NewSocialModel(
		entity.Name,
		entity.Icon,
	)
}

func SocialsModelToEntity(models *[]*models.SocialModel) []*entity.SocialEntity {
	if models == nil {
		return nil
	}
	socials := make([]*entity.SocialEntity, len(*models))
	for i, e := range *models {
		socials[i] = entity.NewSocialEntity(
			e.ID,
			e.Name,
			e.Icon,
		)
	}
	return socials
}
