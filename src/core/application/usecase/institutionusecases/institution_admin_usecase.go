package institutionusecases

import (
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type InstitutionAdminUseCase struct {
	repository repository.UserRoleRepository
}

// @inject
func NewInstitutionAdminUseCase(repository repository.UserRoleRepository) *InstitutionAdminUseCase {
	return &InstitutionAdminUseCase{repository: repository}
}

func (this *InstitutionAdminUseCase) Execute(institutionId uint) ([]*entity.UserRoleWithBasicUserEntity, error) {
	_roles, err := this.repository.GetByInstitutionID(institutionId)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.UserRoleEntity, len(_roles))
	copy(users, _roles)
	return mapper.UserRoleEntityListToWithBasicUserEntityList(users), nil
}
