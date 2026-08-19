package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalRaterModelToEntity(model *models.UserModel) *entity.JournalRaterEntity {
	return entity.NewJournalRaterEntity(model.ID, model.FullName, model.ScienceID)
}
