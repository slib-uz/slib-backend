package schema

import (
	"slib.uz/src/core/domain/entity/enum"
)

type ReviewApplicationRequest struct {
	ApplicationID   uint        `json:"application_id" validate:"required"`
	Status          enum.Status `json:"status" validate:"required"` // -1: rejected, 1: accepted
	RejectionReason *string     `json:"rejection_reason"`
}
