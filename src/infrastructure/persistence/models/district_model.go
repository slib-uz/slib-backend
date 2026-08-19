package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DistrictModel struct {
	gorm.Model

	Name  datatypes.JSON `gorm:"not null;"`
	Soato int            `gorm:"not null;unique;"`

	RegionID uint         `gorm:"not null;"`
	Region   *RegionModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewDistrictModel(name datatypes.JSON, soato int, regionID uint) *DistrictModel {
	return &DistrictModel{Name: name, Soato: soato, RegionID: regionID}
}

func (DistrictModel) TableName() string {
	return "districts"
}
