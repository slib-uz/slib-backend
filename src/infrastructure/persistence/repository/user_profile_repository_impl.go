package repository

import (
	"errors"
	"time"

	"slib.uz/src/infrastructure/db"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/session"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type UserProfileRepositoryImpl struct {
	database *db.Database
}

// @inject
func NewUserProfileRepository(database *db.Database) repository.UserProfileRepository {
	return &UserProfileRepositoryImpl{database: database}
}

func (this *UserProfileRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *UserProfileRepositoryImpl) Create(tx session.Tx, userID uint) error {
	db := tx.(*gorm.DB)

	return _errors.Wrap(db.Create(&models.UserProfileModel{UserID: userID}).Error)

}

func (this *UserProfileRepositoryImpl) GetByUserID(userID uint) (*entity.UserMeEntity, error) {

	var profile models.UserProfileModel

	if result := this.db().
		Preload("User").
		Preload("User.Roles").
		Preload("User.Roles.Journal").
		Preload("User.Roles.Publisher").
		Preload("User.Roles.Institution").
		Where("user_id = ?", userID).
		First(&profile); result.Error != nil {
		return nil, _errors.Wrap(result.Error)
	}
	return mapper.UserMeProfileModelToEntity(&profile), nil
}

func (this *UserProfileRepositoryImpl) GetUserProfileByUserID(userID uint) (*entity.UserProfileEntity, error) {
	var profile models.UserProfileModel

	result := this.db().
		Preload("User").
		Preload("Socials").
		Preload("Socials.Social").
		Where("user_id = ?", userID).
		First(&profile)

	if result.Error != nil {
		return nil, _errors.Wrap(result.Error)
	}
	return mapper.UserProfileModelToEntity(&profile), nil
}

func (this *UserProfileRepositoryImpl) GetProfileByScienceID(scienceID string) (*entity.UserProfileEntity, error) {
	var profile models.UserProfileModel

	result := this.db().
		Preload("User").
		Preload("Socials").
		Preload("Socials.Social").
		Joins("JOIN users ON users.id = user_profiles.user_id").
		Where("users.science_id = ?", scienceID).
		First(&profile)

	if result.Error != nil {
		return nil, _errors.Wrap(result.Error)
	}
	return mapper.UserProfileModelToEntity(&profile), nil
}

func (this *UserProfileRepositoryImpl) Update(id uint, bio *string, email string, photo *string) error {
	return this.db().
		Model(&models.UserProfileModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"bio":   bio,
			"email": email,
			"photo": photo,
		}).Error
}

func (this *UserProfileRepositoryImpl) UpdateLastOnlineAt(userID uint) error {
	now := time.Now()
	threshold := now.Add(-5 * time.Minute)

	err := this.db().
		Model(&models.UserProfileModel{}).
		Where("user_id = ? AND (last_online_at IS NULL OR last_online_at < ?)", userID, threshold).
		Update("last_online_at", now).Error

	return _errors.Wrap(err)
}

func (this *UserProfileRepositoryImpl) GetProfileByUserID(userID uint) (*entity.UserProfileEntity, error) {
	var profile models.UserProfileModel

	result := this.db().
		Preload("User").
		Preload("Socials").
		Preload("Socials.Social").
		Model(&models.UserProfileModel{}).
		Where("user_id = ?", userID).
		First(&profile)
	if result.Error != nil {
		return nil, _errors.Wrap(result.Error)
	}

	if profile.ID == 0 {
		return nil, _errors.Wrap(errors.New("profile not found"))
	}
	return mapper.UserProfileModelToEntity(&profile), nil
}
