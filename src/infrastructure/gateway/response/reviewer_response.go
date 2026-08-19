package response

import (
	"slib.uz/src/core/domain/entity"
)

type FindReviewerResponse struct {
	Success bool                  `json:"success"`
	Data    entity.ReviewerEntity `json:"data"`
	Message string                `json:"message"`
}

type ReviewerResponse struct {
	UserId      uint             `json:"user_id"`
	ScienceId   string           `json:"science_id"`
	FullName    string           `json:"full_name"`
	PhoneNumber string           `json:"phone_number"`
	Subjects    []map[string]any `json:"subjects"`
}

func (this *ReviewerResponse) ToEntity() *entity.ReviewerEntity {
	return entity.NewReviewerEntity(
		0,
		this.UserId,
		this.ScienceId,
		this.FullName,
		this.PhoneNumber,
		this.Subjects,
	)
}
