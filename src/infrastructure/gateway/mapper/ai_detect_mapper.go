package mapper

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/gateway/response"
)

func AiDetectResultResponseToEntity(data *response.AiDetectResultResponse) *entity.AiDetectResultEntity {
	var createdAt *time.Time

	if data.CreatedAt != "" {
		parsedTime, err := time.Parse(time.DateTime, data.CreatedAt)
		if err == nil {
			createdAt = &parsedTime
		}
	}

	return entity.NewAiDetectResultEntity(
		0,
		0,
		nil,
		0,
		nil,
		0,
		nil,
		0,
		nil,
		data.ExternalID,
		data.WordsCount,
		enum.AntiPlagStatus(data.Status),
		data.StatusDisplay,
		data.HumanPercent,
		data.WarnPercent,
		data.AiPercent,
		data.ReportURL,
		createdAt,
		nil,
	)
}
