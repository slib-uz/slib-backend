package newsusecases

import (
	"context"

	"slib.uz/src/core/domain/ports/repository"
)

type NewsDeleteUseCase struct {
	newsRepository repository.NewsRepository
}

// @inject
func NewNewsDeleteUseCase(newsRepository repository.NewsRepository) *NewsDeleteUseCase {
	return &NewsDeleteUseCase{newsRepository: newsRepository}
}

func (this *NewsDeleteUseCase) Execute(ctx context.Context, id uint) error {
	if err := this.newsRepository.Delete(id); err != nil {
		return err
	}

	return nil
}
