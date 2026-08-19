package models

import "gorm.io/gorm"

type OrganizationModel struct {
	gorm.Model

	Tin     string  `gorm:"size:64;not null;uniqueIndex:idx_organizations_tin,where:deleted_at IS NULL"`
	Name    string  `gorm:"size:500; not null"`
	Address *string `gorm:"size:500"`

	SoatoID *uint                 `gorm:"index"`                                           // soato_id
	Soato   *SoatoClassifierModel `gorm:"foreignKey:SoatoID;constraint:OnDelete:SET NULL"` // many2one
}

func (OrganizationModel) TableName() string {
	return "organizations"
}

func NewOrganizationModel(tin, name string, address *string, soatoID *uint) *OrganizationModel {
	return &OrganizationModel{
		Tin:     tin,
		Name:    name,
		Address: address,
		SoatoID: soatoID,
	}
}
