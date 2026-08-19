package activityusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserJobsUseCase struct {
	repository repository.JobRepository
}

// @inject
func NewUserJobsUseCase(repository repository.JobRepository) *UserJobsUseCase {
	return &UserJobsUseCase{repository: repository}
}

func (this *UserJobsUseCase) Execute(userId uint) (*entity.JobEntity, error) {
	_jobs, err := this.repository.GetByUserID(userId)

	if err != nil {
		return nil, err
	}

	if len(_jobs) == 0 {
		return nil, nil
	}

	return _jobs[0], nil
}
