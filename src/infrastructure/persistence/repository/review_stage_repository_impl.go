package repository

import (
	"time"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ReviewStageRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewReviewStageRepository(repository *BaseRepository) repository.ReviewStageRepository {
	return &ReviewStageRepositoryImpl{BaseRepository: repository}
}

func (this *ReviewStageRepositoryImpl) GetByAppID(appID uint) ([]*entity.ReviewStageEntity, error) {
	var reviewStages []*models.ReviewStageModel
	if err := this.db().Where("application_id = ?", appID).Find(&reviewStages).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ReviewStageModelsToEntities(reviewStages), nil
}

func (this *ReviewStageRepositoryImpl) Create(reviewStage *entity.ReviewStageEntity) (uint, error) {
	reviewStageModel := mapper.ReviewStateEntityToModel(reviewStage)
	if err := this.db().Create(reviewStageModel).Error; err != nil {
		return 0, _errors.Wrap(err)
	}
	return reviewStageModel.ID, nil
}

func (this *ReviewStageRepositoryImpl) GetBy(appId uint, stage enum.Stage, status enum.Status) (*entity.ReviewStageEntity, error) {
	var reviewStageModel models.ReviewStageModel
	if err := this.db().Preload("Application").Where("application_id = ? AND stage = ? AND status = ? AND is_old = ?", appId, stage, status, false).Order("created_at DESC").First(&reviewStageModel).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ReviewStageModelToEntity(&reviewStageModel), nil
}

func (this *ReviewStageRepositoryImpl) GetWithApplication(id uint) (*entity.ReviewStageEntity, error) {
	var reviewStageModel models.ReviewStageModel

	if err := this.db().
		Where("id = ?", id).
		Preload("Application").
		First(&reviewStageModel).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.ReviewStageModelToEntity(&reviewStageModel), nil
}

func (this *ReviewStageRepositoryImpl) GetByID(id uint) (*entity.ReviewStageEntity, error) {
	var reviewStageModel models.ReviewStageModel

	if err := this.db().
		Where("id = ?", id).
		Preload("Application").
		Preload("Reviewer").
		First(&reviewStageModel).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ReviewStageModelToEntity(&reviewStageModel), nil
}

func (this *ReviewStageRepositoryImpl) UpdateToNext(reviewStage *entity.ReviewStageEntity, nextStage *entity.ReviewStageEntity) error {
	_model := mapper.ReviewStateEntityToModel(reviewStage)

	return this.db().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ReviewStageModel{}).
			Where("id = ?", _model.ID).
			Updates(_model).
			Error; err != nil {
			return _errors.Wrap(err)
		}

		if nextStage != nil {
			_nextStageModel := mapper.ReviewStateEntityToModel(nextStage)
			//create next review stage with pending status
			return tx.Create(&_nextStageModel).Error
		}
		return nil
	})
}

func (this *ReviewStageRepositoryImpl) GetByStage(appId uint, stage enum.Stage) (*entity.ReviewStageEntity, error) {
	var _model models.ReviewStageModel
	if err := this.db().
		Where("application_id = ? AND stage = ? AND is_old = ?", appId, stage, false).
		Preload("Application").
		Preload("Reviewer").
		First(&_model).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ReviewStageModelToEntity(&_model), nil
}

func (this *ReviewStageRepositoryImpl) IsExists(appId uint, stage enum.Stage, status enum.Status) (bool, error) {
	var count int64
	if err := this.db().
		Model(&models.ReviewStageModel{}).
		Where("application_id = ? AND stage = ? AND status = ? AND is_old = ?", appId, stage, status, false).
		Count(&count).Error; err != nil {
		return false, _errors.Wrap(err)
	}
	return count > 0, nil
}

func (this *ReviewStageRepositoryImpl) GetByIDWithArticle(id uint) (*entity.ReviewStageEntity, error) {
	var reviewStageModel models.ReviewStageModel

	if err := this.db().
		Where("id = ?", id).
		Preload("Application").
		Preload("Application.Article").
		First(&reviewStageModel).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ReviewStageModelToEntity(&reviewStageModel), nil
}

func (this *ReviewStageRepositoryImpl) UpdateDeadline(id uint, deadline time.Time, deadlineType enum.DeadlineType) error {

	var _model = models.ReviewStageModel{}

	switch deadlineType {
	case enum.DeadlineTypeReviewDeadline:
		_model.Deadline = deadline
	case enum.DeadlineTypeResubmitDeadline:
		_model.ResubmitDeadline = &deadline
	}

	return this.db().Where("id = ?", id).Updates(_model).Error
}

func (this *ReviewStageRepositoryImpl) GetExpiringStages() ([]*entity.ReviewStageEntity, error) {
	var _models []*models.ReviewStageModel

	now := time.Now()
	threshold := now.Add(72 * time.Hour)

	if err := this.db().
		Joins("JOIN article_applications ON article_applications.id = review_stages.application_id AND article_applications.deleted_at IS NULL").
		Where("deadline BETWEEN ? AND ? AND status = ? AND is_old = ?", now, threshold, enum.StatusPending, false).
		Preload("Application").
		Find(&_models).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.ReviewStageModelsToEntities(_models), nil
}

