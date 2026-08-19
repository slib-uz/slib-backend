package authorusecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type AuthorByScienceIDUseCase struct {
	repository repository.AuthorRepository
	gateway    gateway.ScienceIDGateway
}

// @inject
func NewAuthorByScienceIDUseCase(repository repository.AuthorRepository, gateway gateway.ScienceIDGateway) *AuthorByScienceIDUseCase {
	return &AuthorByScienceIDUseCase{repository: repository, gateway: gateway}
}

func (this *AuthorByScienceIDUseCase) Execute(scienceID string) (*entity.AuthorEntity, error) {
	author, err := this.repository.GetByScienceID(scienceID)

	if err != nil {
		if errors.Is(err, response.NotFoundError) {
			author, err = this.gateway.GetAuthorByScienceID(scienceID)

			if err != nil {
				return nil, err
			}

			author, err = this.repository.SaveAuthor(author)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return author, nil
}
