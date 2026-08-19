package publisherusecases

import (
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type PublisherAdminUseCase struct {
	repository repository.UserRoleRepository
}

// @inject
func NewPublisherAdminUseCase(repository repository.UserRoleRepository) *PublisherAdminUseCase {
	return &PublisherAdminUseCase{repository: repository}
}

func (this *PublisherAdminUseCase) Execute(publisherId uint) ([]*entity.UserRoleWithBasicUserEntity, error) {
	_roles, err := this.repository.GetByPublisherID(publisherId)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.UserRoleEntity, len(_roles))
	copy(users, _roles)
	return mapper.UserRoleEntityListToWithBasicUserEntityList(users), nil
}
