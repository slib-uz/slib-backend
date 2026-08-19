package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func OTPCodeEntityToModel(it *entity.OTPCodeEntity) *models.OTPCodeModel {
	return models.NewOTPCodeModel(
		it.Phone,
		it.Code,
		it.SessionID,
		it.Purpose,
		it.ExpiresAt,
	)
}

func OTPCodeModelToEntity(it *models.OTPCodeModel) *entity.OTPCodeEntity {
	return entity.NewOTPCodeEntity(
		it.ID,
		it.Phone,
		it.Code,
		it.SessionID,
		it.Purpose,
		it.ExpiresAt,
		it.UsedAt,
		it.CreatedAt,
	)
}
