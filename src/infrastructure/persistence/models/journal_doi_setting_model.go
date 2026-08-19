package models

import "gorm.io/gorm"

type JournalDoiSettingModel struct {
	gorm.Model
	JournalID   uint          `gorm:"uniqueIndex;not null"`
	Journal     *JournalModel `gorm:"constraint:OnDelete:CASCADE;"`
	JournalName string        `gorm:"column:journal_name"`
	Username    string        `gorm:"not null"`
	Password    string        `gorm:"not null"`
	DOIPrefix   string        `gorm:"column:doi_prefix;not null"`
	DOISuffix   string        `gorm:"column:doi_suffix;not null"`
}

func NewJournalDoiSettingModel(journalID uint, journalName string, username string, password string, doiPrefix string, doiSuffix string) *JournalDoiSettingModel {
	return &JournalDoiSettingModel{JournalID: journalID, JournalName: journalName, Username: username, Password: password, DOIPrefix: doiPrefix, DOISuffix: doiSuffix}
}

func (*JournalDoiSettingModel) TableName() string {
	return "journal_doi_settings"
}
