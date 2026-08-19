package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalEditorialModelToEntity(m *models.JournalEditorialModel) *entity.JournalEditorialEntity {
	return &entity.JournalEditorialEntity{
		ID:        m.ID,
		JournalID: m.JournalID,
		FullName:  m.FullName,
		RoleCode:  m.RoleCode,
		RoleTitle: m.RoleTitle,
		Photo:     m.Photo,
		ScienceID: m.ScienceID,
		Workplace: m.Workplace,
		Position:  m.Position,
		Order:     m.Order,
		CreatedAt: m.CreatedAt,
	}
}

func JournalEditorialEntityToModel(e *entity.JournalEditorialEntity) *models.JournalEditorialModel {
	return &models.JournalEditorialModel{
		JournalID: e.JournalID,
		FullName:  e.FullName,
		RoleCode:  e.RoleCode,
		RoleTitle: e.RoleTitle,
		Photo:     e.Photo,
		ScienceID: e.ScienceID,
		Workplace: e.Workplace,
		Position:  e.Position,
		Order:     e.Order,
	}
}
