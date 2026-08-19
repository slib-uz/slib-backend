package repository

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ReviewerRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewReviewerRepository(repository *BaseRepository) repository.ReviewerRepository {
	return &ReviewerRepositoryImpl{BaseRepository: repository}
}

func (this *ReviewerRepositoryImpl) Create(reviewer *entity.ReviewerEntity) (uint, error) {
	var _model = mapper.ReviewerEntityToModel(reviewer)
	if err := this.db().Create(_model).Error; err != nil {
		return 0, infraError.Wrap(err)
	}
	return _model.ID, nil
}

func (this *ReviewerRepositoryImpl) GetByScienceID(scienceID string) (*entity.ReviewerEntity, error) {
	var _model models.ReviewerModel

	if err := this.db().
		Where("science_id = ?", scienceID).
		First(&_model).
		Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.ReviewerModelToEntity(&_model), nil
}

func (this *ReviewerRepositoryImpl) GetIDByExternalID(externalID uint) (uint, error) {
	var _model models.ReviewerModel

	if err := this.db().
		Where("external_id = ?", externalID).
		First(&_model).
		Error; err != nil {
		return 0, infraError.Wrap(err)
	}
	return _model.ID, nil
}

func (this *ReviewerRepositoryImpl) GetByID(id uint) (*entity.ReviewerEntity, error) {
	var _model models.ReviewerModel

	if err := this.db().
		Where("id = ?", id).
		First(&_model).
		Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.ReviewerModelToEntity(&_model), nil
}

func (this *ReviewerRepositoryImpl) GetByJournalID(journalID uint) ([]*entity.ReviewerEntity, error) {
	var journal models.JournalModel

	if err := this.db().Select("id").Preload("Reviewers").Where("id = ?", journalID).First(&journal).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	reviewers := make([]*entity.ReviewerEntity, len(journal.Reviewers))
	for i, reviewer := range journal.Reviewers {
		reviewers[i] = mapper.ReviewerModelToEntity(reviewer)
	}
	return reviewers, nil
}

func (this *ReviewerRepositoryImpl) RemoveReviewerFromJournal(journalID, reviewerID uint) error {
	var journal models.JournalModel

	if err := this.db().Select("id").Preload("Reviewers").Where("id = ?", journalID).First(&journal).Error; err != nil {
		return infraError.Wrap(err)
	}

	return this.db().Model(&journal).Association("Reviewers").Delete(&models.ReviewerModel{Model: gorm.Model{ID: reviewerID}})
}

func (this *ReviewerRepositoryImpl) ListByIDs(ids []uint) ([]*entity.ReviewerEntity, error) {
	var _models []*models.ReviewerModel

	if err := this.db().Where("id IN ?", ids).Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	reviewers := make([]*entity.ReviewerEntity, len(_models))
	for i, r := range _models {
		reviewers[i] = mapper.ReviewerModelToEntity(r)
	}
	return reviewers, nil
}
