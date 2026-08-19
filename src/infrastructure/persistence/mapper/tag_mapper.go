package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func TagModelToEntity(tag *models.TagModel) *entity.TagEntity {
	return entity.NewTagEntity(tag.ID, tag.Name, tag.Lang)
}

func TagModelsToNamesByLang(tags []*models.TagModel) entity.TagNamesByLang {
	if len(tags) == 0 {
		return nil
	}
	entities := make([]*entity.TagEntity, len(tags))
	for i, tag := range tags {
		entities[i] = TagModelToEntity(tag)
	}
	return entity.TagNamesByLangFromEntities(entities)
}
