package jobusecases

import (
	"errors"

	"gorm.io/gorm"

	"slib.uz/src/core/domain/ports/repository"
)

type JobDeleteUseCase struct {
	jobRepository repository.JobRepository
}

// @inject
func NewJobDeleteUseCase(jobRepository repository.JobRepository) *JobDeleteUseCase {
	return &JobDeleteUseCase{
		jobRepository: jobRepository,
	}
}

func (this *JobDeleteUseCase) Execute(userID uint) error {
	err := this.jobRepository.DeleteLastByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("job not found")
		}
		return err
	}

	return nil
}
