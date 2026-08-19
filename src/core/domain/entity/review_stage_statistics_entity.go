package entity

import "slib.uz/src/core/domain/entity/enum"

type ReviewStageStatisticsEntity struct {
	StageName    string     `json:"stage_name"`
	StageNumber  enum.Stage `json:"stage_number" swaggertype:"integer"`
	Rejected     int64      `json:"rejected"`
	Pending      int64      `json:"pending"`
	Accepted     int64      `json:"accepted"`
	Total        int64      `json:"total"`
	OverdueCount int64      `json:"overdue"`
}

func NewReviewStageStatisticsEntity(stageName string, stageNumber enum.Stage, rejected, pending, accepted, total, overdueCount int64) *ReviewStageStatisticsEntity {
	return &ReviewStageStatisticsEntity{
		StageName:    stageName,
		StageNumber:  stageNumber,
		Rejected:     rejected,
		Pending:      pending,
		Accepted:     accepted,
		Total:        total,
		OverdueCount: overdueCount,
	}
}
