package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"

	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ProjectRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewProjectRepository(repository *BaseRepository) repository.ProjectRepository {
	return &ProjectRepositoryImpl{BaseRepository: repository}
}

func (this ProjectRepositoryImpl) GetAll(page, pageSize int) (*entity2.PagingEntity[entity2.ProjectEntity], error) {
	var total int64
	var _model []*models.ProjectModel

	countQuery := this.db().Model(&models.ProjectModel{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("error counting projects: %w", err)
	}

	query := this.db().Model(&models.ProjectModel{})
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_model).Error; err != nil {
		return nil, fmt.Errorf("error fetching projects: %w", err)
	}

	projects := mapper.ProjectModelListToEntityList(_model)

	return entity2.NewPagingEntity(page, pageSize, total, projects), nil
}

func (this *ProjectRepositoryImpl) GetByID(id uint) (*entity2.ProjectEntity, error) {
	var project models.ProjectModel
	if err := this.db().First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, infraError.Wrap(err)
	}
	return mapper.ProjectModelToEntity(&project), nil
}
