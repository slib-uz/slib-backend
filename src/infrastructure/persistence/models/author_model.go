package models

import (
	"gorm.io/gorm"
)

type AuthorModel struct {
	gorm.Model

	ScienceID string `gorm:"not null;uniqueIndex;size:32"`
	FullName  string `gorm:"size:255"`

	ArticleAuthorAffiliations *[]ArticleAuthorAffiliationModel `gorm:"foreignKey:AuthorID"`
	Users                     *[]UserModel                     `gorm:"-"`
}

func (*AuthorModel) TableName() string {
	return "authors"
}
