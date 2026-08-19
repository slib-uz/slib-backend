package userusecases

import (
	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type UserListUseCase struct {
	repository repository.UserRepository
}

// @inject
func NewUserListUseCase(repository repository.UserRepository) *UserListUseCase {
	return &UserListUseCase{repository: repository}
}

func (this *UserListUseCase) Execute(page, pageSize int, search string, user *entity2.UserBasicEntity) (*entity2.PagingEntity[entity2.UserEntity], error) {
	if user == nil {
		return nil, response.UnauthorizedError
	}

	hasAdminRole := false
	if user.Roles != nil {
		for _, role := range user.Roles {
			if role.Role == enum.RoleAdmin {
				hasAdminRole = true
				break
			}
		}
	}

	if !user.IsAdmin && !hasAdminRole {
		return nil, response.PermissionDeniedError
	}
	return this.repository.GetAll(page, pageSize, search)
}
