package repository

import (
	"fmt"
	"time"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type AntiPlagRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewAntiPlagRepository(repository *BaseRepository) repository.AntiPlagRepository {
	return &AntiPlagRepositoryImpl{BaseRepository: repository}
}

func (this *AntiPlagRepositoryImpl) Create(data *entity2.AntiPlagResultEntity) error {
	return this.db().Create(mapper.AntiPlagResultEntityToModel(data)).Error
}

func (this *AntiPlagRepositoryImpl) GetByExternalID(externalID uint) (*entity2.AntiPlagResultEntity, error) {
	var _model models.AntiPlagResultModel

	if err := this.db().Where("external_id = ?", externalID).First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.AntiPlagResultModelToEntity(&_model), nil
}

func (this *AntiPlagRepositoryImpl) GetByID(id uint) (*entity2.AntiPlagResultEntity, error) {
	var _model models.AntiPlagResultModel

	if err := this.db().First(&_model, id).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.AntiPlagResultModelToEntity(&_model), nil
}

func (this *AntiPlagRepositoryImpl) Update(entity *entity2.AntiPlagResultEntity) error {
	var _model = *mapper.AntiPlagResultEntityToModel(entity)
	return this.db().Model(&_model).Updates(_model).Error
}

func (this *AntiPlagRepositoryImpl) FindIDByExternalID(externalID uint) (uint, error) {
	var _model models.AntiPlagResultModel
	if err := this.db().Select("id").Where("external_id = ?", externalID).First(&_model).Error; err != nil {
		return 0, infraError.Wrap(err)
	}
	return _model.ID, nil
}

func (this *AntiPlagRepositoryImpl) GetAllByApplicationID(applicationID uint) ([]*entity2.AntiPlagResultEntity, error) {
	var _models []*models.AntiPlagResultModel

	if err := this.db().Where("application_id = ?", applicationID).Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.AntiPlagResulModelListToEntityList(_models), nil
}

func (this *AntiPlagRepositoryImpl) GetLatestByApplicationID(applicationID uint) (*entity2.AntiPlagResultEntity, error) {
	var _model models.AntiPlagResultModel

	if err := this.db().Where("application_id = ?", applicationID).Order("created_at DESC").First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.AntiPlagResultModelToEntity(&_model), nil
}

func (this *AntiPlagRepositoryImpl) StatsByJournal(publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalAntiPlagStatsEntity], error) {

	var stats []*entity2.JournalAntiPlagStatsEntity

	query := this.db().Table("anti_plag_result_models").
		Select("anti_plag_result_models.journal_id as journal_id, journals.name as journal_name, COUNT(CASE WHEN anti_plag_result_models.status = 7 THEN 1 END) AS success, COUNT(CASE WHEN anti_plag_result_models.status = 3 THEN 1 END) AS failed").
		Joins("JOIN journals ON anti_plag_result_models.journal_id = journals.id").
		Where("anti_plag_result_models.created_at BETWEEN ? AND ?", startDate, endDate)

	if publisherID > 0 {
		query = query.Where("journals.publisher_id = ?", publisherID)
	}

	query = query.Group("anti_plag_result_models.journal_id, journals.name")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("error fetching journal stats: %w", err)
	}

	return entity2.NewPagingEntity(page, pageSize, total, stats), nil
}

func (this *AntiPlagRepositoryImpl) StatsByPublisher() ([]*entity2.PublisherAntiplagStatsEntity, error) {
	var stats []*entity2.PublisherAntiplagStatsEntity

	err := this.db().Table("anti_plag_result_models").
		Select("journals.publisher_id as publisher_id, publishers.name AS publisher_name, COUNT(CASE WHEN anti_plag_result_models.status = 7 THEN 1 END) AS success, COUNT(CASE WHEN anti_plag_result_models.status = 3 THEN 1 END) AS failed").
		Joins("JOIN journals ON anti_plag_result_models.journal_id = journals.id").
		Joins("JOIN publishers ON journals.publisher_id = publishers.id").
		Group("journals.publisher_id, publishers.name").
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching publisher stats: %w", err)
	}
	return stats, nil
}

func (this *AntiPlagRepositoryImpl) GetAllByReviewStageID(reviewStageID uint) ([]*entity2.AntiPlagResultEntity, error) {
	var _models []*models.AntiPlagResultModel

	if err := this.db().Where("review_stage_id = ?", reviewStageID).Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.AntiPlagResulModelListToEntityList(_models), nil
}

func (this *AntiPlagRepositoryImpl) GetByIdWithApplication(id uint) (*entity2.AntiPlagResultEntity, error) {
	var _model models.AntiPlagResultModel

	if err := this.db().
		Preload("Application").
		Preload("Application.User").
		First(&_model, id).
		Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.AntiPlagResultModelToEntity(&_model), nil
}

func (this *AntiPlagRepositoryImpl) UpdateCertificate(externalID uint, certificate *string) error {

	return this.db().
		Model(&models.AntiPlagResultModel{}).
		Where("external_id = ?", externalID).
		Updates(map[string]any{"certificate": certificate}).
		Error

}

func (this *AntiPlagRepositoryImpl) GetByJournalID(journalID uint, page, pageSize int) (*entity2.PagingEntity[entity2.AntiPlagResultEntity], error) {
	var _models []models.AntiPlagResultModel

	query := this.db().
		Preload("Application").
		Preload("Journal").
		Preload("Article").
		Preload("ReviewStage").
		Where("journal_id = ?", journalID)

	var total int64
	if err := query.Model(&models.AntiPlagResultModel{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&_models).Error; err != nil {
		return nil, err
	}

	results := make([]*entity2.AntiPlagResultEntity, len(_models))
	for i, model := range _models {
		results[i] = mapper.AntiPlagResultModelToEntity(&model)
	}

	return entity2.NewPagingEntity(page, pageSize, total, results), nil
}
