package permissionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type PublisherAdminPermissionUseCase struct {
}

// @inject
func NewPublisherAdminPermissionUseCase() *PublisherAdminPermissionUseCase {
	return &PublisherAdminPermissionUseCase{}
}

func (this *PublisherAdminPermissionUseCase) Execute(roles []*entity.UserRoleEntity, publisherID uint) bool {
	for _, role := range roles {
		if role.Role == enum.RolePublisherAdmin && role.PublisherID != nil && *role.PublisherID == publisherID {
			return true
		}
	}
	return false
}
