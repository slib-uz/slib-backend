package models

import (
	"gorm.io/gorm"
)

type ProjectModel struct {
	gorm.Model

	Title    string `gorm:"size:512"`
	LogoPath string `gorm:"size:512"`
	Link     string `gorm:"size:512"`
}

func NewProjectModel(id uint, title string, logoPath string, link string) *ProjectModel {
	return &ProjectModel{Model: gorm.Model{ID: id}, Title: title, LogoPath: logoPath, Link: link}
}

func (ProjectModel) TableName() string { return "projects" }
