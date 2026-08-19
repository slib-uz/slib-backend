package peerreviewusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type PeerReviewSubmissionReviewsUseCase struct {
	gateway    gateway.ETaqrizGateway
	repository repository.PeerReviewSubmissionRepository
}

// @inject
func NewPeerReviewSubmissionReviewsUseCase(gateway gateway.ETaqrizGateway, repository repository.PeerReviewSubmissionRepository) *PeerReviewSubmissionReviewsUseCase {
	return &PeerReviewSubmissionReviewsUseCase{gateway: gateway, repository: repository}
}

func (this *PeerReviewSubmissionReviewsUseCase) Execute(submissionExternalID uint) (*entity.PeerReviewSubmissionEntity, error) {
	submission, err := this.repository.GetByExternalID(submissionExternalID)
	if err != nil {
		return nil, err
	}

	detail, err := this.gateway.GetSubmissionDetail(submissionExternalID)
	if err != nil {
		return nil, err
	}

	submission.Reviews = detail.Reviews

	return submission, nil
}
