package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RegionModel struct {
	gorm.Model

	Name  datatypes.JSON `gorm:"not null;"`
	Soato int            `gorm:"not null;unique;"`
}

func NewRegionModel(name datatypes.JSON, soato int) *RegionModel {
	return &RegionModel{Name: name, Soato: soato}
}

func (RegionModel) TableName() string {
	return "regions"
}
