package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func NewsCategoryModelToEntity(model *models.NewsCategoryModel) *entity.NewsCategoryEntity {
	return entity.NewNewsCategoryEntity(
		model.ID,
		model.Name,
	)
}
