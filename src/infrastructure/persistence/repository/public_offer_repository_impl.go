package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type PublicOfferRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewPublicOfferRepository(repository *BaseRepository) repository.PublicOfferRepository {
	return &PublicOfferRepositoryImpl{BaseRepository: repository}
}

func (this *PublicOfferRepositoryImpl) Get() (*entity.PublicOfferEntity, error) {
	var _model models.PublicOfferModel

	if err := this.db().Order("created_at desc").First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.PublicOfferModelToEntity(&_model), nil
}
