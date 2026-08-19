package permissionusecases

import "slib.uz/src/core/domain/entity"

type HasAnyPublisherRolePermissionUseCase struct {
}

// @inject
func NewHasAnyPublisherRolePermissionUseCase() *HasAnyPublisherRolePermissionUseCase {
	return &HasAnyPublisherRolePermissionUseCase{}
}

func (this *HasAnyPublisherRolePermissionUseCase) Execute(userRoles []*entity.UserRoleEntity, publisherID uint) bool {
	for _, role := range userRoles {
		if role.PublisherID != nil && *role.PublisherID == publisherID {
			return true
		}
	}
	return false
}
