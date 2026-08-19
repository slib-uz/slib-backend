package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
)

func StudyFieldToResponse(entity *entity.StudyFieldEntity) *response.KafkaStudyFieldResponse {
	var parent *response.KafkaStudyFieldResponse
	if entity.Parent != nil {
		parent = StudyFieldToResponse(entity.Parent)
	}
	return response.NewKafkaStudyFieldResponse(entity.ID, entity.Name, entity.Code, entity.ParentID, parent)
}
