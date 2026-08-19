package repository

import (
	"context"

	"slib.uz/src/core/domain/entity"
)

type EditionRepository interface {
	Create(edition *entity.EditionEntity) error
	GetByID(ctx context.Context, id uint) (*entity.EditionEntity, error)
	GetByJournalID(ctx context.Context, journalID uint, page int, pageSize int, search string, year int) (*entity.PagingEntity[entity.EditionEntity], error)
	Update(ctx context.Context, edition *entity.EditionEntity) error
	Delete(ctx context.Context, editionID uint) error
	AttachArticles(ctx context.Context, editionID uint, articleIDs []uint) (int64, error)
	DetachArticles(ctx context.Context, editionID uint, articleIDs []uint) (int64, error)
}
