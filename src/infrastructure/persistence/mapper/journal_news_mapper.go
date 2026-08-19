package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalNewsModelToEntity(m *models.JournalNewsModel) *entity.JournalNewsEntity {
	return &entity.JournalNewsEntity{
		ID:        m.ID,
		JournalID: m.JournalID,
		Title:     FromGormJson[map[string]string](m.Title),
		Body:      FromGormJson[map[string]string](m.Body),
		Image:     m.Image,
		CreatedAt: m.CreatedAt,
	}
}

func JournalNewsEntityToModel(e *entity.JournalNewsEntity) *models.JournalNewsModel {
	return &models.JournalNewsModel{
		JournalID: e.JournalID,
		Title:     ToGormJson(e.Title),
		Body:      ToGormJson(e.Body),
		Image:     e.Image,
	}
}
