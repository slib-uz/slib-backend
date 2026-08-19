package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfileModel struct {
	gorm.Model

	UserID uint       `gorm:"not null;uniqueIndex"`
	User   *UserModel `gorm:"foreignKey:UserID;references:ID"`

	Bio          *string
	Email        string             `gorm:"size:500"`
	Photo        *string            `gorm:"size:500"`
	LastOnlineAt *time.Time         `gorm:"type:timestamptz"`
	Socials      []*UserSocialModel `gorm:"foreignKey:UserProfileID"`
}

func NewUserProfileModel(userID uint, bio *string, email string, photo *string) *UserProfileModel {
	return &UserProfileModel{
		UserID: userID,
		Bio:    bio,
		Email:  email,
		Photo:  photo,
	}
}

func (UserProfileModel) TableName() string {
	return "user_profiles"
}
