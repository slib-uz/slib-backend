package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
)

func LanguageEntityToResponse(entity *entity.LanguageEntity) *response.LanguageResponse {
	return response.NewLanguageResponse(entity.ID, entity.Name, entity.Code)
}
