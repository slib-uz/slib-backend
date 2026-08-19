package entity

import "slib.uz/src/core/domain/entity/enum"

type InstructionEntity struct {
	ID          uint                `json:"id"`
	Key         enum.InstructionKey `json:"key"`
	VideoLink   *string             `json:"video_link"`
	Description *string             `json:"description"`
}
