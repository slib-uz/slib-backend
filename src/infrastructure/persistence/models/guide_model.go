package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GuideModel struct {
	gorm.Model

	Title       datatypes.JSON `gorm:"not null"`
	Description datatypes.JSON `gorm:"not null"`
	FilePath    string
	VideoUrl    string `gorm:"default:''"`
}

func NewGuideModel(id uint, title datatypes.JSON, description datatypes.JSON, filePath string, videoUrl string) *GuideModel {
	return &GuideModel{Model: gorm.Model{ID: id}, Title: title, Description: description, FilePath: filePath, VideoUrl: videoUrl}
}

func (GuideModel) TableName() string { return "guides" }
