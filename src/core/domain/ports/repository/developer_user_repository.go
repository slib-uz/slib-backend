package repository

import (
	"slib.uz/src/core/domain/entity"
)

type DeveloperUserRepository interface {
	GetByID(id uint) (*entity.DeveloperUserEntity, error)
	GetByUsername(username string) (*entity.DeveloperUserEntity, error)
}
