package userusecases

import (
	"errors"
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserArticleStatisticsByYearRangeUseCase struct {
	userRepository repository.UserRepository
}

// @inject
func NewUserArticleStatisticsByYearRangeUseCase(userRepository repository.UserRepository) *UserArticleStatisticsByYearRangeUseCase {
	return &UserArticleStatisticsByYearRangeUseCase{userRepository: userRepository}
}

func (this *UserArticleStatisticsByYearRangeUseCase) Execute(userID uint, fromYear, toYear int) (*entity.UserArticleStatisticsByYearRangeEntity, error) {
	if fromYear == 0 && toYear == 0 {
		activityStartYear, err := this.userRepository.GetUserActivityStartYear(userID)
		if err != nil {
			return nil, err
		}

		if activityStartYear == 0 {
			currentYear := time.Now().Year()
			fromYear = currentYear
			toYear = currentYear
		} else {
			fromYear = activityStartYear
			toYear = time.Now().Year()
		}
	}

	if fromYear > toYear {
		return nil, errors.New("'from' year cannot be greater than 'to' year")
	}

	entity, err := this.userRepository.GetArticleStatisticsByYearRange(userID, fromYear, toYear)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
