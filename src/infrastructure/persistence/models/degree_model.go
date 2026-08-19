package models

import (
	"gorm.io/gorm"
)

type DegreeModel struct {
	gorm.Model

	DegreeTypeID     uint       `gorm:"not null;index"`
	DegreeType       string     `gorm:"size:255"`
	Field            *string    `gorm:"size:64"`
	DegreeStatusID   *uint      `gorm:"index"`
	DegreeStatusName *string    `gorm:"size:255"`
	ConfirmedDate    *string    `gorm:"type:date"`
	Protocol         string     `gorm:"size:64"`
	UserID           *uint      `gorm:"uniqueIndex:idx_user_id_degree_type_id"`
	User             *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // many2one
}

func (DegreeModel) TableName() string {
	return "degrees"
}

func NewDegreeModel(degreeTypeID uint, degreeType string, field *string, degreeStatusID *uint, degreeStatusName *string, confirmedDate *string, protocol string, userID *uint) *DegreeModel {
	return &DegreeModel{
		DegreeTypeID:     degreeTypeID,
		DegreeType:       degreeType,
		Field:            field,
		DegreeStatusID:   degreeStatusID,
		DegreeStatusName: degreeStatusName,
		ConfirmedDate:    confirmedDate,
		Protocol:         protocol,
		UserID:           userID,
	}
}
