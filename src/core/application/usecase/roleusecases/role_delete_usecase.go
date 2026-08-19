package roleusecases

import "slib.uz/src/core/domain/ports/repository"

type RoleDeleteUseCase struct {
	repository repository.UserRoleRepository
}

// @inject
func NewRoleDeleteUseCase(repository repository.UserRoleRepository) *RoleDeleteUseCase {
	return &RoleDeleteUseCase{repository: repository}
}

func (this *RoleDeleteUseCase) Execute(id uint) error {
	return this.repository.Delete(id)
}
