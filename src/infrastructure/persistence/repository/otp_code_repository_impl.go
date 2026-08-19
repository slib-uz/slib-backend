package repository

import (
	"context"
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type OTPCodeRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewOTPCodeRepositoryImpl(baseRepository *BaseRepository) repository.OTPCodeRepository {
	return &OTPCodeRepositoryImpl{BaseRepository: baseRepository}
}

func (this *OTPCodeRepositoryImpl) MarkAsUsed(ctx context.Context, id uint) error {
	err := this.db().WithContext(ctx).Model(&models.OTPCodeModel{}).Where("id = ?", id).Update("used_at", time.Now()).Error
	return _errors.Wrap(err)
}

func (this *OTPCodeRepositoryImpl) Create(ctx context.Context, otp *entity.OTPCodeEntity) (uint, error) {
	var _model = mapper.OTPCodeEntityToModel(otp)
	err := this.db().WithContext(ctx).Create(&_model).Error
	return _model.ID, _errors.Wrap(err)
}

func (this *OTPCodeRepositoryImpl) GetByCodeAndSessionID(ctx context.Context, sessionID string, code string) (*entity.OTPCodeEntity, error) {
	var _model models.OTPCodeModel
	if err := this.db().WithContext(ctx).Where("session_id = ? AND code = ?", sessionID, code).First(&_model).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.OTPCodeModelToEntity(&_model), nil
}
