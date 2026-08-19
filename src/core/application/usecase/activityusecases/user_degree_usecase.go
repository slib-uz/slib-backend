package activityusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserDegreeUsecase struct {
	repository repository.DegreeRepository
}

// @inject
func NewUserDegreeUsecase(repository repository.DegreeRepository) *UserDegreeUsecase {
	return &UserDegreeUsecase{repository: repository}
}

func (this *UserDegreeUsecase) Execute(userId uint) (*entity.DegreeEntity, error) {
	degree, err := this.repository.GetByUserID(userId)

	if err != nil {
		return nil, err
	}

	return degree, nil
}
