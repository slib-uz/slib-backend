package models

import "gorm.io/gorm"

type ReferenceModel struct {
	gorm.Model
	Name string

	ArticleID uint
	Article   *ArticleModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewReferenceModel(name string, articleID uint) *ReferenceModel {
	return &ReferenceModel{Name: name, ArticleID: articleID}
}
