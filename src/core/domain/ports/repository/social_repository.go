package repository

import (
	"slib.uz/src/core/domain/entity"
)

type SocialRepository interface {
	Create(entity *entity.SocialEntity) error
	GetAll() ([]*entity.SocialEntity, error)
	Update(id uint, entity *entity.SocialEntity) error
	Delete(id uint) error
}
