package commentusecases

import (
	"slib.uz/src/core/domain/ports/repository"
)

type CommentDeleteUseCase struct {
	repository repository.CommentRepository
}

// @inject
func NewCommentDeleteUseCase(repository repository.CommentRepository) *CommentDeleteUseCase {
	return &CommentDeleteUseCase{repository: repository}
}

func (this *CommentDeleteUseCase) Execute(userID, id uint) error {
	return this.repository.DeleteByIDAndUserID(id, userID)
}
