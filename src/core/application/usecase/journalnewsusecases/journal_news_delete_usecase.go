package journalnewsusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalNewsDeleteUseCase struct {
	repository       repository.JournalNewsRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalNewsDeleteUseCase(repository repository.JournalNewsRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *JournalNewsDeleteUseCase {
	return &JournalNewsDeleteUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *JournalNewsDeleteUseCase) Execute(id uint, user *entity.UserBasicEntity) error {
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
