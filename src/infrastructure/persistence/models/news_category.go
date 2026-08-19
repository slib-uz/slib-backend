package models

import "gorm.io/gorm"

type NewsCategoryModel struct {
	gorm.Model

	Name string `gorm:"not null;uniqueIndex;size:255"`
}

func (NewsCategoryModel) TableName() string {
	return "news_categories"
}
