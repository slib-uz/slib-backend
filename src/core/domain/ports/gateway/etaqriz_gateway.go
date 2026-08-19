package gateway

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ETaqrizGateway interface {
	FindReviewerByScienceID(id string) (*entity2.ReviewerEntity, error)
	ExtendDeadline(submissionExternalID uint, deadline time.Time) error
	Resubmit(submissionExternalID uint, fileBase64 string) (deadline time.Time, _ error)
	GetSubmissionDetail(submissionExternalID uint) (*entity2.PeerReviewSubmissionEntity, error)
	Submit(
		idempotencyKey string,
		applicationID uint,
		authorFullName string,
		authorScienceID string,
		senderID uint,
		title,
		senderTitle string,
		reviewers []*entity2.ReviewerEntity,
		reviewMethod enum.PeerReviewMethod,
		file string,
	) ([]*entity2.PeerReviewSubmissionEntity, error)
}
