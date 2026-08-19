package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ChiefEditorRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewChiefEditorRepository(baseRepository *BaseRepository) repository.ChiefEditorRepository {
	return &ChiefEditorRepositoryImpl{BaseRepository: baseRepository}
}

func (this *ChiefEditorRepositoryImpl) GetByJournalID(journalID uint) ([]*entity.ChiefEditorEntity, error) {
	var users []*models.UserModel

	if err := this.db().
		Joins("JOIN roles ON roles.user_id = users.id").
		Where("roles.journal_id = ? AND roles.role = ? AND roles.deleted_at IS NULL", journalID, enum.RoleChiefEditor).
		Find(&users).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.ChiefEditorListModelToEntity(users), nil
}
