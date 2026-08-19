package journaleditorialusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalEditorialDeleteUseCase struct {
	repository       repository.JournalEditorialRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalEditorialDeleteUseCase(repository repository.JournalEditorialRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *JournalEditorialDeleteUseCase {
	return &JournalEditorialDeleteUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *JournalEditorialDeleteUseCase) Execute(id uint, user *entity.UserBasicEntity) error {
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
	return this.repository.Delete(id)
}
