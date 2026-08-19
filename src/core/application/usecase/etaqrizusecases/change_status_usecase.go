package etaqrizusecases

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ChangeStatusUseCase struct {
	submissionRepository repository.PeerReviewSubmissionRepository
}

// @inject
func NewChangeStatusUseCase(submissionRepository repository.PeerReviewSubmissionRepository) *ChangeStatusUseCase {
	return &ChangeStatusUseCase{submissionRepository: submissionRepository}
}

func (this *ChangeStatusUseCase) Execute(submissionExternalID uint, status enum.PeerReviewStatus, deadline time.Time, oldDeadline *time.Time, version int64) error {

	submission, err := this.submissionRepository.GetByExternalID(submissionExternalID)
	if err != nil {
		return err
	}

	if submission.Version >= version {
		return nil
	}

	return this.submissionRepository.UpdateStatusAndDeadline(submissionExternalID, status, deadline, oldDeadline, version)

}
