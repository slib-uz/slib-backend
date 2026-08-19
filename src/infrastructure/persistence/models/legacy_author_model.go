package models

import "gorm.io/gorm"

type LegacyAuthorModel struct {
	gorm.Model

	FullName string          `gorm:"size:512;index"`
	Articles []*ArticleModel `gorm:"many2many:legacy_author_articles;"`
}

func (*LegacyAuthorModel) TableName() string {
	return "legacy_authors"
}

func NewLegacyAuthorModel(fullName string) *LegacyAuthorModel {
	return &LegacyAuthorModel{FullName: fullName}
}
