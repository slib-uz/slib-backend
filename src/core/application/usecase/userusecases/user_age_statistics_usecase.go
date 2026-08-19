package userusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserAgeStatisticsUseCase struct {
	userRepository repository.UserRepository
}

// @inject
func NewUserAgeStatisticsUseCase(userRepository repository.UserRepository) *UserAgeStatisticsUseCase {
	return &UserAgeStatisticsUseCase{
		userRepository: userRepository,
	}
}

func (this *UserAgeStatisticsUseCase) Execute() (*entity.UserAgeStatisticsEntity, error) {
	stats, err := this.userRepository.GetAgeStatistics()
	if err != nil {
		return nil, err
	}

	return stats, nil
}
