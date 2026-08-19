package editionusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EditionCreateUseCase struct {
	repository       repository.EditionRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewEditionCreateUseCase(repository repository.EditionRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *EditionCreateUseCase {
	return &EditionCreateUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *EditionCreateUseCase) Execute(edition *entity.EditionEntity, user *entity.UserBasicEntity) error {
	allowed, err := this.memberPermission.Execute(user.Roles, edition.JournalID)
	if err != nil {
		return err
	}
	if !allowed {
		return response.PermissionDeniedError
	}
	return this.repository.Create(edition)
}
