package userusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserArticleStatisticsByYearUseCase struct {
	userRepository repository.UserRepository
}

// @inject
func NewUserArticleStatisticsByYearUseCase(userRepository repository.UserRepository) *UserArticleStatisticsByYearUseCase {
	return &UserArticleStatisticsByYearUseCase{userRepository: userRepository}
}

func (this *UserArticleStatisticsByYearUseCase) Execute(userID uint, year int) (*entity.UserArticleStatisticsByYearEntity, error) {
	entity, err := this.userRepository.GetArticleStatisticsByYear(userID, year)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
