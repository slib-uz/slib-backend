package etaqrizusecases

import (
	"errors"

	core "slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type FindReviewerUseCase struct {
	gateway    gateway.ETaqrizGateway
	repository repository.ReviewerRepository
}

// @inject
func NewFindReviewerUseCase(gateway gateway.ETaqrizGateway, repository repository.ReviewerRepository) *FindReviewerUseCase {
	return &FindReviewerUseCase{gateway: gateway, repository: repository}
}

func (this *FindReviewerUseCase) Execute(scienceID string) (*entity.ReviewerEntity, error) {

	reviewer, err := this.repository.GetByScienceID(scienceID)
	if err == nil {
		return reviewer, nil
	}
	if !errors.Is(err, core.NotFoundError) {
		return nil, err
	}

	reviewer, err = this.gateway.FindReviewerByScienceID(scienceID)
	if err != nil {
		return nil, err
	}

	id, err := this.repository.Create(reviewer)
	if err != nil {
		return nil, err
	}
	reviewer.ID = id

	return reviewer, nil
}
