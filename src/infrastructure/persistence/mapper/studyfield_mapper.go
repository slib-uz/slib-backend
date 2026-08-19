package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func StudyFieldModelToEntity(model *models.StudyFieldModel) *entity.StudyFieldEntity {

	name := FromGormJson[map[string]string](model.Name)

	var parent *entity.StudyFieldEntity
	if model.Parent != nil {
		parent = StudyFieldModelToEntity(model.Parent)
	}

	var children []*entity.StudyFieldEntity
	if len(model.Children) > 0 {
		children = make([]*entity.StudyFieldEntity, len(model.Children))
		for i, child := range model.Children {
			children[i] = StudyFieldModelToEntity(&child)
		}
	}

	studyFieldEntity := entity.NewStudyFieldEntity(
		model.ID,
		name,
		model.ParentID,
		model.Code,
		parent,
	)
	studyFieldEntity.Children = children

	return studyFieldEntity
}

func StudyFieldEntityToModel(entity *entity.StudyFieldEntity) models.StudyFieldModel {
	name := ToGormJson(entity.Name)

	return models.NewStudyFieldModel(
		entity.ID,
		name,
		entity.ParentID,
		entity.Code,
	)
}
