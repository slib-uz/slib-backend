package models

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type InstructionModel struct {
	gorm.Model

	Key         enum.InstructionKey `gorm:"size:64;not null;uniqueIndex:idx_key_deleted_at,where:deleted_at IS NULL"`
	VideoLink   *string             `gorm:"size:1024"`
	Description *string             `gorm:"type:text"`
}

func (InstructionModel) TableName() string {
	return "instructions"
}
