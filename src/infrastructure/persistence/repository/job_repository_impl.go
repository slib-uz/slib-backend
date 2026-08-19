package repository

import (
	"errors"
	"slib.uz/src/infrastructure/db"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type JobRepositoryImpl struct {
	database *db.Database
}

// @inject
func NewJobRepository(database *db.Database) repository.JobRepository {
	return &JobRepositoryImpl{
		database: database,
	}
}

func (this *JobRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *JobRepositoryImpl) GetByUserID(userID uint) ([]*entity.JobEntity, error) {

	var jobs []*models.JobModel
	if err := this.db().Where("user_id = ?", userID).Limit(1).Find(&jobs).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	var jobEntities []*entity.JobEntity
	for _, job := range jobs {
		jobEntities = append(jobEntities, mapper.JobModelToEntity(job))
	}
	return jobEntities, nil
}

func (this *JobRepositoryImpl) GetByUserIDSingle(userID uint) (*entity.JobEntity, error) {
	var job models.JobModel
	if err := this.db().Where("user_id = ?", userID).First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, infraError.Wrap(err)
	}
	return mapper.JobModelToEntity(&job), nil
}

func (this *JobRepositoryImpl) Create(job *entity.JobEntity) (*entity.JobEntity, error) {
	jobModel := mapper.JobEntityToModel(job.UserID, job)
	if err := this.db().Create(jobModel).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.JobModelToEntity(jobModel), nil
}

func (this *JobRepositoryImpl) Update(userID uint, job *entity.JobEntity) (*entity.JobEntity, error) {
	var existingJob models.JobModel
	if err := this.db().Where("user_id = ?", userID).First(&existingJob).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	existingJob.OrganizationID = job.OrganizationID
	existingJob.OrganizationTin = job.OrganizationTin
	existingJob.OrganizationName = job.OrganizationName
	existingJob.PositionName = job.PositionName

	if err := this.db().Save(&existingJob).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.JobModelToEntity(&existingJob), nil
}

func (this *JobRepositoryImpl) UpdateOrCreate(userID uint, entities []*entity.JobEntity) error {
	for _, item := range entities {
		job := mapper.JobEntityToModel(userID, item)

		var existing *models.JobModel
		err := this.db().Where("user_id = ?", userID).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if createErr := this.db().Create(&job).Error; createErr != nil {
					return createErr
				}
			} else {
				return err
			}
		} else {
			existing = mapper.JobUpdateMapper(existing, job)
			if updateErr := this.db().Save(&existing).Error; updateErr != nil {
				return updateErr
			}
		}
	}
	return nil
}

func (this *JobRepositoryImpl) DeleteLastByUserID(userID uint) error {
	var job models.JobModel
	if err := this.db().Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at DESC").First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return infraError.Wrap(err)
	}

	if err := this.db().Delete(&job).Error; err != nil {
		return infraError.Wrap(err)
	}

	return nil
}
