package activityusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserDoctorateUseCase struct {
	repository repository.DoctorateRepository
}

// @inject
func NewUserDoctorateUseCase(repository repository.DoctorateRepository) *UserDoctorateUseCase {
	return &UserDoctorateUseCase{repository: repository}
}

func (this *UserDoctorateUseCase) Execute(userId uint) (*entity.DoctorateEntity, error) {
	_doctorate, err := this.repository.GetByUserID(userId)

	if err != nil {
		return nil, err
	}

	return _doctorate, nil
}
