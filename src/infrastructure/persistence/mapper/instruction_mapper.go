package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func InstructionModelToEntity(m *models.InstructionModel) *entity.InstructionEntity {
	return &entity.InstructionEntity{
		ID:          m.ID,
		Key:         m.Key,
		VideoLink:   m.VideoLink,
		Description: m.Description,
	}
}

func InstructionEntityToModel(e *entity.InstructionEntity) *models.InstructionModel {
	return &models.InstructionModel{
		Key:         e.Key,
		VideoLink:   e.VideoLink,
		Description: e.Description,
	}
}
