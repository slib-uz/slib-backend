package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type NewsModel struct {
	gorm.Model

	Title      datatypes.JSON     `gorm:"not null"`
	Body       datatypes.JSON     `gorm:"not null"`
	CategoryID uint               `gorm:"not null"`
	Category   *NewsCategoryModel `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
	Image      string             `gorm:"not null"`
	ViewsCount uint               `gorm:"not null;default:0"`
}

func NewNewsModel(id uint, title datatypes.JSON, body datatypes.JSON, categoryID uint, image string, viewCounts uint) *NewsModel {
	return &NewsModel{Model: gorm.Model{ID: id}, Title: title, Body: body, CategoryID: categoryID, Image: image, ViewsCount: viewCounts}
}

func (NewsModel) TableName() string { return "news" }
