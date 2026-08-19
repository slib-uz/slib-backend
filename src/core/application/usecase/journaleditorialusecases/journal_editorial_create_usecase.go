package journaleditorialusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalEditorialCreateUseCase struct {
	repository       repository.JournalEditorialRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalEditorialCreateUseCase(repository repository.JournalEditorialRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *JournalEditorialCreateUseCase {
	return &JournalEditorialCreateUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *JournalEditorialCreateUseCase) Execute(editorial *entity.JournalEditorialEntity, user *entity.UserBasicEntity) error {
	allowed, err := this.memberPermission.Execute(user.Roles, editorial.JournalID)
	if err != nil {
		return err
	}
	if !allowed {
		return response.PermissionDeniedError
	}
	return this.repository.Create(editorial)
}
