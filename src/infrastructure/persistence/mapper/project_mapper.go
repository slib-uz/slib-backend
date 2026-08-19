package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ProjectModelListToEntityList(data []*models.ProjectModel) []*entity.ProjectEntity {
	var result = make([]*entity.ProjectEntity, len(data))

	for i, item := range data {
		result[i] = ProjectModelToEntity(item)
	}

	return result
}

func ProjectModelToEntity(data *models.ProjectModel) *entity.ProjectEntity {
	return entity.NewProjectEntity(data.ID, data.Title, data.LogoPath, data.Link)
}
