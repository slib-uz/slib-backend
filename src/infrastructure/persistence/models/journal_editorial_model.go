package models

import "gorm.io/gorm"

type JournalEditorialModel struct {
	gorm.Model

	JournalID uint          `gorm:"not null;index"`
	Journal   *JournalModel `gorm:"foreignKey:JournalID;constraint:OnDelete:CASCADE"`

	FullName  string `gorm:"not null;size:128"`
	RoleCode  int    `gorm:"not null;default:1"`
	RoleTitle string `gorm:"not null;size:128"`
	Photo     string `gorm:"size:128"`
	ScienceID string `gorm:"size:128"`
	Workplace string `gorm:"size:128"`
	Position  string `gorm:"size:128"`
	Order     int    `gorm:"column:order;not null;default:0"`
}

func (JournalEditorialModel) TableName() string {
	return "journal_editorials"
}

func NewJournalEditorialModel(journalID uint, fullName string, roleCode int, roleTitle, photo, scienceID, workplace, position string, order int) *JournalEditorialModel {
	return &JournalEditorialModel{
		JournalID: journalID,
		FullName:  fullName,
		RoleCode:  roleCode,
		RoleTitle: roleTitle,
		Photo:     photo,
		ScienceID: scienceID,
		Workplace: workplace,
		Position:  position,
		Order:     order,
	}
}
