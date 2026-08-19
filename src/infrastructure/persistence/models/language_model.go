package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LanguageModel struct {
	gorm.Model

	Name datatypes.JSON
	Code string `gorm:"size:32;not null;uniqueIndex'"`
}

func (*LanguageModel) TableName() string {
	return "languages"
}
