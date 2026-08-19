package schema

import (
	"slib.uz/src/core/domain/entity"
)

type StudyFieldCreateRequest struct {
	Name     map[string]string `json:"name" validate:"required"`
	ParentID *uint             `json:"parent_id" validate:"omitempty"`
	Code     *uint             `json:"code" validate:"omitempty"`
}

func (this *StudyFieldCreateRequest) ToEntity() *entity.StudyFieldEntity {
	return entity.NewStudyFieldEntity(
		0,
		this.Name,
		this.ParentID,
		this.Code,
		nil,
	)
}
