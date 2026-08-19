package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DraftModel struct {
	gorm.Model

	UserID uint
	User   *UserModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Key  string `gorm:"size:1024;uniqueIndex;not null;"`
	Data datatypes.JSON
}

func NewDraftModel(userID uint, key string, data datatypes.JSON) *DraftModel {
	return &DraftModel{UserID: userID, Key: key, Data: data}
}

func (DraftModel) TableName() string {
	return "drafts"
}
