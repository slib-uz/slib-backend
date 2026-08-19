package repository

import (
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
	"slib.uz/src/infrastructure/persistence/sorting"
)

type NewsRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewNewsRepository(baseRepository *BaseRepository) repository.NewsRepository {
	return &NewsRepositoryImpl{BaseRepository: baseRepository}
}

func (this *NewsRepositoryImpl) GetByID(newsID uint) (*entity2.NewsEntity, error) {
	var _model models.NewsModel

	if err := this.db().Where("id = ?", newsID).First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.NewsModelToEntity(&_model), nil
}

// NewsSortFields — yangiliklar ro'yxati uchun ruxsat etilgan tartiblash
// maydonlari. NewsModel.TableName() "news" qaytaradi.
var NewsSortFields = sorting.New("news.created_at DESC", map[string]string{
	"created_at": "news.created_at",
})

func (this *NewsRepositoryImpl) GetByPaging(page, pageSize int, ordering string, categoryID uint) (*entity2.PagingEntity[entity2.NewsEntity], error) {
	orderExpr, err := NewsSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}

	var total int64
	var _models []*models.NewsModel

	countQuery := this.db().Model(&models.NewsModel{})
	query := this.db().Model(&models.NewsModel{})

	if categoryID > 0 {
		countQuery = countQuery.Where("category_id = ?", categoryID)
		query = query.Where("category_id = ?", categoryID)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	query = query.Order(orderExpr)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_models).Error; err != nil {
		return nil, fmt.Errorf("error fetching news: %w", err)
	}

	news := mapper.NewsListModelToEntity(_models)

	return entity2.NewPagingEntity(page, pageSize, total, news), nil
}

func (this *NewsRepositoryImpl) UpdateViewsCount(counts map[uint]int64) error {

	if len(counts) == 0 {
		return nil
	}

	ids := make([]uint, 0, len(counts))
	params := make([]any, 0, len(counts)*2)
	caseSQL := "CASE id"

	for id, d := range counts {
		ids = append(ids, id)

		caseSQL += " WHEN ? THEN views_count + ?::bigint"
		params = append(params, id, d)
	}
	caseSQL += " ELSE views_count END"

	return this.db().
		Model(&models.NewsModel{}).
		Where("id IN ?", ids).
		Update("views_count", gorm.Expr(caseSQL, params...)).
		Error
}

func (this *NewsRepositoryImpl) Create(news *entity2.NewsEntity) error {
	model := models.NewNewsModel(
		news.ID,
		datatypes.JSON(mapper.ToGormJson(news.Title)),
		datatypes.JSON(mapper.ToGormJson(news.Body)),
		news.CategoryID,
		news.Image,
		news.ViewsCount,
	)

	if err := this.db().Create(model).Error; err != nil {
		return infraError.Wrap(err)
	}

	news.ID = model.ID
	news.CreatedAt = model.CreatedAt

	return nil
}

func (this *NewsRepositoryImpl) Update(id uint, news *entity2.NewsEntity) error {
	updates := map[string]interface{}{
		"title":       datatypes.JSON(mapper.ToGormJson(news.Title)),
		"body":        datatypes.JSON(mapper.ToGormJson(news.Body)),
		"category_id": news.CategoryID,
		"image":       news.Image,
	}

	if err := this.db().Model(&models.NewsModel{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return infraError.Wrap(err)
	}

	return nil
}

func (this *NewsRepositoryImpl) Delete(id uint) error {
	if err := this.db().Delete(&models.NewsModel{}, id).Error; err != nil {
		return infraError.Wrap(err)
	}

	return nil
}
