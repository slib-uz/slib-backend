package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/session"
)

type UserProfileRepository interface {
	GetByUserID(userID uint) (*entity2.UserMeEntity, error)
	GetUserProfileByUserID(userID uint) (*entity2.UserProfileEntity, error)
	Create(tx session.Tx, userID uint) error
	GetProfileByScienceID(scienceID string) (*entity2.UserProfileEntity, error)
	Update(id uint, bio *string, email string, photo *string) error
	UpdateLastOnlineAt(userID uint) error
	GetProfileByUserID(userID uint) (*entity2.UserProfileEntity, error)
}
