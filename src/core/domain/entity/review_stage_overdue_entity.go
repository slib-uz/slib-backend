package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type ReviewStageOverdueEntity struct {
	ReviewStageID     uint              `json:"review_stage_id"`
	ApplicationID     uint              `json:"application_id"`
	ApplicationNumber string            `json:"application_number"`
	ArticleID         uint              `json:"article_id"`
	Article           *ArticleEntity    `json:"article"`
	StageNumber       enum.Stage        `json:"stage_number" swaggertype:"integer"`
	StageName         string            `json:"stage_name"`
	Deadline          time.Time         `json:"deadline"`
	CreatedAt         time.Time         `json:"created_at"`
	OverdueDays       int               `json:"overdue_days"`
	JournalID         uint              `json:"journal_id"`
	JournalName       map[string]string `json:"journal_name"`
}
