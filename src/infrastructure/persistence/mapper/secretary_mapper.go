package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func SecretaryListModelToEntity(data []*models.UserModel) []*entity.SecretaryEntity {
	var result = make([]*entity.SecretaryEntity, len(data))
	for i, item := range data {
		result[i] = SecretaryModelToEntity(item)
	}
	return result
}

func SecretaryModelToEntity(model *models.UserModel) *entity.SecretaryEntity {
	return entity.NewSecretaryEntity(model.ID, model.FullName, model.ScienceID)
}
