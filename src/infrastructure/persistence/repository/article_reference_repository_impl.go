package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/db"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ArticleReferenceRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewArticleReferenceRepository(db *db.Database) repository.ArticleReferenceRepository {
	return &ArticleReferenceRepositoryImpl{BaseRepository: NewBaseRepository(db)}
}

func (this *ArticleReferenceRepositoryImpl) GetListByArticleID(articleID uint) ([]*entity.ReferenceEntity, error) {
	var references []*models.ReferenceModel

	if err := this.db().Model(&models.ReferenceModel{}).
		Select("id", "name", "article_id").
		Where("article_id = ?", articleID).
		Find(&references).Error; err != nil {
		return nil, err
	}

	return mapper.ReferenceModelListToEntityList(references), nil
}
