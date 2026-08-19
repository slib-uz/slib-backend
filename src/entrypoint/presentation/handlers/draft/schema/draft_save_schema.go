package schema

import (
	"slib.uz/src/core/domain/entity"
)

type DraftSaveReqeust struct {
	Key  string         `json:"key" validate:"required"`
	Data map[string]any `json:"data" validate:"required"`
}

func (this *DraftSaveReqeust) ToEntity(userID uint) *entity.DraftEntity {
	return entity.NewDraftEntity(0, userID, this.Key, this.Data)
}
