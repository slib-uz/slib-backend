package entity

import "slib.uz/src/core/domain/entity/enum"

type ReviewStageOverdueStatisticsEntity struct {
	StageName    string     `json:"stage_name"`
	StageNumber  enum.Stage `json:"stage_number" swaggertype:"integer"`
	OverdueCount int64      `json:"overdue_count"`
}
