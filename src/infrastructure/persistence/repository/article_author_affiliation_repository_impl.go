package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slib.uz/src/core/domain/entity"

	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ArticleAuthorAffiliationRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewArticleAuthorAffiliationRepository(baseRepository *BaseRepository) repository.ArticleAuthorAffiliationRepository {
	return &ArticleAuthorAffiliationRepositoryImpl{BaseRepository: baseRepository}
}

func (this *ArticleAuthorAffiliationRepositoryImpl) Create(articleAuthorAffiliation *entity.ArticleAuthorAffiliationEntity) (*entity.ArticleAuthorAffiliationEntity, error) {
	articleAuthorAffiliationModel := mapper.ArticleAuthorAffiliationEntityToModel(articleAuthorAffiliation)
	if err := this.db().Create(articleAuthorAffiliationModel).Error; err != nil {
		return nil, err
	}
	return mapper.ArticleAuthorAffiliationModelToEntity(articleAuthorAffiliationModel), nil
}

func (this *ArticleAuthorAffiliationRepositoryImpl) GetByArticleID(articleID uint) ([]*entity.ArticleAuthorAffiliationEntity, error) {
	var articleAuthorAffiliations []*models.ArticleAuthorAffiliationModel
	if err := this.db().Where("article_id = ?", articleID).Find(&articleAuthorAffiliations).Error; err != nil {
		return nil, err
	}
	return mapper.ArticleAuthorAffiliationListModelToEntity(articleAuthorAffiliations), nil
}

func (this *ArticleAuthorAffiliationRepositoryImpl) GetAuthorIdsByArticleAuthorAffiliationIDs(articleAuthorAffiliationIDs []uint) ([]uint, error) {

	var authorIds []uint
	if err := this.db().Model(&models.ArticleAuthorAffiliationModel{}).Where("id IN ?", articleAuthorAffiliationIDs).Pluck("author_id", &authorIds).Error; err != nil {
		return nil, err
	}
	return authorIds, nil
}

func (this *ArticleAuthorAffiliationRepositoryImpl) ExistingIds(ids []uint) ([]uint, error) {
	var articleAuthorAffiliationIds []uint
	if err := this.db().Model(&models.ArticleAuthorAffiliationModel{}).Where("id IN ?", ids).Pluck("id", &articleAuthorAffiliationIds).Error; err != nil {
		return nil, err
	}
	return articleAuthorAffiliationIds, nil
}

func (this *ArticleAuthorAffiliationRepositoryImpl) BulkUpdateArticleIds(articleID uint, affiliationIDs []uint) ([]uint, error) {
	return this.BulkUpdateArticleIdsWithTransaction(this.db(), articleID, affiliationIDs)
}

// BulkUpdateArticleIdsWithTransaction updates article IDs for affiliations using provided transaction
func (this *ArticleAuthorAffiliationRepositoryImpl) BulkUpdateArticleIdsWithTransaction(tx *gorm.DB, articleID uint, affiliationIDs []uint) ([]uint, error) {
	var authorIDs []uint

	err := tx.
		Model(&models.ArticleAuthorAffiliationModel{}).
		Where("id IN ?", affiliationIDs).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "author_id"}}}).
		UpdateColumn("article_id", articleID).
		Pluck("author_id", &authorIDs).Error

	if err != nil {
		return nil, infraError.Wrap(err)
	}

	return authorIDs, nil
}

func (this *ArticleAuthorAffiliationRepositoryImpl) BulkCreateForAuthorsNeedArticleAuthorAffiliation(affiliations []*entity.ArticleAuthorAffiliationEntity) error {
	return this.BulkCreateForAuthorsNeedArticleAuthorAffiliationWithTransaction(this.db(), affiliations)
}

// BulkCreateForAuthorsNeedArticleAuthorAffiliationWithTransaction creates affiliations using provided transaction
func (this *ArticleAuthorAffiliationRepositoryImpl) BulkCreateForAuthorsNeedArticleAuthorAffiliationWithTransaction(tx *gorm.DB, affiliations []*entity.ArticleAuthorAffiliationEntity) error {
	articleAuthorAffiliationModels := make([]*models.ArticleAuthorAffiliationModel, len(affiliations))
	for i, affiliation := range affiliations {
		articleAuthorAffiliationModels[i] = mapper.ArticleAuthorAffiliationEntityToModel(affiliation)
	}
	if err := tx.Create(articleAuthorAffiliationModels).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}

func (this *ArticleAuthorAffiliationRepositoryImpl) GetListByAuthorIDs(authorIDs []uint) ([]*entity.ArticleAuthorAffiliationEntity, error) {
	var _models []*models.ArticleAuthorAffiliationModel
	if err := this.db().Where("author_id IN ?", authorIDs).Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.ArticleAuthorAffiliationListModelToEntity(_models), nil
}
