package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func PeerReviewSubmissionModelToEntity(s *models.PeerReviewSubmissionModel) *entity2.PeerReviewSubmissionEntity {
	var reviewer *entity2.ReviewerEntity
	var reviews []*entity2.PeerReviewResultEntity
	var sender *entity2.UserEntity
	var application *entity2.ApplicationEntity

	if s.Sender != nil {
		sender = UserModelToEntity(s.Sender)

	}

	if s.Reviewer != nil {
		reviewer = ReviewerModelToEntity(s.Reviewer)
	}

	if s.Application != nil {
		application = ApplicationModelToEntity(s.Application)
	}

	return entity2.NewPeerReviewSubmissionEntity(
		s.ID,
		s.ExternalIdempotencyKey,
		s.ExternalID,
		s.ApplicationID,
		application,
		s.ReviewerInternalID,
		reviewer,
		s.ReviewerExternalID,
		s.SenderTitle,
		s.Title,
		FromGormJson[map[string]any](s.ExtraData),
		s.OldDeadline,
		s.Deadline,
		s.Status,
		s.ReviewMethod,
		s.SenderID,
		sender,
		s.CreatedAt,
		reviews,
		s.Version,
	)
}

func PeerReviewSubmissionEntityToModel(s *entity2.PeerReviewSubmissionEntity) *models.PeerReviewSubmissionModel {

	return models.NewPeerReviewSubmissionModel(
		s.ExternalIdempotencyKey,
		s.ExternalID,
		s.ReviewerInternalID,
		s.ReviewerExternalID,
		s.ApplicationID,
		s.SenderTitle,
		s.Title,
		ToGormJson(s.ExtraData),
		s.OldDeadline,
		s.Deadline,
		s.Status,
		s.ReviewMethod,
		s.SenderID,
		s.Version,
	)
}
