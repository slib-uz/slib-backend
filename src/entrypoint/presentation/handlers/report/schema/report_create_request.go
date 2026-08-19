package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ReportCreateRequest struct {
	Reason     string          `json:"reason"`
	TargetID   uint            `json:"target_id"`
	TargetType enum.ReportType `json:"target_type"`
}

func (this ReportCreateRequest) ToEntity() *entity.ReportEntity {
	return entity.NewReportEntity(
		0,
		this.Reason,
		this.TargetID,
		nil,
		this.TargetType,
		0,
		nil,
		time.Now())
}
