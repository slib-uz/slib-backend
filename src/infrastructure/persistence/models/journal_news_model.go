package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JournalNewsModel struct {
	gorm.Model

	JournalID uint          `gorm:"not null;index"`
	Journal   *JournalModel `gorm:"foreignKey:JournalID;constraint:OnDelete:CASCADE"`
	Title     datatypes.JSON
	Body      datatypes.JSON
	Image     *string `gorm:"size:1024"`
}

func (JournalNewsModel) TableName() string {
	return "journal_news"
}
