package repository

import (
	"fmt"
	"time"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ArticlePurchaseRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewArticlePurchaseRepository(repository *BaseRepository) repository.ArticlePurchaseRepository {
	return &ArticlePurchaseRepositoryImpl{BaseRepository: repository}
}

func (this *ArticlePurchaseRepositoryImpl) Create(_purchase *entity2.ArticlePurchaseEntity) error {

	var purchase = mapper.ArticlePurchaseEntityToModel(_purchase)
	return this.db().Create(&purchase).Error
}

func (this *ArticlePurchaseRepositoryImpl) IsExists(articleID, userID uint) (bool, error) {
	var count int64
	if err := this.db().
		Model(&models.ArticlePurchaseModel{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (this *ArticlePurchaseRepositoryImpl) StatsByJournal(publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalArticlePurchaseStatsEntity], error) {
	var stats []*entity2.JournalArticlePurchaseStatsEntity

	query := this.db().Table("article_purchases").
		Select("article_purchases.journal_id as journal_id, journals.name as journal_name, COUNT(article_purchases.id) AS count, SUM(article_purchases.amount) AS total_amount").
		Joins("JOIN journals ON article_purchases.journal_id = journals.id").
		Where("article_purchases.created_at BETWEEN ? AND ?", startDate, endDate)

	if publisherID > 0 {
		query = query.Where("journals.publisher_id = ?", publisherID)
	}
	query = query.Group("article_purchases.journal_id, journals.name")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("error counting journal stats: %w", err)
	}
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("error fetching journal stats: %w", err)
	}

	return entity2.NewPagingEntity(page, pageSize, total, stats), nil
}

func (this *ArticlePurchaseRepositoryImpl) GetByJournalID(journalID uint, page, pageSize int, startDate, endDate time.Time) (*entity2.PagingEntity[entity2.ArticlePurchaseEntity], error) {
	var purchases []*models.ArticlePurchaseModel

	query := this.db().
		Model(&models.ArticlePurchaseModel{}).
		Where("journal_id = ? AND created_at BETWEEN ? AND ?", journalID, startDate, endDate).
		Order("created_at DESC")

	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("error counting article purchases: %w", err)
	}

	// Apply pagination
	if err := query.
		Preload("Article").
		Preload("User").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&purchases).Error; err != nil {
		return nil, fmt.Errorf("error fetching article purchases: %w", err)
	}

	return mapper.PagingMapper(purchases, mapper.ArticlePurchaseModelToEntity, page, pageSize, total), nil
}

func (this *ArticlePurchaseRepositoryImpl) GetByUserID(userID uint, page, pageSize int, search string) (*entity2.PagingEntity[entity2.ArticlePurchaseEntity], error) {
	var purchases []*models.ArticlePurchaseModel

	query := this.db().
		Model(&models.ArticlePurchaseModel{}).
		Where("user_id = ?", userID)

	if search != "" {
		query = query.Joins("JOIN articles ON article_purchases.article_id = articles.id").
			Where(
				"(articles.name->>'uz' ILIKE ? OR articles.name->>'ru' ILIKE ? OR articles.name->>'en' ILIKE ?)",
				"%"+search+"%", "%"+search+"%", "%"+search+"%",
			)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("error counting article purchases: %w", err)
	}

	// Apply pagination
	if err := query.
		Preload("Article").
		Preload("Journal").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&purchases).Error; err != nil {
		return nil, fmt.Errorf("error fetching article purchases: %w", err)
	}

	return mapper.PagingMapper(purchases, mapper.ArticlePurchaseModelToEntity, page, pageSize, total), nil

}
