package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func SandboxUserModelToEntity(it *models.SandboxUserModel) *entity.SandboxUserEntity {
	return &entity.SandboxUserEntity{
		ID:          it.ID,
		FullName:    it.FullName,
		ScienceID:   it.ScienceID,
		PhoneNumber: it.PhoneNumber,
		Otp:         it.Otp,
	}
}
