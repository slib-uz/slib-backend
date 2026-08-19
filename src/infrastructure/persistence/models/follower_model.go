package models

import (
	"gorm.io/gorm"
)

type FollowerModel struct {
	gorm.Model

	FollowingID uint `gorm:"column:following_id;not null;index"` // following_id is who wants to follow anyone profile
	FollowedID  uint `gorm:"column:followed_id;not null;index"`  // followed_id is follower to profile

	Following *UserModel `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE"`
	Followed  *UserModel `gorm:"foreignKey:FollowedID;constraint:OnDelete:CASCADE"`
}

func (FollowerModel) TableName() string {
	return "followers"
}

func NewFollowerModel(followingID, followedID uint) *FollowerModel {
	return &FollowerModel{
		FollowingID: followingID,
		FollowedID:  followedID,
	}
}
