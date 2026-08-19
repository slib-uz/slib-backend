package peerreviewusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type PeerReviewSubmissionsListUseCase struct {
	repository repository.PeerReviewSubmissionRepository
}

// @inject
func NewPeerReviewSubmissionsListUseCase(repository repository.PeerReviewSubmissionRepository) *PeerReviewSubmissionsListUseCase {
	return &PeerReviewSubmissionsListUseCase{repository: repository}
}

func (this *PeerReviewSubmissionsListUseCase) Execute(applicationID uint) ([]*entity.PeerReviewSubmissionEntity, error) {
	return this.repository.GetByApplicationID(applicationID)
}
