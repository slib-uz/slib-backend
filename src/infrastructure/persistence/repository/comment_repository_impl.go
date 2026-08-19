package repository

import (
	"errors"

	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type CommentRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewCommentRepository(baseRepository *BaseRepository) repository.CommentRepository {
	return &CommentRepositoryImpl{BaseRepository: baseRepository}
}

func (this *CommentRepositoryImpl) Create(comment *entity2.CommentEntity) error {

	var exists *models.CommentModel
	var article *models.ArticleModel
	if err := this.db().Model(&models.ArticleModel{}).Where("id = ?", comment.ArticleID).First(&article).Error; err != nil {
		return infraError.Wrap(err)
	}

	// Check if user already has a non-deleted comment for this article
	if err := this.db().Model(&models.CommentModel{}).
		Where("article_id = ? AND user_id = ? AND deleted_at IS NULL", comment.ArticleID, comment.UserID).
		First(&exists).Error; err == nil && exists.ID != 0 {
		return infraError.Wrap(errors.New("comment already exists"))
	}

	return this.db().Transaction(func(tx *gorm.DB) error {
		var _model = mapper.CommentEntityToModel(comment)
		if err := tx.Create(_model).Error; err != nil {
			return err
		}

		return tx.
			Model(&models.ArticleModel{}).
			Where("id = ?", comment.ArticleID).
			UpdateColumns(map[string]any{
				"rating_sum":   gorm.Expr("rating_sum + ?", comment.Rating),
				"rating_count": gorm.Expr("rating_count + 1"),
			}).Error
	})
}

func (this *CommentRepositoryImpl) GetByArticleID(articleID uint, page, pageSize int) (*entity2.PagingEntity[entity2.CommentEntity], error) {

	var comments []*models.CommentModel
	query := this.db().Model(&models.CommentModel{}).Order("created_at DESC").Where("article_id = ?", articleID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := query.
		Preload("User").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).
		Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.PagingMapper(comments, mapper.CommentModelToEntity, page, pageSize, total), nil

}

func (this *CommentRepositoryImpl) GetByUserID(userID uint, page, pageSize int) (*entity2.PagingEntity[entity2.CommentEntity], error) {

	var comments []*models.CommentModel
	query := this.db().Model(&models.CommentModel{}).Order("created_at DESC").Where("user_id = ? = ?", userID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.PagingMapper(comments, mapper.CommentModelToEntity, page, pageSize, total), nil

}

func (this *CommentRepositoryImpl) GetStatsByArticleID(articleID uint) (*entity2.CommentStatsEntity, error) {
	var stats entity2.CommentStatsEntity
	type ratingCount struct {
		Rating int
		Count  uint
	}

	var rows []ratingCount
	if err := this.db().
		Model(&models.CommentModel{}).
		Select("rating, COUNT(*) as count").
		Where("article_id = ?", articleID).
		Group("rating").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, r := range rows {
		stats.Total += r.Count
		switch r.Rating {
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

func (this *CommentRepositoryImpl) DeleteByIDAndUserID(id, userID uint) error {
	var comment models.CommentModel
	if err := this.db().Where("id = ? AND user_id = ?", id, userID).First(&comment).Error; err != nil {
		return err
	}

	return this.db().Transaction(func(tx *gorm.DB) error {
		// Delete the comment
		if err := tx.Delete(&comment).Error; err != nil {
			return err
		}

		// Update article rating
		return tx.
			Model(&models.ArticleModel{}).
			Where("id = ?", comment.ArticleID).
			UpdateColumns(map[string]any{
				"rating_sum":   gorm.Expr("rating_sum - ?", comment.Rating),
				"rating_count": gorm.Expr("rating_count - 1"),
			}).Error
	})
}
