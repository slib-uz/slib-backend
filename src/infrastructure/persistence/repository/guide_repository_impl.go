package repository

import (
	"fmt"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type GuideRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewGuideRepository(repository *BaseRepository) repository.GuideRepository {
	return &GuideRepositoryImpl{BaseRepository: repository}
}

func (this *GuideRepositoryImpl) GetAll(page, pageSize int) (*entity2.PagingEntity[entity2.GuideListEntity], error) {
	var total int64
	var guideModels []*models.GuideModel

	countQuery := this.db().Model(&models.GuideModel{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("error counting guides: %w", err)
	}

	query := this.db().Model(&models.GuideModel{})
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&guideModels).Error; err != nil {
		return nil, fmt.Errorf("error fetching guide: %w", err)
	}

	guides := mapper.GuideModelListToEntityList(guideModels)

	return entity2.NewPagingEntity(page, pageSize, total, guides), nil
}

func (this *GuideRepositoryImpl) GetByID(id uint) (*entity2.GuideRetrieveEntity, error) {
	var guideModel models.GuideModel

	query := this.db().Model(&models.GuideModel{})
	if err := query.Where("id = ?", id).Find(&guideModel).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	guide := mapper.GuideRetrieveModelToEntity(&guideModel)

	return entity2.NewGuideRetrieveEntity(guide.ID, guide.Title, guide.Description, guide.FilePath, guide.VideoUrl), nil
}
