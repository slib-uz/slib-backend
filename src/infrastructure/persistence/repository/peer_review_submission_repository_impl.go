package repository

import (
	"time"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type PeerReviewSubmissionRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewPeerReviewSubmissionRepository(baseRepository *BaseRepository) repository.PeerReviewSubmissionRepository {
	return &PeerReviewSubmissionRepositoryImpl{BaseRepository: baseRepository}
}

func (this *PeerReviewSubmissionRepositoryImpl) BulkCreate(submissions []*entity.PeerReviewSubmissionEntity) error {
	var _models = make([]*models.PeerReviewSubmissionModel, len(submissions))

	for i, s := range submissions {
		_models[i] = mapper.PeerReviewSubmissionEntityToModel(s)
	}

	return this.db().Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(_models, 10).Error
	})
}

func (this *PeerReviewSubmissionRepositoryImpl) IsExistsByExternalIdempotencyKey(idempotencyKey string) (bool, error) {
	var count int64

	err := this.db().
		Model(&models.PeerReviewSubmissionModel{}).
		Where("external_idempotency_key = ?", idempotencyKey).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (this *PeerReviewSubmissionRepositoryImpl) UpdateStatusAndDeadline(externalID uint, status enum.PeerReviewStatus, deadline time.Time, oldDeadline *time.Time, version int64) error {

	updatedFields := map[string]any{
		"status":     status,
		"version":    version,
		"updated_at": time.Now(),
		"deadline":   deadline,
	}

	if oldDeadline != nil {
		updatedFields["old_deadline"] = oldDeadline
	}

	return this.db().
		Model(&models.PeerReviewSubmissionModel{}).
		Where("external_id = ?", externalID).
		Updates(updatedFields).
		Error
}

func (this *PeerReviewSubmissionRepositoryImpl) GetByExternalID(externalID uint) (*entity.PeerReviewSubmissionEntity, error) {
	var _model models.PeerReviewSubmissionModel
	if err := this.db().
		Where("external_id = ?", externalID).
		First(&_model).Error; err != nil {
		return nil, err
	}
	return mapper.PeerReviewSubmissionModelToEntity(&_model), nil
}

func (this *PeerReviewSubmissionRepositoryImpl) GetByApplicationID(applicationID uint) ([]*entity.PeerReviewSubmissionEntity, error) {

	var _models []*models.PeerReviewSubmissionModel

	if err := this.db().
		Preload("Reviewer").
		Where("application_id = ?", applicationID).
		Find(&_models).Error; err != nil {
		return nil, err
	}

	var submissions = make([]*entity.PeerReviewSubmissionEntity, len(_models))

	for i, _model := range _models {
		submissions[i] = mapper.PeerReviewSubmissionModelToEntity(_model)
	}

	return submissions, nil

}

func (this *PeerReviewSubmissionRepositoryImpl) GetByIDWithApplication(id uint) (*entity.PeerReviewSubmissionEntity, error) {
	var _model models.PeerReviewSubmissionModel

	if err := this.db().
		Preload("Application").
		Where("id = ?", id).
		First(&_model).Error; err != nil {
		return nil, err
	}

	return mapper.PeerReviewSubmissionModelToEntity(&_model), nil
}
