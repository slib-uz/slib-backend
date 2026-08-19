package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type AiDetectRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewAiDetectRepository(repository *BaseRepository) repository.AiDetectRepository {
	return &AiDetectRepositoryImpl{BaseRepository: repository}
}

func (this *AiDetectRepositoryImpl) Create(data *entity2.AiDetectResultEntity) error {
	return this.db().Create(mapper.AiDetectResultEntityToModel(data)).Error
}

func (this *AiDetectRepositoryImpl) GetByExternalID(externalID uint) (*entity2.AiDetectResultEntity, error) {
	var _model models.AiDetectResultModel
	if err := this.db().Where("external_id = ?", externalID).First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.AiDetectResultModelToEntity(&_model), nil
}

func (this *AiDetectRepositoryImpl) GetByID(id uint) (*entity2.AiDetectResultEntity, error) {
	var _model models.AiDetectResultModel
	if err := this.db().First(&_model, id).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.AiDetectResultModelToEntity(&_model), nil
}

func (this *AiDetectRepositoryImpl) GetByIdWithApplication(id uint) (*entity2.AiDetectResultEntity, error) {
	var _model models.AiDetectResultModel
	if err := this.db().
		Preload("Application").
		Preload("Application.User").
		First(&_model, id).
		Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.AiDetectResultModelToEntity(&_model), nil
}

func (this *AiDetectRepositoryImpl) Update(entity *entity2.AiDetectResultEntity) error {
	var _model = *mapper.AiDetectResultEntityToModel(entity)
	return this.db().Model(&_model).Updates(_model).Error
}

func (this *AiDetectRepositoryImpl) FindIDByExternalID(externalID uint) (uint, error) {
	var _model models.AiDetectResultModel
	if err := this.db().Select("id").Where("external_id = ?", externalID).First(&_model).Error; err != nil {
		return 0, infraError.Wrap(err)
	}
	return _model.ID, nil
}

func (this *AiDetectRepositoryImpl) GetAllByReviewStageID(reviewStageID uint) ([]*entity2.AiDetectResultEntity, error) {
	var _models []*models.AiDetectResultModel
	if err := this.db().Where("review_stage_id = ?", reviewStageID).Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	result := make([]*entity2.AiDetectResultEntity, len(_models))
	for i, item := range _models {
		result[i] = mapper.AiDetectResultModelToEntity(item)
	}
	return result, nil
}

func (this *AiDetectRepositoryImpl) GetByJournalID(journalID uint, page, pageSize int) (*entity2.PagingEntity[entity2.AiDetectResultEntity], error) {
	var _models []models.AiDetectResultModel

	query := this.db().
		Preload("Application").
		Preload("Journal").
		Preload("Article").
		Preload("ReviewStage").
		Where("journal_id = ?", journalID)

	var total int64
	if err := query.Model(&models.AiDetectResultModel{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&_models).Error; err != nil {
		return nil, err
	}

	results := make([]*entity2.AiDetectResultEntity, len(_models))
	for i, model := range _models {
		results[i] = mapper.AiDetectResultModelToEntity(&model)
	}

	return entity2.NewPagingEntity(page, pageSize, total, results), nil
}
