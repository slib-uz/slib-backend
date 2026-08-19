package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ChiefEditorListModelToEntity(data []*models.UserModel) []*entity.ChiefEditorEntity {
	var result = make([]*entity.ChiefEditorEntity, len(data))
	for i, item := range data {
		result[i] = ChiefEditorModelToEntity(item)
	}
	return result
}

func ChiefEditorModelToEntity(model *models.UserModel) *entity.ChiefEditorEntity {
	return entity.NewChiefEditorEntity(model.ID, model.FullName, model.ScienceID, model.Photo)
}
