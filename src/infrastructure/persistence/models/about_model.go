package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AboutModel struct {
	gorm.Model

	Content datatypes.JSON `gorm:"not null"`
}

func NewAboutModel(id uint, content datatypes.JSON) *AboutModel {
	return &AboutModel{Model: gorm.Model{ID: id}, Content: content}
}

func (AboutModel) TableName() string { return "abouts" }
