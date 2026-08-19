package models

import (
	"gorm.io/gorm"
)

type JournalRatingModel struct {
	gorm.Model
	UserID    uint          `gorm:"not null;index"`
	User      *UserModel    `gorm:"foreignKey:UserID"`
	Journal   *JournalModel `gorm:"foreignKey:JournalID;constraint:OnDelete:CASCADE;"`
	JournalID uint          `gorm:"not null;index"`
	Stars     uint          `gorm:"type:smallint;check:stars >= 1 AND stars <= 5;not null"`
	Review    string        `gorm:"type:text;size:512;"`
}

func NewJournalRatingModel(userID uint, journalID uint, stars uint, review string) *JournalRatingModel {
	return &JournalRatingModel{UserID: userID, JournalID: journalID, Stars: stars, Review: review}
}

func (JournalRatingModel) TableName() string {
	return "journal_ratings"
}
