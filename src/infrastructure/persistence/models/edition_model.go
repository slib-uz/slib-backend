package models

import (
	"time"

	"gorm.io/gorm"
)

type EditionModel struct {
	gorm.Model

	JournalID uint
	Journal   *JournalModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	File        string `gorm:"size:1024"`
	CoverImage  string `gorm:"size:1024"`
	Name        string `gorm:"size:255"`
	Number      string `gorm:"size:255"`
	Volume      string `gorm:"size:255"`
	PublishedAt time.Time
}

func NewEditionModel(journalID uint, file string, coverImage string, name string, number string, volume string, publishedAt time.Time) *EditionModel {
	return &EditionModel{JournalID: journalID, File: file, CoverImage: coverImage, Name: name, Number: number, Volume: volume, PublishedAt: publishedAt}
}

func (EditionModel) TableName() string {
	return "editions"
}