func (this *ReviewStageRepositoryImpl) GetStatistics(journalID *uint) ([]*entity.ReviewStageStatisticsEntity, error) {
	var results []*entity.ReviewStageStatisticsEntity

	query := `
SELECT
    CASE rs.stage
        WHEN 10 THEN 'Technical Review'
        WHEN 20 THEN 'Peer Review'
        WHEN 30 THEN 'Payment'
        WHEN 40 THEN 'Publish'
    END AS stage_name,
    rs.stage AS stage_number,
    COUNT(*) FILTER (WHERE rs.status = -1) AS rejected,
    COUNT(*) FILTER (WHERE rs.status = 0)  AS pending,
    COUNT(*) FILTER (WHERE rs.status = 1)  AS accepted,
    COUNT(*) AS total
FROM review_stages rs
JOIN article_applications aa
    ON aa.id = rs.application_id
   AND aa.deleted_at IS NULL
WHERE rs.is_old = false
  AND (?::bigint IS NULL OR aa.journal_id = ?)
GROUP BY rs.stage
ORDER BY rs.stage`

	if err := this.db().Raw(query, journalID, journalID).Scan(&results).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return results, nil
}

func (this *ReviewStageRepositoryImpl) GetOverdueStatistics(journalID *uint) ([]*entity.ReviewStageOverdueStatisticsEntity, error) {
	var results []*entity.ReviewStageOverdueStatisticsEntity

	query := `
SELECT
    CASE rs.stage
        WHEN 10 THEN 'Technical Review'
        WHEN 20 THEN 'Peer Review'
        WHEN 40 THEN 'Publish'
    END AS stage_name,
    rs.stage AS stage_number,
    COUNT(*) AS overdue_count
FROM review_stages rs
JOIN article_applications aa
    ON aa.id = rs.application_id
   AND aa.deleted_at IS NULL
JOIN journals j
    ON j.id = aa.journal_id
   AND j.deleted_at IS NULL
WHERE rs.is_old = false
  AND rs.status = ?
  AND rs.deadline < NOW()
  AND rs.stage IN (?, ?, ?)
  AND (?::bigint IS NULL OR aa.journal_id = ?)
GROUP BY rs.stage
ORDER BY rs.stage`

	if err := this.db().Raw(query, enum.StatusPending, enum.TechnicalReviewStage, enum.PeerReviewStage, enum.PublishStage, journalID, journalID).Scan(&results).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return results, nil
}

func (this *ReviewStageRepositoryImpl) GetOverdueList(page, pageSize int, journalID *uint) (*entity.PagingEntity[entity.ReviewStageOverdueEntity], error) {
	countQuery := `
SELECT COUNT(*)
FROM review_stages rs
JOIN article_applications aa
    ON aa.id = rs.application_id
   AND aa.deleted_at IS NULL
JOIN journals j
    ON j.id = aa.journal_id
   AND j.deleted_at IS NULL
WHERE rs.is_old = false
  AND rs.status = ?
  AND rs.deadline < NOW()
  AND rs.stage IN (?, ?, ?)
  AND (?::bigint IS NULL OR aa.journal_id = ?)`

	var total int64
	if err := this.db().Raw(countQuery, enum.StatusPending, enum.TechnicalReviewStage, enum.PeerReviewStage, enum.PublishStage, journalID, journalID).Scan(&total).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	var rows []mapper.ReviewStageOverdueRow

	query := `
SELECT
    rs.id AS review_stage_id,
    rs.application_id,
    aa.number AS application_number,
    aa.article_id,
    a.name AS article_name,
    rs.stage AS stage_number,
    CASE rs.stage
        WHEN 10 THEN 'Technical Review'
        WHEN 20 THEN 'Peer Review'
        WHEN 40 THEN 'Publish'
    END AS stage_name,
    rs.deadline,
    rs.created_at,
    EXTRACT(DAY FROM NOW() - rs.deadline)::int AS overdue_days,
    aa.journal_id,
    j.name AS journal_name
FROM review_stages rs
JOIN article_applications aa
    ON aa.id = rs.application_id
   AND aa.deleted_at IS NULL
JOIN articles a
    ON a.id = aa.article_id
   AND a.deleted_at IS NULL
JOIN journals j
    ON j.id = aa.journal_id
   AND j.deleted_at IS NULL
WHERE rs.is_old = false
  AND rs.status = ?
  AND rs.deadline < NOW()
  AND rs.stage IN (?, ?, ?)
  AND (?::bigint IS NULL OR aa.journal_id = ?)
ORDER BY rs.deadline
LIMIT ? OFFSET ?`

	offset := (page - 1) * pageSize
	if err := this.db().Raw(query, enum.StatusPending, enum.TechnicalReviewStage, enum.PeerReviewStage, enum.PublishStage, journalID, journalID, pageSize, offset).Scan(&rows).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	results := mapper.ReviewStageOverdueRowsToEntities(rows)

	return entity.NewPagingEntity(page, pageSize, total, results), nil
}
