package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PublicOfferModel struct {
	gorm.Model

	Description datatypes.JSON `gorm:"not null"`
}

func NewPublicOfferModel(id uint, description datatypes.JSON) *PublicOfferModel {
	return &PublicOfferModel{Model: gorm.Model{ID: id}, Description: description}
}

func (PublicOfferModel) TableName() string { return "public_offers" }
