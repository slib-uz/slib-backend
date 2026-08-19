package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
)

type NotificationSendResultRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewNotificationSendResultRepository(baseRepository *BaseRepository) repository.NotificationSendResultRepository {
	return &NotificationSendResultRepositoryImpl{BaseRepository: baseRepository}
}

func (this *NotificationSendResultRepositoryImpl) Create(entity *entity.NotificationSendResultEntity) error {
	var _model = mapper.NotificationSendResultEntityToModel(entity)
	return this.db().Create(_model).Error
}
