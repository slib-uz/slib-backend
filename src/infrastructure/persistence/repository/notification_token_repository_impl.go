package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/models"
)

type NotificationTokenRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewNotificationTokenRepository(repository *BaseRepository) repository.NotificationTokenRepository {
	return &NotificationTokenRepositoryImpl{BaseRepository: repository}
}

func (this *NotificationTokenRepositoryImpl) CreateOrUpdate(t *entity.NotificationTokenEntity) (isCreated bool, err error) {
	var _model models.NotificationTokenModel

	result := this.db().Where("token = ?", t.Token).First(&_model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return true, this.db().Create(
				models.NewNotificationTokenModel(t.UserID, t.Token, t.Segment),
			).Error
		}
		return false, result.Error
	}

	_model.UserID = t.UserID
	return false, this.db().Updates(&_model).Error
}

func (this *NotificationTokenRepositoryImpl) GetTokensByUserID(userID uint) ([]string, error) {
	var tokens []string

	err := this.db().Model(&models.NotificationTokenModel{}).
		Where("user_id = ?", userID).
		Pluck("token", &tokens).Error

	return tokens, err
}

func (this *NotificationTokenRepositoryImpl) Delete(userID uint, token string) error {
	return this.db().Unscoped().Where("token = ? AND user_id = ?", token, userID).Delete(&models.NotificationTokenModel{}).Error
}
