package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type ReportEntity struct {
	ID         uint            `json:"id"`
	Reason     string          `json:"reason"`
	TargetID   uint            `json:"target_id"`
	Target     interface{}     `json:"target"`
	TargetType enum.ReportType `json:"target_type"`
	ReporterID uint            `json:"reporter_id"`
	Reporter   *ReporterEntity `json:"reporter"`
	CreatedAt  time.Time       `json:"created_at"`
}

func NewReportEntity(ID uint, reason string, targetID uint, target interface{}, targetType enum.ReportType, reporterID uint, reporter *ReporterEntity, createdAt time.Time) *ReportEntity {
	return &ReportEntity{ID: ID, Reason: reason, TargetID: targetID, Target: target, TargetType: targetType, ReporterID: reporterID, Reporter: reporter, CreatedAt: createdAt}
}
