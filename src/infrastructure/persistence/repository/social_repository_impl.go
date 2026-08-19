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

type SocialRepositoryImpl struct {
	database *db.Database
}

// @inject
func NewSocialRepository(database *db.Database) repository.SocialRepository {
	return &SocialRepositoryImpl{database: database}
}

func (this *SocialRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this SocialRepositoryImpl) Create(entity *entity.SocialEntity) error {

	socialModel := mapper.SocialEntityToModel(entity)
	if err := this.db().Create(&socialModel).Error; err != nil {
		return err
	}
	return nil
}

func (this SocialRepositoryImpl) GetAll() ([]*entity.SocialEntity, error) {
	var socials []*models.SocialModel
	if err := this.db().Find(&socials).Error; err != nil {
		return nil, err
	}
	socialsEntity := mapper.SocialsModelToEntity(&socials)
	return socialsEntity, nil
}

func (this SocialRepositoryImpl) Update(id uint, entity *entity.SocialEntity) error {
	var social models.SocialModel
	if err := this.db().First(&social, id).Error; err != nil {
		return infraError.Wrap(err)
	}
	updatedModel := mapper.SocialEntityToModel(entity)
	social.Name = updatedModel.Name
	social.Icon = updatedModel.Icon

	if err := this.db().Save(&social).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}

func (this SocialRepositoryImpl) Delete(id uint) error {
	var social models.SocialModel
	if err := this.db().First(&social, id).Error; err != nil {
		return infraError.Wrap(err)
	}
	if err := this.db().Unscoped().Delete(&social).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}
