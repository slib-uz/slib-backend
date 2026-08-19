package schema

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type UpdateStatusRequest struct {
	Status       enum.ClaimStatus `json:"status"`
	RejectReason string           `json:"reject_reason"`
}

func (this *UpdateStatusRequest) ToEntity(claimID uint, reviewerID uint) *entity.UpdateAuthorshipClaimStatusInputEntity {
	return entity.NewUpdateAuthorshipClaimStatusInputEntity(claimID, this.Status, this.RejectReason, reviewerID)
}

func (this *UpdateStatusRequest) Validate() (bool, error) {
	switch this.Status {
	case enum.ClaimStatusApproved, enum.ClaimStatusRejected:
		return true, nil
	default:
		return false, response.NewFailResponse(400, "invalid status value")
	}
}
