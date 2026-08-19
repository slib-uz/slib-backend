package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type NotificationTemplateRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewNotificationTemplateRepository(repository *BaseRepository) repository.NotificationTemplateRepository {
	return &NotificationTemplateRepositoryImpl{BaseRepository: repository}
}

func (this *NotificationTemplateRepositoryImpl) GetByKey(key enum.NotificationTemplate) (*entity.NotificationTemplateEntity, error) {
	var _model models.NotificationTemplateModel

	if err := this.db().Where("key = ?", key).First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.NotificationTemplateModelToEntity(&_model), nil
}
