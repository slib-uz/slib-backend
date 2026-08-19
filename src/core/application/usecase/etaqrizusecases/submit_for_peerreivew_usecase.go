package etaqrizusecases

import (
	"fmt"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type SubmitForPeerReviewUseCase struct {
	gateway               gateway.ETaqrizGateway
	applicationRepository repository.ApplicationRepository
	storage               storage.FileStorage
	reviewerRepository    repository.ReviewerRepository
	submissionRepository  repository.PeerReviewSubmissionRepository
}

// @inject
func NewSubmitForPeerReviewUseCase(gateway gateway.ETaqrizGateway, applicationRepository repository.ApplicationRepository, storage storage.FileStorage, reviewerRepository repository.ReviewerRepository, submissionRepository repository.PeerReviewSubmissionRepository) *SubmitForPeerReviewUseCase {
	return &SubmitForPeerReviewUseCase{gateway: gateway, applicationRepository: applicationRepository, storage: storage, reviewerRepository: reviewerRepository, submissionRepository: submissionRepository}
}

func (this *SubmitForPeerReviewUseCase) Execute(idempotencyKey string, userID uint, applicationID uint, reviewerIds []uint) error {

	exists, err := this.submissionRepository.IsExistsByExternalIdempotencyKey(idempotencyKey)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	app, err := this.applicationRepository.GetWithArticleAndJournal(applicationID)
	if err != nil {
		return err
	}

	file, err := this.storage.GetObjectAsBase64(enum.BucketPrivate, app.Article.ContentFile)
	if err != nil {
		return err
	}

	reviewers, err := this.getReviewers(reviewerIds)
	if err != nil {
		return err
	}

	submissions, err := this.gateway.Submit(
		idempotencyKey,
		applicationID,
		app.User.FullName,
		app.User.ScienceID,
		userID,
		app.Article.Name["uz"],
		app.Journal.Name["uz"],
		reviewers,
		app.Journal.PeerReviewMethod,
		file,
	)
	if err != nil {
		return fmt.Errorf("failed to submit for peer review: %w", err)

	}

	return this.submissionRepository.BulkCreate(submissions)
}

func (this *SubmitForPeerReviewUseCase) getReviewers(ids []uint) ([]*entity.ReviewerEntity, error) {
	reviewers, err := this.reviewerRepository.ListByIDs(ids)
	if err != nil {
		return nil, err
	}

	if len(reviewers) != len(ids) {
		return nil, response.NewFailResponse(400, "Some reviewers not found. Please check the reviewer IDs and try again.")
	}

	return reviewers, nil
}
