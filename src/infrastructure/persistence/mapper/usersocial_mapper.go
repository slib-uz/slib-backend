package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/handlers/social/schema"
	"slib.uz/src/infrastructure/persistence/models"
)

func UserSocialsModelToEntity(models []*models.UserSocialModel) []*entity.UserSocialEntity {
	if models == nil {
		return nil
	}
	socials := make([]*entity.UserSocialEntity, 0)
	for _, e := range models {
		socials = append(socials, entity.NewUserSocialEntity(
			e.ID,
			e.UserProfileID,
			e.Social.Name,
			e.Link,
			e.Social.Icon,
			e.SocialID,
		))
	}
	return socials
}

func UserSocialModelToEntity(model *models.UserSocialModel) *entity.UserSocialEntity {
	return entity.NewUserSocialEntity(
		model.ID,
		model.UserProfileID,
		model.Social.Name,
		model.Link,
		model.Social.Icon,
		model.SocialID,
	)
}

func UserSocialRequestToDTO(req *schema.UserSocialRequest) *entity.UserSocialInputEntity {
	return entity.NewUserSocialInputEntity(
		0,
		0,
		req.SocialID,
		req.Link,
	)
}

func UserSocialEntityToModel(entity *entity.UserSocialInputEntity) *models.UserSocialModel {
	return models.NewUserSocialModel(
		entity.UserProfileID,
		entity.SocialID,
		entity.Link,
	)
}
