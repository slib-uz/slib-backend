package models

import "gorm.io/gorm"

type SocialModel struct {
	gorm.Model

	Name string `gorm:"size:255"`
	Icon string `gorm:"size:512"`
}

func (SocialModel) TableName() string {
	return "socials"
}

func NewSocialModel(name, icon string) *SocialModel {
	return &SocialModel{
		Name: name,
		Icon: icon,
	}
}

type UserSocialModel struct {
	gorm.Model

	UserProfileID uint   `gorm:"not null;index"` // user_profile_id
	SocialID      uint   `gorm:"not null;index"` // social_id
	Link          string `gorm:"size:255"`       // link

	Social      *SocialModel      `gorm:"foreignKey:SocialID"`                    // many2one
	UserProfile *UserProfileModel `gorm:"foreignKey:UserProfileID;references:ID"` // many2one
}

func (UserSocialModel) TableName() string {
	return "user_socials"
}

func NewUserSocialModel(uerProfileID, socialID uint, link string) *UserSocialModel {
	return &UserSocialModel{
		UserProfileID: uerProfileID,
		SocialID:      socialID,
		Link:          link,
	}
}
