package schema

import (
	"slib.uz/src/core/domain/entity/enum"
)

type TechnicalReviewRequest struct {
	ApplicationID uint        `json:"application_id"`
	Status        enum.Status `json:"status"`
	Reason        *string     `json:"reason"`
}

type PeerReviewRequest struct {
	ApplicationID uint        `json:"application_id"`
	Status        enum.Status `json:"status"`
	Reason        *string     `json:"reason"`
}
