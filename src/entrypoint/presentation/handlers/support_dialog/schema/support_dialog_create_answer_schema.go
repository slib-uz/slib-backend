package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type SupportDialogCreateAnswerRequest struct {
	Message string `json:"message"`
	ChatID  uint   `json:"chat_id"`
}

func (this SupportDialogCreateAnswerRequest) ToEntity() *entity.SupportDialogEntity {
	return entity.NewSupportDialogEntity(
		0,
		enum.Answer,
		0,
		nil,
		this.Message,
		this.ChatID,
		true,
		time.Now())
}
