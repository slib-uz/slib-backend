package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type SupportDialogCreateQuestionRequest struct {
	Message string `json:"message"`
}

func (this SupportDialogCreateQuestionRequest) ToEntity() *entity.SupportDialogEntity {
	return entity.NewSupportDialogEntity(
		0,
		enum.Question,
		0,
		nil,
		this.Message,
		0,
		false,
		time.Now())
}
