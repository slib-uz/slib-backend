package userusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserGenderStatisticsUseCase struct {
	userRepository repository.UserRepository
}

// @inject
func NewUserGenderStatisticsUseCase(userRepository repository.UserRepository) *UserGenderStatisticsUseCase {
	return &UserGenderStatisticsUseCase{
		userRepository: userRepository,
	}
}

func (this *UserGenderStatisticsUseCase) Execute() (*entity.UserGenderStatisticsEntity, error) {
	stats, err := this.userRepository.GetGenderStatistics()
	if err != nil {
		return nil, err
	}

	return stats, nil
}
