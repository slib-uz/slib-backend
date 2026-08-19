package commentusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type CommentStatsUsecase struct {
	repository repository.CommentRepository
}

// @inject
func NewCommentStatsUsecase(repository repository.CommentRepository) *CommentStatsUsecase {
	return &CommentStatsUsecase{repository: repository}
}

func (this *CommentStatsUsecase) Execute(articleID uint) (*entity.CommentStatsEntity, error) {
	stats, err := this.repository.GetStatsByArticleID(articleID)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
