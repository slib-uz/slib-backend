package commentusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticleCommentsListUseCase struct {
	repository repository.CommentRepository
}

// @inject
func NewArticleCommentsListUseCase(repository repository.CommentRepository) *ArticleCommentsListUseCase {
	return &ArticleCommentsListUseCase{repository: repository}
}

func (this *ArticleCommentsListUseCase) Execute(articleID uint, page, pageSize int) (*entity2.PagingEntity[entity2.CommentEntity], error) {
	paging, err := this.repository.GetByArticleID(articleID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return paging, nil
}
