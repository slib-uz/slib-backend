package models

import (
	"gorm.io/gorm"
)

type PartnerModel struct {
	gorm.Model

	Title    string `gorm:"size:512"`
	LogoPath string `gorm:"size:512"`
	Link     string `gorm:"size:512"`
}

func NewPartnerModel(id uint, title string, logoPath string, link string) *PartnerModel {
	return &PartnerModel{Model: gorm.Model{ID: id}, Title: title, LogoPath: logoPath, Link: link}
}

func (PartnerModel) TableName() string { return "partners" }
