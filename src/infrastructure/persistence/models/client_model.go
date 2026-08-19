package models

import "gorm.io/gorm"

type ClientModel struct {
	gorm.Model

	ClientID     string `gorm:"size:255;not null;uniqueIndex;"`
	ClientSecret string `gorm:"size:255;not null"`

	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	CallbackUrl string `gorm:"type:text"`

	JournalID *uint           `gorm:"index;"`
	Journal   *JournalModel   `gorm:"foreignKey:JournalID;constraint:OnDelete:SET NULL;"`

	IsActive bool `gorm:"default:true"`
}

func NewClientModel(name string, description string, callbackUrl string, clientID string, clientSecret string, journalID *uint, isActive bool) *ClientModel {
	return &ClientModel{
		Name:         name,
		Description:  description,
		CallbackUrl:  callbackUrl,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		JournalID:    journalID,
		IsActive:     isActive,
	}
}

func (ClientModel) TableName() string {
	return "clients"
}
