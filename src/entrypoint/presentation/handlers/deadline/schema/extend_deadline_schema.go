package schema

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type ExtendSchemaRequest struct {
	ReviewStageID uint              `query:"review_stage_id"`
	Deadline      time.Time         `json:"deadline"`
	DeadlineType  enum.DeadlineType `json:"deadline_type"`
}
