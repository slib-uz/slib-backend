package repository

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/db"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type UserSocialRepositoryImpl struct {
	database *db.Database
}

// @inject
func NewUserSocialRepository(database *db.Database) repository.UserSocialRepository {
	return &UserSocialRepositoryImpl{database: database}
}

func (this *UserSocialRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *UserSocialRepositoryImpl) GetByUserID(UserID uint) (*entity.UserSocialEntity, error) {
	var social models.UserSocialModel

	result := this.db().Model(&models.UserSocialModel{}).
		Preload("Social").
		Joins("LEFT JOIN user_profiles up ON user_socials.user_profile_id = up.id").
		Where("up.user_id = ?", UserID).
		First(&social)

	if result.Error != nil {
		return nil, infraError.Wrap(result.Error)
	}
	socialEntity := mapper.UserSocialModelToEntity(&social)
	return socialEntity, nil
}

func (this *UserSocialRepositoryImpl) GetByID(id uint) (*entity.UserSocialEntity, error) {
	var social models.UserSocialModel

	result := this.db().Preload("Social").First(&social, id)
	if result.Error != nil {
		return nil, infraError.Wrap(result.Error)
	}
	socialEntity := mapper.UserSocialModelToEntity(&social)
	return socialEntity, nil
}

func (this *UserSocialRepositoryImpl) Create(entity *entity.UserSocialInputEntity) error {
	socialModel := mapper.UserSocialEntityToModel(entity)
	if err := this.db().Create(&socialModel).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}

func (this *UserSocialRepositoryImpl) GetAll() ([]*entity.UserSocialEntity, error) {
	var socials []*models.UserSocialModel
	if err := this.db().Preload("Social").Find(&socials).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	var socialModels []*models.UserSocialModel
	for _, social := range socials {
		socialModels = append(socialModels, social)
	}

	socialsEntity := mapper.UserSocialsModelToEntity(socialModels)

	return socialsEntity, nil
}

func (this *UserSocialRepositoryImpl) Update(id uint, entity *entity.UserSocialInputEntity) error {
	var social models.UserSocialModel

	if err := this.db().First(&social, id).Error; err != nil {
		return infraError.Wrap(err)
	}

	updatedModel := mapper.UserSocialEntityToModel(entity)
	social.Link = updatedModel.Link
	social.SocialID = updatedModel.SocialID

	if err := this.db().Save(&social).Error; err != nil {
		return infraError.Wrap(err)
	}

	return nil
}

func (this *UserSocialRepositoryImpl) Delete(id uint) error {
	var social models.UserSocialModel
	if err := this.db().First(&social, id).Error; err != nil {
		return infraError.Wrap(err)
	}
	if err := this.db().Unscoped().Delete(&social).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}

func (this *UserSocialRepositoryImpl) Exists(userProfileID uint, socialID uint) (bool, error) {
	var count int64
	if err := this.db().Model(&models.UserSocialModel{}).Where("user_profile_id = ? AND social_id = ?", userProfileID, socialID).Count(&count).Error; err != nil {
		return false, infraError.Wrap(err)
	}
	return count > 0, nil
}
