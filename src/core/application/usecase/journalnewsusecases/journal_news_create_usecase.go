package journalnewsusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalNewsCreateUseCase struct {
	repository       repository.JournalNewsRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalNewsCreateUseCase(repository repository.JournalNewsRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *JournalNewsCreateUseCase {
	return &JournalNewsCreateUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *JournalNewsCreateUseCase) Execute(news *entity.JournalNewsEntity, user *entity.UserBasicEntity) error {
	allowed, err := this.memberPermission.Execute(user.Roles, news.JournalID)
	if err != nil {
		return err
	}
	if !allowed {
		return response.PermissionDeniedError
	}
	return this.repository.Create(news)
}
