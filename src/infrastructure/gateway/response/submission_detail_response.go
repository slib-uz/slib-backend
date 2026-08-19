package response

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type SubmissionDetailResponse struct {
	Data struct {
		Id          int            `json:"id"`
		ContentType string         `json:"content_type"`
		ExtraData   map[string]any `json:"extra_data"`
		File        string         `json:"file"`
		FileUrl     string         `json:"file_url"`

		OldDeadline  *string `json:"old_deadline"`
		Deadline     string  `json:"deadline"`
		ReviewMethod int     `json:"review_method"`

		Reviews     []*ReviewResponse     `json:"reviews"`
		SenderTitle string                `json:"sender_title"`
		Status      enum.PeerReviewStatus `json:"status"`
		Title       string                `json:"title"`
	} `json:"data"`
}

func (this *SubmissionDetailResponse) ToEntity() *entity2.PeerReviewSubmissionEntity {
	s := this.Data

	var deadline time.Time
	var oldDeadline *time.Time
	var reviews []*entity2.PeerReviewResultEntity

	if s.Deadline != "" {
		d := parseTime(s.Deadline)
		if d != nil {
			deadline = *d
		}
	}

	if s.OldDeadline != nil && *s.OldDeadline != "" {
		oldDeadline = parseTime(*s.OldDeadline)
	}

	if s.Reviews != nil {
		reviews = make([]*entity2.PeerReviewResultEntity, len(s.Reviews))
		for i, r := range s.Reviews {
			reviews[i] = r.ToEntity()
		}
	}
	return &entity2.PeerReviewSubmissionEntity{Status: s.Status, Deadline: deadline, OldDeadline: oldDeadline, Reviews: reviews}
}

type ReviewResponse struct {
	ID                uint                  `json:"id"`
	Comment           string                `json:"comment"`
	CreatedAt         string                `json:"created_at"`
	ReportFileUrl     string                `json:"report_file_url"`
	SubmissionFileUrl string                `json:"submission_file_url"`
	Review            map[string]string     `json:"review"`
	Status            enum.PeerReviewStatus `json:"status"`
	Submission        string                `json:"submission"`
	SubmissionID      uint                  `json:"submission_id"`
}

func (this *ReviewResponse) ToEntity() *entity2.PeerReviewResultEntity {

	var createdAt time.Time
	d := parseTime(this.CreatedAt)
	if d != nil {
		createdAt = *d
	}

	return entity2.NewPeerReviewResultEntity(
		this.ID,
		this.SubmissionID,
		this.Status,
		this.Review,
		this.Comment,
		this.ReportFileUrl,
		this.SubmissionFileUrl,
		createdAt,
	)
}
