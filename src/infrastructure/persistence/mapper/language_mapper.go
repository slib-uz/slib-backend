package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func LanguageModelToEntity(model *models.LanguageModel) *entity.LanguageEntity {

	var name = make(map[string]string)
	if model.Name != nil {
		name = FromGormJson[map[string]string](model.Name)
	}

	return entity.NewLanguageEntity(
		model.ID,
		name,
		model.Code,
	)
}
