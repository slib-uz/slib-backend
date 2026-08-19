package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type AuthorshipClaimRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewAuthorshipClaimRepository(repository *BaseRepository) repository.AuthorshipClaimRepository {
	return &AuthorshipClaimRepositoryImpl{BaseRepository: repository}
}

func (this *AuthorshipClaimRepositoryImpl) Create(claim *entity.AuthorshipClaimEntity) error {
	model := mapper.AuthorshipClaimEntityToModel(claim)
	if err := this.db().Create(model).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}

func (this *AuthorshipClaimRepositoryImpl) CreateBatch(claims []*entity.AuthorshipClaimEntity) error {
	_models := make([]*models.AuthorshipClaimModel, len(claims))
	for i, claim := range claims {
		_models[i] = mapper.AuthorshipClaimEntityToModel(claim)
	}

	return infraError.Wrap(this.db().CreateInBatches(_models, 100).Error)
}

func (this *AuthorshipClaimRepositoryImpl) GetList(page, size int, filters map[string]interface{}) (*entity.PagingEntity[entity.AuthorshipClaimEntity], error) {
	var claims []*models.AuthorshipClaimModel
	query := this.db().Model(&models.AuthorshipClaimModel{}).
		Preload("Sender").
		Preload("Article").
		Preload("Article.Journal").
		Preload("Article.Journal.Publisher").
		Preload("ReviewedBy")

	if val, ok := filters["article_id"]; ok {
		query = query.Where("article_id = ?", val)
	}

	if val, ok := filters["status"]; ok {
		query = query.Where("status = ?", val)
	}

	if val, ok := filters["journal_id"]; ok {
		query = query.Joins("JOIN articles ON articles.id = authorship_claims.article_id").
			Where("articles.journal_id = ?", val)
	}

	if val, ok := filters["publisher_id"]; ok {
		query = query.
			Joins("JOIN articles ON articles.id = authorship_claims.article_id").
			Joins("JOIN journals ON journals.id = articles.journal_id").
			Where("journals.publisher_id = ?", val)
	}

	if val, ok := filters["sender_id"]; ok {
		query = query.Where("sender_id = ?", val)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := query.Limit(size).Offset((page - 1) * size).Order("created_at desc").Find(&claims).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	entities := mapper.AuthorshipClaimModelListToEntityList(claims)
	return entity.NewPagingEntity(page, size, total, entities), nil
}

func (this *AuthorshipClaimRepositoryImpl) GetByID(id uint) (*entity.AuthorshipClaimEntity, error) {
	var claim models.AuthorshipClaimModel
	err := this.db().
		Preload("Sender").
		Preload("Article").
		Preload("Article.Journal").
		Preload("Article.Journal.Publisher").
		Preload("ReviewedBy").
		First(&claim, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return standard Not Found error
		}
		return nil, infraError.Wrap(err)
	}
	return mapper.AuthorshipClaimModelToEntity(&claim), nil
}

func (this *AuthorshipClaimRepositoryImpl) FindPendingByArticleAndUser(articleID, userID uint) (*entity.AuthorshipClaimEntity, error) {
	var claim models.AuthorshipClaimModel
	err := this.db().
		Where("article_id = ?", articleID).
		Where("sender_id = ?", userID).
		Where("status = ?", "pending").
		First(&claim).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, infraError.Wrap(err)
	}
	return mapper.AuthorshipClaimModelToEntity(&claim), nil
}

func (this *AuthorshipClaimRepositoryImpl) FindPendingByArticleIDsAndUser(articleIDs []uint, userID uint) ([]*entity.AuthorshipClaimEntity, error) {
	var claims []*models.AuthorshipClaimModel
	err := this.db().
		Where("article_id IN ?", articleIDs).
		Where("sender_id = ?", userID).
		Where("status = ?", "pending").
		Find(&claims).Error

	if err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.AuthorshipClaimModelListToEntityList(claims), nil
}

func (this *AuthorshipClaimRepositoryImpl) Update(claim *entity.AuthorshipClaimEntity) error {
	model := mapper.AuthorshipClaimEntityToModel(claim)
	if err := this.db().Model(&models.AuthorshipClaimModel{}).Where("id = ?", model.ID).Updates(model).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}
