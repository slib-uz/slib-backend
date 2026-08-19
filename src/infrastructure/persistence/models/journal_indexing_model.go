package models

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type JournalIndexingModel struct {
	gorm.Model
	IndexingType enum.IndexingType `gorm:"size:32;not null;index:idx_journal_indexing"`
	URL          string            `gorm:"not null;size:255;"`
	JournalID    uint              `gorm:"index:idx_journal_indexing;not null"`
	Journal      *JournalModel     `gorm:"foreignKey:JournalID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func NewJournalIndexingModel(indexingType enum.IndexingType, url string, journalID uint) *JournalIndexingModel {
	return &JournalIndexingModel{IndexingType: indexingType, URL: url, JournalID: journalID}
}

func (JournalIndexingModel) TableName() string {
	return "journal_indexing"
}
