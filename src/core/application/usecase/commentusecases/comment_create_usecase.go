package commentusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type CommentCreateUseCase struct {
	repository repository.CommentRepository
}

// @inject
func NewCommentCreateUseCase(repository repository.CommentRepository) *CommentCreateUseCase {
	return &CommentCreateUseCase{repository: repository}
}

func (this *CommentCreateUseCase) Execute(user *entity2.UserBasicEntity, comment *entity2.CommentEntity) error {
	comment.UserID = user.ID
	return this.repository.Create(comment)
}
