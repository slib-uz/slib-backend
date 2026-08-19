package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func SupportDialogOwnerModelToEntity(data *models.UserModel) *entity.SupportDialogOwnerEntity {
	return entity.NewSupportDialogOwnerEntity(data.ID, data.FullName, data.ScienceID)
}
