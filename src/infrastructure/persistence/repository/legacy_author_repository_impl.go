package repository

import (
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/models"
)

type LegacyAuthorRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewLegacyAuthorRepository(baseRepository *BaseRepository) repository.LegacyAuthorRepository {
	return &LegacyAuthorRepositoryImpl{BaseRepository: baseRepository}
}

func (this *LegacyAuthorRepositoryImpl) GetIDsByFullName(fullName string) ([]uint, error) {
	if len(fullName) == 0 {
		return []uint{}, nil
	}

	var ids []uint

	threshold := 0.4

	err := this.db().Model(&models.LegacyAuthorModel{}).
		Where("similarity(full_name, ?) > ?", fullName, threshold).
		Order(this.db().Raw("similarity(full_name, ?) DESC", fullName)).
		Limit(50).
		Pluck("id", &ids).
		Error

	if err != nil {
		return nil, infraError.Wrap(err)
	}

	return ids, nil
}
