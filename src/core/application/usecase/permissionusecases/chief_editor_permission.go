package permissionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ChiefEditorPermissionUseCase struct {
}

// @inject
func NewChiefEditorPermissionUseCase() *ChiefEditorPermissionUseCase {
	return &ChiefEditorPermissionUseCase{}
}

func (this *ChiefEditorPermissionUseCase) Execute(userRoles []*entity.UserRoleEntity, journalID uint) bool {
	for _, role := range userRoles {
		if role.JournalID != nil && *role.JournalID == journalID && role.Role == enum.RoleChiefEditor {
			return true
		}
	}
	return false
}
