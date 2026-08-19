package mapper

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/gateway/response"
)

func AntiPlagResultResponseToEntity(data *response.AntiPlagResultResponse) *entity.AntiPlagResultEntity {

	var createAt *time.Time

	if data.ExternalCreatedAt != "" {
		parsedTime, err := time.Parse(time.DateTime, data.ExternalCreatedAt)
		if err == nil {
			createAt = &parsedTime
		}
	}

	return entity.NewAntiPlagResultEntity(
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
		enum.AntiPlagStatus(data.Status),
		data.StatusDisplay,
		data.PlagiarismPercent,
		data.LegalPercent,
		data.SelfCitePercent,
		data.UnknownPercent,
		data.ShortReportURL,
		data.FullReportURL,
		createAt,
		nil,
		nil,
	)
}
