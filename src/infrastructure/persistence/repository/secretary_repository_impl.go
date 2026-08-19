package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type SecretaryRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewSecretaryRepository(repository *BaseRepository) repository.SecretaryRepository {
	return &SecretaryRepositoryImpl{BaseRepository: repository}
}

func (this *SecretaryRepositoryImpl) GetByJournalID(journalID uint) ([]*entity.SecretaryEntity, error) {
	var users []*models.UserModel

	if err := this.db().
		Joins("JOIN roles ON roles.user_id = users.id").
		Where("roles.journal_id = ? AND roles.role = ? AND roles.deleted_at IS NULL", journalID, enum.RoleSecretary).
		Find(&users).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.SecretaryListModelToEntity(users), nil
}
