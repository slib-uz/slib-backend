package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type SandboxUserRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewSandboxUserRepositoryImpl(baseRepository *BaseRepository) repository.SandboxUserRepository {
	return &SandboxUserRepositoryImpl{BaseRepository: baseRepository}
}

func (this *SandboxUserRepositoryImpl) GetByPhoneNumber(phoneNumber string) (*entity.SandboxUserEntity, error) {
	var _model models.SandboxUserModel
	if err := this.db().Where("phone_number = ?", phoneNumber).First(&_model).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.SandboxUserModelToEntity(&_model), nil
}
