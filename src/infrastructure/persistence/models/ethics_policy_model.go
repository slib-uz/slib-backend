package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EthicsPolicyModel struct {
	gorm.Model

	Content datatypes.JSON `gorm:"not null"`
}

func (EthicsPolicyModel) TableName() string {
	return "ethics_policies"
}
