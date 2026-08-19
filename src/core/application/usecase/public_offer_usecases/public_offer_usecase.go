package public_offer_usecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type PublicOfferUseCase struct {
	repository repository.PublicOfferRepository
}

// @inject
func NewPublicOfferUseCase(repository repository.PublicOfferRepository) *PublicOfferUseCase {
	return &PublicOfferUseCase{repository: repository}
}

func (this *PublicOfferUseCase) Execute() (*entity.PublicOfferEntity, error) {
	publicOffer, err := this.repository.Get()

	if err != nil {
		return nil, err
	}

	return publicOffer, nil
}
