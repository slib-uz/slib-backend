package schema

import (
	"slib.uz/src/core/domain/entity"
)

type StudyFieldUpdateRequest struct {
	ID       uint              `json:"id" validate:"required"`
	Name     map[string]string `json:"name" validate:"omitempty"`
	ParentID *uint             `json:"parent_id" validate:"omitempty"`
	Code     *uint             `json:"code" validate:"omitempty"`
}

func (this *StudyFieldUpdateRequest) ToEntity() *entity.StudyFieldEntity {
	return entity.NewStudyFieldEntity(
		this.ID,
		this.Name,
		this.ParentID,
		this.Code,
		nil,
	)
}
