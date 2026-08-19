package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
)

func JournalEntityToResponse(entity *entity.JournalEntity) *response.KafkaJournalResponse {
	return response.NewKafkaJournalResponse(entity.ID, entity.Name, entity.ShortName, entity.ISSNPaper, entity.ISSNOnline)
}
