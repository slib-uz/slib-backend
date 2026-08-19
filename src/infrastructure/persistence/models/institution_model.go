package models

import "gorm.io/gorm"

type InstitutionModel struct {
	gorm.Model

	Name string  `gorm:"size:500;not null"`
	Tin  *string `gorm:"size:32"`
	Logo *string `gorm:"size:500"`

	Publishers []PublisherModel `gorm:"foreignKey:InstitutionID"`
}

func (InstitutionModel) TableName() string {
	return "institutions"
}

func NewInstitutionModel(name string, tin, logo *string) *InstitutionModel {
	return &InstitutionModel{
		Name: name,
		Tin:  tin,
		Logo: logo,
	}
}
