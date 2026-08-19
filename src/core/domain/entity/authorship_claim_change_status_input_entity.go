package entity

import "slib.uz/src/core/domain/entity/enum"

type UpdateAuthorshipClaimStatusInputEntity struct {
	ClaimID      uint             `json:"claim_id"`
	Status       enum.ClaimStatus `json:"status"`
	RejectReason string           `json:"reject_reason"`
	ReviewerID   uint             `json:"reviewer_id"`
}

func NewUpdateAuthorshipClaimStatusInputEntity(claimID uint, status enum.ClaimStatus, rejectReason string, reviewerID uint) *UpdateAuthorshipClaimStatusInputEntity {
	return &UpdateAuthorshipClaimStatusInputEntity{
		ClaimID:      claimID,
		Status:       status,
		RejectReason: rejectReason,
		ReviewerID:   reviewerID,
	}
}
