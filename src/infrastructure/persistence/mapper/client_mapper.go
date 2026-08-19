package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ClientModelToEntity(model *models.ClientModel) *entity.ClientEntity {
	return entity.NewClientEntity(
		model.ID,
		model.ClientID,
		model.ClientSecret,
		model.Name,
		model.Description,
		model.CallbackUrl,
		model.JournalID,
		model.IsActive,
	)
}
