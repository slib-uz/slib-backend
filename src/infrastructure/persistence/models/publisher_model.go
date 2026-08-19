package models

import (
	"gorm.io/gorm"
)

type PublisherModel struct {
	gorm.Model

	Tin         string  `gorm:"size:32;uniqueIndex"`
	Name        *string `gorm:"size:500"`
	ShortName   *string `gorm:"size:255"`
	Logo        *string `gorm:"size:500"`
	PhoneNumber *string `gorm:"size:128"`
	FaxNumber   *string `gorm:"size:128"`
	Email       *string `gorm:"size:255"`
	Website     *string `gorm:"size:64"`
	Address     *string `gorm:"size:500"`
	Description *string `gorm:"size:32512"`
	IsActive    bool

	InstitutionID *uint             `gorm:"index"`
	Institution   *InstitutionModel `gorm:"foreignKey:InstitutionID;constraint:OnDelete:SET NULL;"`

	RegionID   *uint          `gorm:"index"`
	Region     *RegionModel   `gorm:"foreignKey:RegionID;constraint:OnDelete:SET NULL;"`
	DistrictID *uint          `gorm:"index"`
	District   *DistrictModel `gorm:"foreignKey:DistrictID;constraint:OnDelete:SET NULL;"`
}

func (PublisherModel) TableName() string {
	return "publishers"
}

func NewPublisherModel(
	tin string,
	name *string,
	shortName *string,
	logo *string,
	phoneNumber *string,
	faxNumber *string,
	email *string,
	website *string,
	address *string,
	description *string,
	isActive bool,
) *PublisherModel {
	return &PublisherModel{
		Tin:         tin,
		Name:        name,
		ShortName:   shortName,
		Logo:        logo,
		PhoneNumber: phoneNumber,
		FaxNumber:   faxNumber,
		Email:       email,
		Website:     website,
		Address:     address,
		Description: description,
		IsActive:    isActive,
	}
}
