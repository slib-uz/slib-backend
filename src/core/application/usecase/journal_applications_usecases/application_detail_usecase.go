package journal_applications_usecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ApplicationDetailUsecase struct {
	repository repository.JournalApplicationRepository
}

// @inject
func NewApplicationDetailUsecase(repository repository.JournalApplicationRepository) *ApplicationDetailUsecase {
	return &ApplicationDetailUsecase{repository: repository}
}

func (this *ApplicationDetailUsecase) Execute(id uint) (*entity.JournalApplicationEntity, error) {
	app, err := this.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return app, nil
}
