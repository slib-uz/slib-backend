package repository

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type SpellCheckRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewSpellCheckRepository(baseRepository *BaseRepository) repository.SpellCheckResultRepository {
	return &SpellCheckRepositoryImpl{BaseRepository: baseRepository}
}

func (this *SpellCheckRepositoryImpl) Create(result *entity2.SpellCheckResultEntity) (*entity2.SpellCheckResultEntity, error) {
	_model := mapper.SpellCheckResultEntityToModel(result)

	if err := this.db().Create(&_model).Error; err != nil {
		return nil, err
	}
	return mapper.SpellCheckResultModelToEntity(_model), nil
}

func (this *SpellCheckRepositoryImpl) Update(id uint, status int, resultFile *string, resultTime *time.Time) error {

	return this.db().
		Model(&models.SpellCheckResultModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      status,
			"result_file": resultFile,
			"result_time": resultTime,
		}).Error
}

func (this *SpellCheckRepositoryImpl) GetByApplicationID(applicationID uint) ([]*entity2.SpellCheckResultEntity, error) {
	var _models []models.SpellCheckResultModel

	if err := this.db().
		Preload("Submitter").
		Where("application_id = ?", applicationID).
		Find(&_models).Error; err != nil {
		return nil, err
	}
	var results = make([]*entity2.SpellCheckResultEntity, len(_models))
	for i, model := range _models {
		results[i] = mapper.SpellCheckResultModelToEntity(&model)
	}
	return results, nil
}

func (this *SpellCheckRepositoryImpl) GetLatestByApplicationID(applicationID uint) (*entity2.SpellCheckResultEntity, error) {
	var _model models.SpellCheckResultModel

	if err := this.db().
		Preload("Submitter").
		Where("application_id = ?", applicationID).
		Order("created_at DESC").
		First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.SpellCheckResultModelToEntity(&_model), nil
}

func (this *SpellCheckRepositoryImpl) StatsByJournal(publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalSpellcheckStatsEntity], error) {

	var stats []*entity2.JournalSpellcheckStatsEntity

	query := this.db().Table("spellcheck_results").
		Select("spellcheck_results.journal_id as journal_id, journals.name as journal_name, COUNT(CASE WHEN spellcheck_results.status = 1 THEN 1 END) AS success, COUNT(CASE WHEN spellcheck_results.status = -1 THEN 1 END) AS failed").
		Joins("JOIN journals ON spellcheck_results.journal_id = journals.id").
		Where("spellcheck_results.created_at BETWEEN ? AND ?", startDate, endDate)

	if publisherID > 0 {
		query = query.Where("journals.publisher_id = ?", publisherID)
	}
	query = query.Group("spellcheck_results.journal_id, journals.name")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Scan(&stats).Error; err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(page, pageSize, total, stats), nil
}

func (this *SpellCheckRepositoryImpl) GetByReviewStageID(reviewStageID uint) ([]*entity2.SpellCheckResultEntity, error) {
	var _models []models.SpellCheckResultModel

	if err := this.db().
		Preload("Submitter").
		Where("review_stage_id = ?", reviewStageID).
		Find(&_models).Error; err != nil {
		return nil, err
	}
	var results = make([]*entity2.SpellCheckResultEntity, len(_models))
	for i, model := range _models {
		results[i] = mapper.SpellCheckResultModelToEntity(&model)
	}
	return results, nil
}

func (this *SpellCheckRepositoryImpl) GetByIDWithApplication(id uint) (*entity2.SpellCheckResultEntity, error) {
	var _model models.SpellCheckResultModel
	if err := this.db().Preload("Application").First(&_model, id).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.SpellCheckResultModelToEntity(&_model), nil
}

func (this *SpellCheckRepositoryImpl) StatsByPublisher() ([]*entity2.PublisherSpellcheckStatsEntity, error) {
	var stats []*entity2.PublisherSpellcheckStatsEntity

	if err := this.db().Table("spellcheck_results").
		Select("journals.publisher_id as publisher_id, publishers.name AS publisher_name, COUNT(CASE WHEN spellcheck_results.status = 1 THEN 1 END) AS success, COUNT(CASE WHEN spellcheck_results.status = -1 THEN 1 END) AS failed").
		Joins("JOIN journals ON spellcheck_results.journal_id = journals.id").
		Joins("JOIN publishers ON journals.publisher_id = publishers.id").
		Group("journals.publisher_id, publishers.name").
		Scan(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (this *SpellCheckRepositoryImpl) GetByJournalID(journalID uint, page, pageSize int) (*entity2.PagingEntity[entity2.SpellCheckResultEntity], error) {
	var _models []models.SpellCheckResultModel

	query := this.db().
		Preload("Submitter").
		Preload("Application").
		Preload("Application.Article").
		Preload("Journal").
		Where("journal_id = ?", journalID)

	var total int64
	if err := query.Model(&models.SpellCheckResultModel{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&_models).Error; err != nil {
		return nil, err
	}

	results := make([]*entity2.SpellCheckResultEntity, len(_models))
	for i, model := range _models {
		results[i] = mapper.SpellCheckResultModelToEntity(&model)
	}

	return entity2.NewPagingEntity(page, pageSize, total, results), nil
}
