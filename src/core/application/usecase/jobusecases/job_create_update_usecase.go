package jobusecases

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"

	"slib.uz/src/core/domain/ports/repository"
)

type JobCreateUpdateUseCase struct {
	jobRepository          repository.JobRepository
	organizationRepository repository.OrganizationRepository
}

// @inject
func NewJobCreateUpdateUseCase(jobRepository repository.JobRepository, organizationRepository repository.OrganizationRepository) *JobCreateUpdateUseCase {
	return &JobCreateUpdateUseCase{
		jobRepository:          jobRepository,
		organizationRepository: organizationRepository,
	}
}

func (this *JobCreateUpdateUseCase) Execute(userID uint, createUpdateDTO *entity2.JobCreateUpdateEntity) (*entity2.JobEntity, error) {
	existingJob, err := this.jobRepository.GetByUserIDSingle(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return this.createJob(userID, createUpdateDTO)
		}
		return nil, err
	}

	return this.updateJob(userID, existingJob, createUpdateDTO)
}

func (this *JobCreateUpdateUseCase) createJob(userID uint, createUpdateDTO *entity2.JobCreateUpdateEntity) (*entity2.JobEntity, error) {
	organizationTin := ""
	organizationName := ""

	if createUpdateDTO.OrganizationID != nil {
		organization, err := this.organizationRepository.GetByID(*createUpdateDTO.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("organization not found: %w", err)
		}
		if organization.Tin != nil {
			organizationTin = *organization.Tin
		}
		organizationName = organization.Name
	}

	jobEntity := entity2.NewJobEntity(
		createUpdateDTO.OrganizationID,
		userID,
		organizationTin,
		organizationName,
		createUpdateDTO.PositionName,
	)

	createdJob, err := this.jobRepository.Create(jobEntity)
	if err != nil {
		return nil, err
	}

	return createdJob, nil
}

func (this *JobCreateUpdateUseCase) updateJob(userID uint, existingJob *entity2.JobEntity, createUpdateDTO *entity2.JobCreateUpdateEntity) (*entity2.JobEntity, error) {
	organizationTin := existingJob.OrganizationTin
	organizationName := existingJob.OrganizationName

	if createUpdateDTO.OrganizationID != nil {
		if existingJob.OrganizationID == nil || *existingJob.OrganizationID != *createUpdateDTO.OrganizationID {
			organization, err := this.organizationRepository.GetByID(*createUpdateDTO.OrganizationID)
			if err != nil {
				return nil, fmt.Errorf("organization not found: %w", err)
			}
			if organization.Tin != nil {
				organizationTin = *organization.Tin
			}
			organizationName = organization.Name
		}
	}

	updateEntity := entity2.NewJobEntity(
		createUpdateDTO.OrganizationID,
		userID,
		organizationTin,
		organizationName,
		createUpdateDTO.PositionName,
	)

	updatedJob, err := this.jobRepository.Update(userID, updateEntity)
	if err != nil {
		return nil, err
	}

	return updatedJob, nil
}
