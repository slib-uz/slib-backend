package repository

import (
	"errors"

	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
	"slib.uz/src/infrastructure/persistence/sorting"
)

type JournalRatingRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewJournalRatingRepository(baseRepository *BaseRepository) repository.JournalRatingRepository {
	return &JournalRatingRepositoryImpl{BaseRepository: baseRepository}
}

func (this *JournalRatingRepositoryImpl) Create(rating *entity2.JournalRatingEntity) error {
	var _model = mapper.JournalRatingEntityToModel(rating)
	var exists models.JournalRatingModel
	var journal models.JournalModel

	if err := this.db().Model(&models.JournalModel{}).Where("id = ?", _model.JournalID).First(&journal).Error; err != nil {
		return infraError.Wrap(err)
	}

	// Check if user already has a non-deleted rating for this journal
	if err := this.db().Model(&models.JournalRatingModel{}).
		Where("user_id = ? AND journal_id = ? AND deleted_at IS NULL", _model.UserID, _model.JournalID).
		First(&exists).Error; err == nil && exists.ID != 0 {
		return errors.New("journal rating already exists")
	}

	if err := this.db().Create(&_model).Error; err != nil {
		return infraError.Wrap(err)
	}

	if err := this.db().Model(&models.JournalModel{}).
		Where("id = ?", _model.JournalID).
		UpdateColumns(map[string]interface{}{
			"rating_count": gorm.Expr("rating_count + ?", 1),
			"rating_sum":   gorm.Expr("rating_sum + ?", _model.Stars),
		}).Error; err != nil {
		return infraError.Wrap(err)
	}

	return nil
}

// JournalRatingSortFields — jurnal reytinglari ro'yxati uchun ruxsat etilgan
// tartiblash maydonlari.
var JournalRatingSortFields = sorting.New("journal_ratings.created_at DESC", map[string]string{
	"created_at": "journal_ratings.created_at",
})

func (this *JournalRatingRepositoryImpl) GetByJournalID(journalID uint, page, pageSize int, ordering string) (*entity2.PagingEntity[entity2.JournalRatingEntity], error) {
	orderExpr, err := JournalRatingSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}

	var total int64
	var _models []*models.JournalRatingModel

	query := this.db().Model(&models.JournalRatingModel{}).Preload("User")
	query = query.Where("journal_id = ?", journalID)

	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	query = query.Order(orderExpr)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_models).Error; err != nil {
		return nil, err
	}

	ratings := mapper.JournalRatingModelListToEntityList(_models)

	return entity2.NewPagingEntity(page, pageSize, total, ratings), nil
}

func (this *JournalRatingRepositoryImpl) GetByID(id uint) (*entity2.JournalRatingEntity, error) {
	var _model models.JournalRatingModel
	if err := this.db().Where("id = ?", id).First(&_model).Error; err != nil {
		return nil, err
	}
	return mapper.JournalRatingModelToEntity(&_model), nil
}

func (this *JournalRatingRepositoryImpl) GetStatsByJournalID(journalID uint) (*entity2.JournalRatingStatsEntity, error) {
	var stats entity2.JournalRatingStatsEntity

	var rows []struct {
		Stars uint
		Count uint
	}
	if err := this.db().Model(&models.JournalRatingModel{}).Select("stars, COUNT(*) as count").Where("journal_id = ?", journalID).Group("stars").Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, r := range rows {
		stats.Total += r.Count
		switch r.Stars {
		case 1:
			stats.Rating1 = r.Count
		case 2:
			stats.Rating2 = r.Count
		case 3:
			stats.Rating3 = r.Count
		case 4:
			stats.Rating4 = r.Count
		case 5:
			stats.Rating5 = r.Count
		}
	}

	return &stats, nil
}

func (this *JournalRatingRepositoryImpl) DeleteByIDAndUserID(id, userID uint) error {
	var rating models.JournalRatingModel
	if err := this.db().Where("id = ? AND user_id = ?", id, userID).First(&rating).Error; err != nil {
		return err
	}

	return this.db().Transaction(func(tx *gorm.DB) error {
		// Delete the rating
		if err := tx.Delete(&rating).Error; err != nil {
			return err
		}

		// Update journal rating
		return tx.
			Model(&models.JournalModel{}).
			Where("id = ?", rating.JournalID).
			UpdateColumns(map[string]interface{}{
				"rating_count": gorm.Expr("rating_count - ?", 1),
				"rating_sum":   gorm.Expr("rating_sum - ?", rating.Stars),
			}).Error
	})
}
