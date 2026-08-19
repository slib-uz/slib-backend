package usecase

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ApplicationCheckForPaymentUsecase struct {
	repository repository.ApplicationRepository
}

// @inject
func NewApplicationCheckForPaymentUsecase(repository repository.ApplicationRepository) *ApplicationCheckForPaymentUsecase {
	return &ApplicationCheckForPaymentUsecase{repository: repository}
}

func (this *ApplicationCheckForPaymentUsecase) Execute(applicationID uint) (*entity.ReviewStageEntity, error) {
	return this.repository.GetCheckForPaymentByApplicationID(applicationID)
}
