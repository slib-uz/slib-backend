package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func SoatoModelToEntity(model *models.SoatoClassifierModel) *entity.SoatoClassifierEntity {
	if model == nil {
		return nil
	}

	e := &entity.SoatoClassifierEntity{
		ID:       model.ID,
		Code:     model.Code,
		Name:     model.Name(),
		ParentID: model.ParentID,
	}

	if model.Parent != nil {
		e.Parent = SoatoModelToEntity(model.Parent)
	}

	if len(model.Children) > 0 {
		children := make([]*entity.SoatoClassifierEntity, len(model.Children))
		for i, child := range model.Children {
			children[i] = SoatoModelToEntity(&child)
		}
		e.Children = children
	}

	return e
}

func SoatoModelListToEntityList(models []*models.SoatoClassifierModel) []*entity.SoatoClassifierEntity {
	entities := make([]*entity.SoatoClassifierEntity, len(models))
	for i, model := range models {
		entities[i] = SoatoModelToEntity(model)
	}
	return entities
}
