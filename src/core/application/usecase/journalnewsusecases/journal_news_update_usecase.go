package journalnewsusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalNewsUpdateUseCase struct {
	repository       repository.JournalNewsRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalNewsUpdateUseCase(repository repository.JournalNewsRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *JournalNewsUpdateUseCase {
	return &JournalNewsUpdateUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *JournalNewsUpdateUseCase) Execute(id uint, news *entity.JournalNewsEntity, user *entity.UserBasicEntity) error {
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
	return this.repository.Update(id, news)
}
