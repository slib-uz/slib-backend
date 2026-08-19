package roleusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserRoleAllUseCase struct {
	repository repository.UserRoleRepository
}

// @inject
func NewUserRoleAllUseCase(repository repository.UserRoleRepository) *UserRoleAllUseCase {
	return &UserRoleAllUseCase{
		repository: repository,
	}
}

func (this *UserRoleAllUseCase) Execute(userID uint) ([]*entity.UserRoleEntity, error) {
	userRoles, err := this.repository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	var results = make([]*entity.UserRoleEntity, len(userRoles))
	for i, userRole := range userRoles {
		results[i] = userRole
	}
	return results, nil

}
