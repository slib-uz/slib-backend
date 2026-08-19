package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JournalConfigModel struct {
	gorm.Model

	JournalID uint `gorm:"uniqueIndex;"`

	CreatorID uint
	Creator   *UserModel `gorm:"foreignKey:CreatorID;references:ID;constraint:OnDelete:SET NULL"`

	Journal *JournalModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	WebsiteURL string `gorm:"uniqueIndex;"`
	Conf       datatypes.JSON
	IsActive   bool `gorm:"default:true;not null;"`
}

func NewJournalConfigModel(journalID uint, creatorID uint, websiteURL string, conf datatypes.JSON, isActive bool) *JournalConfigModel {
	return &JournalConfigModel{JournalID: journalID, CreatorID: creatorID, WebsiteURL: websiteURL, Conf: conf, IsActive: isActive}
}

func (JournalConfigModel) TableName() string {
	return "journal_configs"
}
