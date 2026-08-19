package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StudyFieldModel struct {
	gorm.Model

	Name     datatypes.JSON    `gorm:"type:jsonb;index"`
	Code     *uint             `gorm:"uniqueIndex"`
	ParentID *uint             `gorm:"index"`
	Parent   *StudyFieldModel  `gorm:"foreignKey:ParentID;references:ID;constraint,OnDelete:SET NULL"`
	Children []StudyFieldModel `gorm:"foreignKey:ParentID"`
}

func NewStudyFieldModel(id uint, name datatypes.JSON, parentID *uint, code *uint) StudyFieldModel {
	return StudyFieldModel{Model: gorm.Model{ID: id}, Name: name, ParentID: parentID, Code: code}
}

func (*StudyFieldModel) TableName() string {
	return "study_fields"
}
