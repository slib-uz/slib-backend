package repository

import (
	"slib.uz/src/core/domain/entity"
)

type UserSocialRepository interface {
	GetByUserID(UserID uint) (*entity.UserSocialEntity, error)
	GetByID(id uint) (*entity.UserSocialEntity, error)
	Create(entity *entity.UserSocialInputEntity) error
	GetAll() ([]*entity.UserSocialEntity, error)
	Update(id uint, entity *entity.UserSocialInputEntity) error
	Delete(id uint) error
	Exists(userProfileID uint, socialID uint) (bool, error)
}
