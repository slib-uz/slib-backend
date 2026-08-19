package journaleditorialusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalEditorialUpdateUseCase struct {
	repository       repository.JournalEditorialRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalEditorialUpdateUseCase(repository repository.JournalEditorialRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *JournalEditorialUpdateUseCase {
	return &JournalEditorialUpdateUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *JournalEditorialUpdateUseCase) Execute(id uint, editorial *entity.JournalEditorialEntity, user *entity.UserBasicEntity) error {
	existing, err := this.repository.GetByID(id)
	if err != nil {
		return err
	}
	allowed, err := this.memberPermission.Execute(user.Roles, existing.JournalID)
	if err != nil {
		return err
	}
	if !allowed {
		return response.PermissionDeniedError
	}
	return this.repository.Update(id, editorial)
}
